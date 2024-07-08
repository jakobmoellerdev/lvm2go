package lvm2go

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type (
	ConfigOptions struct {
		ConfigType
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
	if err := opts.ConfigType.ApplyToArgs(args); err != nil {
		return err
	}

	return nil
}
