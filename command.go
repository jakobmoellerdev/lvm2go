package lvm2go

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"sync"
	"time"
)

const (
	nsenter               = "/usr/bin/nsenter"
	DefaultVolumeGroupEnv = "LVM_VG_NAME"
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

	if DefaultVolumeGroup(ctx) != "" {
		c.Env = append(c.Env, fmt.Sprintf("%s=%s", DefaultVolumeGroupEnv, DefaultVolumeGroup(ctx)))
	}

	return CommandWithCustomEnvironment(ctx, c)
}

var defaultVolumeGroupKey = struct{}{}

func WithDefaultVolumeGroup(ctx context.Context, vg string) context.Context {
	return context.WithValue(ctx, defaultVolumeGroupKey, vg)
}

func DefaultVolumeGroup(ctx context.Context) string {
	if vg, ok := ctx.Value(defaultVolumeGroupKey).(string); ok {
		return vg
	}
	return ""
}

var (
	isContainerized     bool
	detectContainerized sync.Once
)

func IsContainerized(ctx context.Context) bool {
	detectContainerized.Do(func() {
		if _, err := os.Stat("/.dockerenv"); err == nil {
			isContainerized = true
		} else if _, err := os.Stat("/.containerenv"); err == nil {
			isContainerized = true
		}
		if isContainerized {
			slog.DebugContext(ctx, "lvm2go is running in docker environment")
		}
	})
	return isContainerized
}

var envContextKey = struct{}{}

func WithCustomEnvironment(ctx context.Context, env map[string]string) context.Context {
	return context.WithValue(ctx, envContextKey, env)
}

func GetCustomEnvironment(ctx context.Context) map[string]string {
	if env, ok := ctx.Value(envContextKey).(map[string]string); ok {
		return env
	}
	return nil
}

func CommandWithCustomEnvironment(ctx context.Context, cmd *exec.Cmd) *exec.Cmd {
	if UseStandardLocale() {
		cmd.Env = append(cmd.Env, "LC_ALL=C")
	}
	if env := GetCustomEnvironment(ctx); env != nil {
		for k, v := range env {
			cmd.Env = append(cmd.Env, k+"="+v)
		}
	} else {
		cmd.Env = os.Environ()
	}
	return cmd
}
