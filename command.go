package lvm2go

import (
	"context"
	"os/exec"
	"time"
)

const (
	nsenter = "/usr/bin/nsenter"
)

var waitDelayKey = struct{}{}

var DefaultWaitDelay = time.Duration(0)

func SetProcessCancelWaitDelay(ctx context.Context, delay time.Duration) context.Context {
	return context.WithValue(ctx, waitDelayKey, delay)
}

func GetProcessCancelWaitDelay(ctx context.Context) time.Duration {
	if delay, ok := ctx.Value(waitDelayKey).(time.Duration); ok {
		return delay
	}
	return DefaultWaitDelay
}

// CommandContext creates exec.Cmd with custom args. it is equivalent to exec.Command(cmd, args...) when not containerized.
// When containerized, it calls nsenter with the provided command and args.
func CommandContext(ctx context.Context, cmd string, args ...string) *exec.Cmd {
	var c *exec.Cmd

	if IsContainerized(ctx) {
		args = append([]string{"-m", "-u", "-i", "-n", "-p", "-t", "1", cmd}, args...)
		c = exec.CommandContext(ctx, nsenter, args...)
	} else {
		c = exec.CommandContext(ctx, cmd, args...)
	}
	c.WaitDelay = GetProcessCancelWaitDelay(ctx)

	return CommandWithCustomEnvironment(ctx, c)
}
