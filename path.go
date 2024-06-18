package lvm2go

import (
	"sync"
)

var (
	lvmBinaryPathLock = &sync.RWMutex{}
	lvmBinaryPath     = "/sbin/lvm"
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
	lvmBinaryPathLock.RLock()
	defer lvmBinaryPathLock.RUnlock()
	return lvmBinaryPath
}
