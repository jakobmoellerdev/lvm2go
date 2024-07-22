package lvm2go

import (
	"regexp"
)

var (
	// NotFoundPattern is a regular expression that matches the error message when a volume group or logical volume is not found.
	// The volume group might not be present or the logical volume might not be present in the volume group.
	NotFoundPattern = regexp.MustCompile(`Volume group "(.*?)" not found|Failed to find logical volume "(.*?)"`)

	NoSuchCommandPattern = regexp.MustCompile(`RequestConfirm such command`)
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

// IsLVMNoSuchCommand returns true if the error is a LVM recognized error and it determined that the command
// is not found.
func IsLVMNoSuchCommand(err error) bool {
	lvmErr, ok := AsExitCodeError(err)

	// If the exit code is not 2, it is guaranteed that the error is not a not found error.
	if !ok || lvmErr.ExitCode() != 2 {
		return false
	}

	return NoSuchCommandPattern.Match([]byte(lvmErr.Error()))
}
