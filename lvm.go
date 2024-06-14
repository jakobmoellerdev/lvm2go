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

// RunLVM calls lvmBinaryPath sub-commands and prints the output to the log.
func RunLVM(ctx context.Context, args ...string) error {
	return RunLVMInto(ctx, nil, args...)
}

// RunLVMInto calls lvmBinaryPath sub-commands and decodes the output via JSON into the provided struct pointer.
// if the struct pointer is nil, the output will be printed to the log instead.
func RunLVMInto(ctx context.Context, into any, args ...string) error {
	output, err := RunLVMStreamed(ctx, args...)
	if err != nil {
		return fmt.Errorf("failed to execute command: %v", err)
	}

	// if we don't decode the output into a struct, we can still log the command results from stdout.
	if into == nil {
		scanner := bufio.NewScanner(output)
		for scanner.Scan() {
			slog.Info(strings.TrimSpace(scanner.Text()))
		}
		err = scanner.Err()
	} else {
		err = json.NewDecoder(output).Decode(&into)
	}
	closeErr := output.Close()

	return errors.Join(closeErr, err)
}

// RunLVMStreamed calls lvmBinaryPath sub-commands and returns the output as a ReadCloser.
// The caller is responsible for closing the ReadCloser, which will cause the command to complete.
// Not calling close on this method will result in a resource leak.
func RunLVMStreamed(ctx context.Context, args ...string) (io.ReadCloser, error) {
	return runCommand(ctx, Command(GetLVMPath(), args...))
}
