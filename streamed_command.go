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
	"io"
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

// StreamedCommand runs the command and returns the stdout as a ReadCloser that also Waits for the command to finish.
// After the Close command is called the cmd is closed and the resources are released.
// Not calling close on this method will result in a resource leak.
func StreamedCommand(ctx context.Context, cmd *exec.Cmd) (io.ReadCloser, error) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stdoutClose := func() error {
		return ignoreClosed(stdout.Close())
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, errors.Join(err, stdoutClose())
	}
	stderrClose := func() error {
		return ignoreClosed(stderr.Close())
	}

	slog.DebugContext(ctx, "running command", slog.String("command", strings.Join(cmd.Args, " ")))

	cmd.Cancel = func() error {
		slog.WarnContext(ctx, "killing streamed command process due to ctx cancel")

		return errors.Join(cmd.Process.Kill(), stdoutClose(), stderrClose())
	}

	if err := cmd.Start(); err != nil {
		return nil, errors.Join(err, stdoutClose(), stderrClose())
	}

	// Return a read closer that will wait for the command to finish when closed to release all resources.
	return &commandReadCloser{cmd: cmd, ReadCloser: stdout, stderr: stderr}, nil
}

// commandReadCloser is a ReadCloser that calls the Wait function of the command when Close is called.
// This is used to wait for the command the pipe before waiting for the command to finish.
type commandReadCloser struct {
	cmd *exec.Cmd
	io.ReadCloser
	stderr io.ReadCloser
}

// Close closes stdout and stderr and waits for the command to exit. Close
// should not be called before all reads from stdout have completed.
func (p *commandReadCloser) Close() error {
	var err error

	// Fully Read the pipes before waiting for the command to finish.
	stderr, stderrReadAllErr := io.ReadAll(p.stderr)
	err = errors.Join(err, stderrReadAllErr)
	stdout, stdoutReadAllErr := io.ReadAll(p.ReadCloser)
	err = errors.Join(err, stdoutReadAllErr)

	// create an error out of the stderr output if necessary
	err = errors.Join(err, NewLVMStdErr(stderr))

	// wait can result in an exit code error
	err = errors.Join(err, NewExitCodeError(p.cmd.Wait()))

	if len(stdout) > 0 {
		slog.Warn("STDOUT still contained data after the command finished")
		scanner := bufio.NewScanner(bytes.NewReader(stdout))
		for scanner.Scan() {
			slog.Warn(strings.TrimSpace(scanner.Text()))
		}
	}

	return err
}

func ignoreClosed(err error) error {
	if errors.Is(err, os.ErrClosed) {
		return nil
	}
	return err
}
