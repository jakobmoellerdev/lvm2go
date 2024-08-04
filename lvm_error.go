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
	"regexp"
	"slices"
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

	// CouldNotFindDeviceWithUUIDPattern is a regular expression that matches the error message when a device with a specific UUID is not found.
	CouldNotFindDeviceWithUUIDPattern = regexp.MustCompile(`Couldn't find device with uuid [\w-]+\.`)

	// VGMissingPVsPattern is a regular expression that matches the error message when a volume group is missing physical volumes.
	VGMissingPVsPattern = regexp.MustCompile(`VG (.*?) is missing PV (.*?) \(last written to (.*?)\)`)

	// ThereAreStillPartialLVsPattern is a regular expression that matches the error message when there are still partial logical volumes in a volume group.
	ThereAreStillPartialLVsPattern = regexp.MustCompile(`There are still partial LVs in VG (.*?)\.`)

	// PartialLVNeedsRepairOrRemovePattern is a regular expression that matches the error message when a logical volume needs repair or remove.
	PartialLVNeedsRepairOrRemovePattern = regexp.MustCompile(`Partial LV (.*?) needs to be repaired or removed\.`)

	// NoDataToMovePattern is a regular expression that matches the error message when there is no data to move for a specific volume group during a pvmove operation.
	NoDataToMovePattern = regexp.MustCompile(`No data to move for (.*?)`)

	// NoFreeExtentsPattern is a regular expression that matches the error message when there are no free extents on a physical volume.
	NoFreeExtentsPattern = regexp.MustCompile(`No free extents on physical volume "(.*?)"`)
)

func isLVMError(err error, pattern *regexp.Regexp, validExitCodes ...int) bool {
	if err == nil {
		return false
	}
	lvmErr, ok := AsExitCodeError(err)
	if !ok || !slices.Contains(validExitCodes, lvmErr.ExitCode()) {
		return false
	}
	return pattern.Match([]byte(lvmErr.Error()))
}

func IsLVMErrNotFound(err error) bool {
	return isLVMError(err, NotFoundPattern, 5)
}

func IsLVMErrNoSuchCommand(err error) bool {
	return isLVMError(err, NoSuchCommandPattern, 2)
}

func IsLVMErrMaximumLogicalVolumesReached(err error) bool {
	return isLVMError(err, MaximumNumberOfLogicalVolumesPattern, 5)
}

func IsLVMErrMaximumPhysicalVolumesReached(err error) bool {
	return isLVMError(err, MaximumNumberOfPhysicalVolumesPattern, 5)
}

func IsLVMErrVGImmutableDueToMissingPVs(err error) bool {
	return isLVMError(err, CannotChangeVGWhilePVsAreMissingPattern, 5)
}

func IsLVMCouldNotFindDeviceWithUUID(err error) bool {
	return isLVMError(err, CouldNotFindDeviceWithUUIDPattern, 5)
}

func IsLVMErrVGMissingPVs(err error) bool {
	return isLVMError(err, VGMissingPVsPattern, 5, 3)
}

func LVMErrVGMissingPVsDetails(err error) (vg string, pv string, lastWrittenTo string, ok bool) {
	submatches := VGMissingPVsPattern.FindStringSubmatch(err.Error())
	if submatches == nil {
		return "", "", "", false
	}
	return submatches[1], submatches[2], submatches[3], true
}

func IsLVMPartialLVNeedsRepairOrRemove(err error) bool {
	return isLVMError(err, PartialLVNeedsRepairOrRemovePattern, 5)
}

func IsLVMErrThereAreStillPartialLVs(err error) bool {
	return isLVMError(err, ThereAreStillPartialLVsPattern, 5)
}

func IsLVMErrNoDataToMove(err error) bool {
	return isLVMError(err, NoDataToMovePattern, 5)
}

func IsLVMNoFreeExtents(err error) bool {
	return isLVMError(err, NoFreeExtentsPattern, 5)
}
