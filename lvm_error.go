package lvm2go

import (
	"fmt"
	"regexp"
	"runtime"
)

var (
	// NotFoundPattern is a regular expression that matches the error message when a volume group or logical volume is not found.
	// The volume group might not be present or the logical volume might not be present in the volume group.
	NotFoundPattern = regexp.MustCompile(`Volume group "(.*?)" not found|Failed to find logical volume "(.*?)"`)
)

// IsLVMNotFound returns true if the error is a LVM recognized error and it determined that either
// the underlying volume group or logical volume is not found.
func IsLVMNotFound(err error) bool {
	lvmErr, ok := AsExitCodeError(err)

	// If the exit code is not 5, it is guaranteed that the error is not a not found error.
	if !ok || lvmErr.ExitCode() != 5 {
		return false
	}

	return NotFoundPattern.Match([]byte(lvmErr.Error()))
}

func errFromToArgs(err error) error {
	pc, _, _, _ := runtime.Caller(2) // skip 2 frames (this function and the ApplyToArgs call)
	return fmt.Errorf("%s: %v", runtime.FuncForPC(pc).Name(), err)
}
