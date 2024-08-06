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
		return fmt.Errorf("failed to get struct processor and query for config Decode: %v", err)
	}

	queryArgs := append(query, args.GetRaw()...)

	return c.RunLVMRaw(ctx, processor, append([]string{"config"}, queryArgs...)...)
}

func (c *client) WriteAndEncodeConfig(_ context.Context, v any, writer io.Writer) error {
	return NewLexingConfigEncoder(writer).Encode(v)
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

func (c *client) UpdateGlobalConfig(ctx context.Context, v any) error {
	return c.UpdateConfigFromPath(ctx, v, LVMGlobalConfiguration)
}

func (c *client) UpdateLocalConfig(ctx context.Context, v any) error {
	return c.UpdateConfigFromPath(ctx, v, LVMLocalConfiguration)
}

func (c *client) UpdateProfileConfig(ctx context.Context, v any, profile Profile) error {
	path, err := c.GetProfilePath(ctx, profile)
	if err != nil {
		return fmt.Errorf("failed to get profile path: %v", err)
	}
	return c.UpdateConfigFromPath(ctx, v, path)
}

func (c *client) UpdateConfigFromPath(ctx context.Context, v any, path string) error {
	fileMode := os.FileMode(0600)
	profileFile, err := os.OpenFile(path, os.O_RDWR, fileMode)
	if err != nil {
		return fmt.Errorf("failed to open config for read/write (%s): %v", fileMode, err)
	}
	defer func() {
		err = errors.Join(err, profileFile.Close())
	}()
	if err = updateConfig(ctx, v, profileFile); err != nil {
		err = fmt.Errorf("failed to update config at %s: %v", path, err)
	}
	return err
}

// updateConfig updates the configuration file with the new values from the struct v.
// The configuration file is read and written from the provided io.ReadWriteSeeker.
// The configuration file is updated with the new values from the struct v.
// If a field is not present in the configuration file, it is added with a comment to indicate it was added.
// If a field is present in the configuration file, it is updated with the new value and a comment to indicate it was edited.
// If the resulting configuration is smaller than the original, the difference is padded with empty bytes.
// The configuration file is written back to the start of original configuration file.
func updateConfig(ctx context.Context, v any, rw io.ReadWriteSeeker) error {
	structMappings, err := DecodeLVMStructTagFieldMappings(v)
	if err != nil {
		return fmt.Errorf("failed to read lvm struct tag: %v", err)
	}
	tokensToModify := StructMappingsToConfigTokens(structMappings)

	data, err := io.ReadAll(rw)
	if err != nil {
		return fmt.Errorf("failed to read configuration: %v", err)
	}
	reader := bytes.NewReader(data)

	tokensFromFile, err := NewBufferedConfigLexer(reader).Lex()
	if err != nil {
		return fmt.Errorf("failed to read configuration: %v", err)
	}

	// First merge all assignments from the new struct into the existing configuration
	newTokens := assignmentsWithSections(tokensFromFile).
		overrideWith(assignmentsWithSections(tokensToModify))

	// Then append any new assignments at the end of the sections
	tokens := appendAssignmentsAtEndOfSections(tokensFromFile, newTokens)

	// Write the new configuration to a buffer
	buf := bytes.NewBuffer(make([]byte, 0, len(data)))
	if err := NewLexingConfigEncoder(buf).Encode(tokens); err != nil {
		return fmt.Errorf("failed to encode new configuration: %v", err)
	}

	if diff := buf.Cap() - len(data); diff < 0 {
		// If the old configuration is smaller than the new configuration, we need to append the difference
		// with empty bytes to ensure we do not have leftover data from the old configuration
		buf.Write(make([]byte, -diff))
	}
	// We want to write from the start, so seek back to the start of the configuration
	if _, err = rw.Seek(int64(-len(data)), io.SeekCurrent); err != nil {
		return fmt.Errorf("failed to seek to start of configuration: %v", err)
	}

	// Write the new configuration based on the old configuration with the new fields
	return copyWithTimeout(ctx, rw, buf, 10*time.Second)
}

// generateLVMConfigEditComment generates a comment to be added to the configuration file
// This comment is used to indicate that the field was edited by the client.
func generateLVMConfigEditComment() string {
	return fmt.Sprintf(`This field was edited by %s at %s`, ModuleID(), time.Now().Format(time.RFC3339))
}

func generateLVMConfigCreateComment() string {
	return fmt.Sprintf(`configuration created by %s at %s`, ModuleID(), time.Now().Format(time.RFC3339))
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
	fieldsForConfigQuery, err := DecodeLVMStructTagFieldMappings(v)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read lvm struct tag: %v", err)
	}

	var query []string
	for _, field := range fieldsForConfigQuery {
		query = append(query, fmt.Sprintf("%s/%s", field.prefix, field.name))
	}

	return func(out io.Reader) error {
		return newLexingConfigDecoderWithFieldMapping(out, fieldsForConfigQuery).Decode()
	}, query, nil
}

// copyWithTimeout copies data from r to w with a timeout.
// If the operation takes longer than the timeout, an error is returned.
// If the operation completes before the timeout, the error as returned by io.Copy is returned.
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
