/*
 Copyright 2024 The lvm2go Authors.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package lvm2go

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var ErrProfileNameEmpty = errors.New("profile name is empty")

const LVMConfigStructTag = "lvm"
const LVMProfileExtension = ".profile"

type (
	ConfigOptions struct {
		ConfigType
		Profile
	}
	ConfigOption interface {
		ApplyToConfigOptions(opts *ConfigOptions)
	}
	ConfigOptionsList []ConfigOption
)

var (
	_ ArgumentGenerator = ConfigOptionsList{}
	_ Argument          = (*ConfigOptions)(nil)
)

type RawConfig map[string]any

func (c *client) RawConfig(ctx context.Context, opts ...ConfigOption) (RawConfig, error) {
	args, err := ConfigOptionsList(opts).AsArgs()
	if err != nil {
		return nil, err
	}

	entries := map[string]any{}
	processor := RawOutputProcessor(func(out io.Reader) error {
		scanner := bufio.NewScanner(out)
		start := scanner.Scan()
		if !start {
			return errors.New("no output")
		}

		for scanner.Scan() {
			if bytes.TrimSpace(scanner.Bytes())[0] == '}' {
				startingToken := scanner.Scan()
				if !startingToken {
					break // end of output
				} else {
					continue // skip the closing brace
				}
			}

			trimmed := bytes.TrimSpace(scanner.Bytes())
			entrySepIdx := bytes.Index(trimmed, []byte("="))
			if entrySepIdx == -1 {
				return fmt.Errorf("unexpected end of entry without assignment: %s", scanner.Text())
			}
			key := trimmed[:entrySepIdx]
			value := trimmed[entrySepIdx+1:]

			if len(value) == 0 {
				entries[string(key)] = nil
				continue
			}

			if bytes.Contains(value, []byte("\"")) {
				if bytes.Contains(value, []byte("[")) {
					elems := strings.Split(string(bytes.Trim(value, "[]")), ",")
					for i, elem := range elems {
						elems[i] = strings.Trim(elem, "\"")
					}
					entries[string(key)] = elems
				} else {
					entries[string(key)] = string(bytes.Trim(value, "\""))
				}

				continue
			} else if bytes.Contains(value, []byte("[]")) {
				entries[string(key)] = []any{}
				continue
			}

			if parsed, err := strconv.ParseInt(string(value), 10, 64); err != nil {
				return err
			} else {
				entries[string(key)] = parsed
			}
		}

		return scanner.Err()
	})

	if err := c.RunLVMRaw(ctx, processor, append([]string{"config"}, args.GetRaw()...)...); err != nil {
		return nil, err
	}

	return entries, nil
}

func (c *client) ReadAndDecodeConfig(ctx context.Context, v any, opts ...ConfigOption) error {
	args, err := ConfigOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	processor, query, err := getStructProcessorAndQuery(v)
	if err != nil {
		return fmt.Errorf("failed to get struct processor and query for config decode: %v", err)
	}

	queryArgs := append(query, args.GetRaw()...)

	return c.RunLVMRaw(ctx, processor, append([]string{"config"}, queryArgs...)...)
}

func (c *client) WriteAndEncodeConfig(ctx context.Context, v any, writer io.Writer) error {
	fieldsForConfigQuery, err := readLVMStructTag(v)
	if err != nil {
		return fmt.Errorf("failed to read lvm struct tag: %v", err)
	}

	fieldsByPrefix := map[string][]lvmStructTagFieldSpec{}
	for _, field := range fieldsForConfigQuery {
		fieldsByPrefix[field.prefix] = append(fieldsByPrefix[field.prefix], field)
	}

	for prefix, field := range fieldsByPrefix {
		buf := bytes.Buffer{}
		if _, err = buf.WriteString(fmt.Sprintf("%s {\n", prefix)); err != nil {
			return fmt.Errorf("failed to block start %s: %v", prefix, err)
		}
		for _, f := range field {
			switch f.Kind() {
			case reflect.Int:
				if _, err = buf.WriteString(fmt.Sprintf("\t%s=%d\n", f.name, f.Value.Int())); err != nil {
					return fmt.Errorf("failed to write field %s: %v", f.name, err)
				}
			case reflect.String:
				if _, err = buf.WriteString(fmt.Sprintf("\t%s=%q\n", f.name, f.Value.String())); err != nil {
					return fmt.Errorf("failed to write field %s: %v", f.name, err)
				}
			default:
				return fmt.Errorf("unsupported field type %s", f.Kind())
			}
		}
		if _, err = buf.WriteString("}\n"); err != nil {
			return fmt.Errorf("failed to block end %s: %v", prefix, err)
		}

		if err = copyWithTimeout(ctx, writer, &buf, 10*time.Second); err != nil {
			return fmt.Errorf("failed to write config block: %v", err)
		}
	}

	return nil
}

func (c *client) GetProfileDirectory(ctx context.Context) (string, error) {
	type lvmConfig struct {
		Config struct {
			ProfileDir string `lvm:"profile_dir"`
		} `lvm:"config"`
	}
	cfg := &lvmConfig{}
	if err := c.ReadAndDecodeConfig(ctx, cfg, ConfigTypeFull); err != nil {
		return "", fmt.Errorf("failed to get lvm profile directory: %v", err)
	}
	return cfg.Config.ProfileDir, nil
}

func (c *client) CreateProfile(ctx context.Context, v any, profile Profile) (string, error) {
	path, err := c.GetProfilePath(ctx, profile)
	if err != nil {
		return "", err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return "", fmt.Errorf("failed to create profile file: %v", err)
	}
	defer func() {
		err = errors.Join(err, file.Close())
	}()

	if err = c.WriteAndEncodeConfig(ctx, v, file); err != nil {
		err = errors.Join(err, fmt.Errorf("failed to write config to profile: %v", err))
	}

	return path, err
}

func (c *client) RemoveProfile(ctx context.Context, profile Profile) error {
	path, err := c.GetProfilePath(ctx, profile)
	if err != nil {
		return err
	}
	return os.Remove(path)
}

func (c *client) GetProfilePath(ctx context.Context, profile Profile) (string, error) {
	if profile == "" {
		return "", ErrProfileNameEmpty
	}

	dir, err := c.GetProfileDirectory(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get profile directory: %v", err)
	}
	path := string(profile)

	split := strings.Split(string(profile), string(filepath.Separator))
	if len(split) > 1 {
		splitDir := strings.Join(split[:len(split)-1], string(filepath.Separator))
		if splitDir != "" && dir != splitDir {
			return "", fmt.Errorf("unexpected profile directory: %s, should be %s", splitDir, dir)
		}
		path = split[len(split)-1]
	}

	if ext := filepath.Ext(path); ext == "" {
		path = fmt.Sprintf("%s%s", path, LVMProfileExtension)
	} else if ext != ".profile" {
		return "", fmt.Errorf("%q is an invalid profile extension: %w", ext, ErrInvalidProfileExtension)
	}

	return filepath.Join(dir, path), nil
}

// GetFromRawConfig retrieves a value from a RawConfig by key and attempts to cast it to the type of T.
// If the key is not found, an error is returned.
// If the key is found but the value is not of type T, an error is returned.
func GetFromRawConfig[T any](config RawConfig, key string) (T, error) {
	if value, ok := config[key]; !ok {
		return *new(T), fmt.Errorf("key %s not found in config", key)
	} else {
		if typed, ok := value.(T); !ok {
			return *new(T), fmt.Errorf("key %s is not of type %T", key, typed)
		} else {
			return typed, nil
		}
	}
}

func (opts *ConfigOptions) ApplyToConfigOptions(new *ConfigOptions) {
	*new = *opts
}

func (list ConfigOptionsList) AsArgs() (Arguments, error) {
	args := NewArgs(ArgsTypeGeneric)
	options := ConfigOptions{}
	for _, opt := range list {
		opt.ApplyToConfigOptions(&options)
	}
	if err := options.ApplyToArgs(args); err != nil {
		return nil, err
	}
	return args, nil
}

func (opts *ConfigOptions) ApplyToVersionOptions(new *ConfigOptions) {
	*new = *opts
}

func (opts *ConfigOptions) ApplyToArgs(args Arguments) error {
	for _, arg := range []Argument{
		opts.ConfigType,
		opts.Profile,
	} {
		if err := arg.ApplyToArgs(args); err != nil {
			return err
		}
	}

	return nil
}

func getStructProcessorAndQuery(v any) (RawOutputProcessor, []string, error) {
	fieldsForConfigQuery, err := readLVMStructTag(v)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read lvm struct tag: %v", err)
	}

	var query []string
	for _, field := range fieldsForConfigQuery {
		query = append(query, fmt.Sprintf("%s/%s", field.prefix, field.name))
	}

	return func(out io.Reader) error {
		scanner := bufio.NewScanner(out)
		for scanner.Scan() {
			split := strings.Split(scanner.Text(), "=")
			if len(split) != 2 {
				return fmt.Errorf("unexpected line (no key value identification): %s", scanner.Text())
			}
			k, v := split[0], strings.Trim(split[1], "\"")

			if field, ok := fieldsForConfigQuery[k]; ok {
				if field.Kind() == reflect.String {
					field.SetString(v)
				} else if field.Kind() == reflect.Int64 {
					if parsed, err := strconv.ParseInt(v, 10, 64); err != nil {
						return fmt.Errorf("failed to parse int64 for field %s: %v", k, err)
					} else {
						field.SetInt(parsed)
					}
				} else {
					return fmt.Errorf("unsupported field type %s", field.Kind())
				}
			}
		}

		return scanner.Err()
	}, query, nil
}

type lvmStructTagFieldSpec struct {
	prefix string
	name   string
	reflect.Value
}

func readLVMStructTag(v any) (map[string]lvmStructTagFieldSpec, error) {
	fields, typeAccessor, valueAccessor, err := accessStructOrPointerToStruct(v)
	if err != nil {
		return nil, err
	}

	tagOrIgnore := func(tag reflect.StructTag) (string, bool) {
		return tag.Get(LVMConfigStructTag), tag.Get(LVMConfigStructTag) == "-"
	}

	fieldSpecs := make(map[string]lvmStructTagFieldSpec)
	for i := range fields {
		outerField := typeAccessor(i)
		prefix, ignore := tagOrIgnore(outerField.Tag)
		if ignore {
			continue
		}
		fields, typeAccessor, valueAccessor, err := accessStructOrPointerToStruct(valueAccessor(i))
		if err != nil {
			return nil, err
		}
		for j := range fields {
			innerField := typeAccessor(j)
			name, ignore := tagOrIgnore(innerField.Tag)
			if ignore {
				continue
			}
			fieldSpecs[name] = lvmStructTagFieldSpec{
				prefix,
				name,
				valueAccessor(j),
			}
		}
	}
	return fieldSpecs, nil
}

func copyWithTimeout(ctx context.Context, w io.Writer, r io.Reader, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		_, err := io.Copy(w, r)
		done <- err
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("write operation timed out")
	case err := <-done:
		return err
	}
}
