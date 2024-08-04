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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
)

// RunLVM calls lvm2 sub-commands and prints the output to the log.
func (c *client) RunLVM(ctx context.Context, args ...string) error {
	return c.RunLVMInto(ctx, nil, args...)
}

// RunLVMInto calls lvm2 sub-commands and decodes the output via JSON into the provided struct pointer.
// if the struct pointer is nil, the output will be printed to the log instead.
func (c *client) RunLVMInto(ctx context.Context, into any, args ...string) error {
	cmd := CommandContext(ctx, GetLVMPath(), args...)

	output, err := StreamedCommand(ctx, cmd)
	if err != nil {
		return fmt.Errorf("failed to execute command: %v", err)
	}

	// if we don't decode the output into a struct, we can still log the command results from stdout.
	if into == nil {
		scanner := bufio.NewScanner(output)
		for scanner.Scan() {
			slog.InfoContext(ctx, strings.TrimSpace(scanner.Text()))
		}
		err = scanner.Err()
	} else {
		err = json.NewDecoder(output).Decode(&into)
	}

	err = errors.Join(output.Close(), err)

	if IsLVMErrNoSuchCommand(err) {
		return fmt.Errorf("%q is not a valid command: %w", strings.Join(args, " "), err)
	}

	return err
}

func (c *client) RunLVMRaw(ctx context.Context, process RawOutputProcessor, args ...string) error {
	return c.RunRaw(ctx, process, append([]string{GetLVMPath()}, args...)...)
}

type RawOutputProcessor func(out io.Reader) error

func NoOpRawOutputProcessor(expectOutput bool) RawOutputProcessor {
	return func(out io.Reader) error {
		data, err := io.ReadAll(out)
		if err != nil {
			return fmt.Errorf("failed to read output: %v", err)
		}
		if expectOutput && len(data) == 0 {
			return fmt.Errorf("expected output but got none")
		}
		if !expectOutput && len(data) > 0 {
			return fmt.Errorf("expected no output but got: %s", string(data))
		}
		return nil
	}
}

func (c *client) RunRaw(ctx context.Context, process RawOutputProcessor, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("no command provided")
	}
	cmd := CommandContext(ctx, args[0], args[1:]...)

	output, err := StreamedCommand(ctx, cmd)
	if err != nil {
		return fmt.Errorf("failed to execute command: %v", err)
	}
	err = process(output)
	closeErr := output.Close()
	return errors.Join(closeErr, err)
}
