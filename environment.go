package lvm2go

import (
	"context"
	"os"
	"os/exec"
)

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
	cmd.Env = os.Environ()
	if UseStandardLocale() {
		cmd.Env = append(cmd.Env, "LC_ALL=C")
	}
	if env := GetCustomEnvironment(ctx); env != nil {
		for k, v := range env {
			cmd.Env = append(cmd.Env, k+"="+v)
		}
	}
	return cmd
}
