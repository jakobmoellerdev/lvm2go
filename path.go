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
