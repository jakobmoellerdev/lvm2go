package lvm2go

import (
	"fmt"
	"os"
	"sync"
)

const (
	DefaultLVMSystemDir            = "/etc/lvm"
	LVMSystemDirEnv                = "LVM_SYSTEM_DIR"
	LVMGlobalConfigurationFileName = "lvm.conf"
	LVMLocalConfigurationFileName  = "lvmlocal.conf"
)

// LVMSystemDir returns the system directory for LVM configuration files.
// If the LVM_SYSTEM_DIR environment variable is set, its value will be returned.
// Otherwise, the default value "/etc/lvm" will be returned.
var LVMSystemDir = sync.OnceValue[string](lvmSystemDir)

func lvmSystemDir() string {
	if dir := os.Getenv(LVMSystemDirEnv); dir != "" {
		return dir
	}
	return DefaultLVMSystemDir
}

// LVMGlobalConfiguration is the path to the global LVM configuration file.
// It usually defaults to "/etc/lvm/lvm.conf", but can be changed by setting the LVM_SYSTEM_DIR environment variable
var LVMGlobalConfiguration = fmt.Sprintf("%s/%s", LVMSystemDir(), LVMGlobalConfigurationFileName)

// LVMLocalConfiguration is the path to the local LVM configuration file.
// It usually defaults to "/etc/lvm/lvmlocal.conf", but can be changed by setting the LVM_SYSTEM_DIR environment variable
var LVMLocalConfiguration = fmt.Sprintf("%s/%s", LVMSystemDir(), LVMLocalConfigurationFileName)
