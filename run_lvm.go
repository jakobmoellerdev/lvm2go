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
	output, err := StreamedCommand(ctx, CommandContext(ctx, GetLVMPath(), args...))
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
	closeErr := output.Close()

	return errors.Join(closeErr, err)
}

type RawOutputProcessor func(out io.Reader) error

func (c *client) RunLVMRaw(ctx context.Context, process RawOutputProcessor, args ...string) error {
	output, err := StreamedCommand(ctx, CommandContext(ctx, GetLVMPath(), args...))
	if err != nil {
		return fmt.Errorf("failed to execute command: %v", err)
	}
	err = process(output)
	closeErr := output.Close()
	return errors.Join(closeErr, err)
}
