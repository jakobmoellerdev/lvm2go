package lvm2go

import (
	"os/exec"
	"sync"
)

var (
	lvmBinaryPathLock = &sync.Mutex{}
	lvmBinaryPath     = ""
)

// SetLVMPath sets the Path to the lvmBinaryPath command.
func SetLVMPath(path string) {
	lvmBinaryPathLock.Lock()
	defer lvmBinaryPathLock.Unlock()
	if path != "" {
		lvmBinaryPath = path
	}
}

// GetLVMPath returns the Path to the lvmBinaryPath command.
func GetLVMPath() string {
	lvmBinaryPathLock.Lock()
	defer lvmBinaryPathLock.Unlock()

	if lvmBinaryPath == "" {
		lvmBinaryPath = resolveLVMPathFromHost()
	}

	return lvmBinaryPath
}

var resolveLVMPathFromHost = sync.OnceValue(func() string {
	if path, err := exec.LookPath("lvm"); err != nil {
		return "/usr/sbin/lvm"
	} else {
		return path
	}
})
