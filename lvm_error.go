package lvm2go

import (
	"regexp"
)

var (
	// NotFoundPattern is a regular expression that matches the error message when a volume group or logical volume is not found.
	// The volume group might not be present or the logical volume might not be present in the volume group.
	NotFoundPattern = regexp.MustCompile(`Volume group "(.*?)" not found|Failed to find logical volume "(.*?)"`)

	// NoSuchCommandPattern is a regular expression that matches the error message when a command is not found.
	NoSuchCommandPattern = regexp.MustCompile(`no such command`)

	// MaximumNumberOfLogicalVolumesPattern is a regular expression that matches the error message when the maximum number of logical volumes is reached.
	MaximumNumberOfLogicalVolumesPattern = regexp.MustCompile(`Maximum number of logical volumes \(\d+\) reached in volume group (.*?)`)

	// MaximumNumberOfPhysicalVolumesPattern is a regular expression that matches the error message when the maximum number of physical volumes is reached.
	MaximumNumberOfPhysicalVolumesPattern = regexp.MustCompile(`No space for '(.*?)' - volume group '(.*?)' holds max \d+ physical volume\(s\)\.`)

	// CannotChangeVGWhilePVsAreMissingPattern is a regular expression that matches the error message when the volume group is immutable because physical volumes are missing.
	CannotChangeVGWhilePVsAreMissingPattern = regexp.MustCompile(`Cannot change VG (.*?) while PVs are missing\.`)

	// WarningCouldNotFindDeviceWithUUIDPattern is a regular expression that matches the error message when a device with a specific UUID is not found.
	WarningCouldNotFindDeviceWithUUIDPattern = regexp.MustCompile(`WARNING: Couldn't find device with uuid [\w-]+\.`)

	// VGMissingPVsPattern is a regular expression that matches the error message when a volume group is missing physical volumes.
	VGMissingPVsPattern = regexp.MustCompile(`VG (.*?) is missing PV (.*?) \(last written to (.*?)\)`)

	// ThereAreStillPartialLVsPattern is a regular expression that matches the error message when there are still partial logical volumes in a volume group.
	ThereAreStillPartialLVsPattern = regexp.MustCompile(`There are still partial LVs in VG (.*?)\.`)

	// WarningPartialLVNeedsRepairOrRemovePattern is a regular expression that matches the error message when a logical volume needs repair or remove.
	WarningPartialLVNeedsRepairOrRemovePattern = regexp.MustCompile(`WARNING: Partial LV (.*?) needs to be repaired or removed\.`)

	// NoDataToMovePattern is a regular expression that matches the error message when there is no data to move for a specific volume group during a pvmove operation.
	NoDataToMovePattern = regexp.MustCompile(`No data to move for (.*?)`)

	// WarningNoFreeExtentsPattern is a regular expression that matches the error message when there are no free extents on a physical volume.
	WarningNoFreeExtentsPattern = regexp.MustCompile(`WARNING: No free extents on physical volume "(.*?)"`)
)

func isLVMError(err error, exitCode int, pattern *regexp.Regexp) bool {
	if err == nil {
		return false
	}
	lvmErr, ok := AsExitCodeError(err)
	if !ok || lvmErr.ExitCode() != exitCode {
		return false
	}
	return pattern.Match([]byte(lvmErr.Error()))
}

func IsLVMErrNotFound(err error) bool {
	return isLVMError(err, 5, NotFoundPattern)
}

func IsLVMErrNoSuchCommand(err error) bool {
	return isLVMError(err, 2, NoSuchCommandPattern)
}

func IsLVMErrMaximumLogicalVolumesReached(err error) bool {
	return isLVMError(err, 5, MaximumNumberOfLogicalVolumesPattern)
}

func IsLVMErrMaximumPhysicalVolumesReached(err error) bool {
	return isLVMError(err, 5, MaximumNumberOfPhysicalVolumesPattern)
}

func IsLVMErrVGImmutableDueToMissingPVs(err error) bool {
	return isLVMError(err, 5, CannotChangeVGWhilePVsAreMissingPattern)
}

func IsLVMWarningCouldNotFindDeviceWithUUID(err error) bool {
	return isLVMError(err, 5, WarningCouldNotFindDeviceWithUUIDPattern)
}

func IsLVMErrVGMissingPVs(err error) bool {
	return isLVMError(err, 5, VGMissingPVsPattern)
}

func LVMErrVGMissingPVsDetails(err error) (vg string, pv string, lastWrittenTo string, ok bool) {
	submatches := VGMissingPVsPattern.FindStringSubmatch(err.Error())
	if submatches == nil {
		return "", "", "", false
	}
	return submatches[1], submatches[2], submatches[3], true
}

func IsLVMWarningPartialLVNeedsRepairOrRemove(err error) bool {
	return isLVMError(err, 5, WarningPartialLVNeedsRepairOrRemovePattern)
}

func IsLVMErrThereAreStillPartialLVs(err error) bool {
	return isLVMError(err, 5, ThereAreStillPartialLVsPattern)
}

func IsLVMErrNoDataToMove(err error) bool {
	return isLVMError(err, 5, NoDataToMovePattern)
}

func IsLVMWarningNoFreeExtents(err error) bool {
	return isLVMError(err, 5, WarningNoFreeExtentsPattern)
}
