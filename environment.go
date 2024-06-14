package lvm2go

import (
	"os"
	"os/exec"
)

func CommandWithOSEnvironment(cmd *exec.Cmd) *exec.Cmd {
	cmd.Env = os.Environ()
	if UseStandardLocale() {
		cmd.Env = append(cmd.Env, "LC_ALL=C")
	}
	return cmd
}
