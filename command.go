package lvm2go

import (
	"context"
	"os/exec"
)

const (
	nsenter = "/usr/bin/nsenter"
)

// CommandContext creates exec.Cmd with custom args. it is equivalent to exec.Command(cmd, args...) when not containerized.
// When containerized, it calls nsenter with the provided command and args.
func CommandContext(ctx context.Context, cmd string, args ...string) *exec.Cmd {
	var c *exec.Cmd

	if IsContainerized() {
		args = append([]string{"-m", "-u", "-i", "-n", "-p", "-t", "1", cmd}, args...)
		c = exec.CommandContext(ctx, nsenter, args...)
	} else {
		c = exec.CommandContext(ctx, cmd, args...)
	}

	return CommandWithOSEnvironment(c)
}
