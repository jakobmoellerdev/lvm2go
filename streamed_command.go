package lvm2go

import (
	"context"
	"io"
	"log/slog"
	"os/exec"
)

// StreamedCommand runs the command and returns the stdout as a ReadCloser that also Waits for the command to finish.
// After the Close command is called the cmd is closed and the resources are released.
// Not calling close on this method will result in a resource leak.
func StreamedCommand(ctx context.Context, cmd *exec.Cmd) (io.ReadCloser, error) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		_ = stdout.Close()
		return nil, err
	}

	slog.DebugContext(ctx, "invoking command", "args", cmd.Args, "env", cmd.Env, "pwd", cmd.Dir)

	cmd.Cancel = func() error {
		slog.WarnContext(ctx, "killing streamed command process due to ctx cancel")
		if err := stdout.Close(); err != nil {
			return err
		}
		if err := stderr.Close(); err != nil {
			return err
		}
		return cmd.Process.Kill()
	}

	if err := cmd.Start(); err != nil {
		_ = stdout.Close()
		_ = stderr.Close()
		return nil, err
	}
	// Return a read closer that will wait for the command to finish when closed to release all resources.
	return commandReadCloser{cmd: cmd, ReadCloser: stdout, stderr: stderr}, nil
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
func (p commandReadCloser) Close() error {
	// Read the stderr output after the read has finished since we are sure by then the command must have run.
	stderr, err := io.ReadAll(p.stderr)
	if err != nil {
		return err
	}

	if err := p.cmd.Wait(); err != nil {
		// wait can result in an exit code error
		return NewExitCodeError(err, stderr)
	}
	return nil
}
