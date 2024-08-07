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
	"fmt"
	"regexp"
)

var (

	// NotFoundPatterns are regular expressions that matches the error message when a device, volume group or logical volume is not found.

	volumeGroupNotFoundPattern   = `Volume group "(.*?)" not found`
	VolumeGroupNotFoundPattern   = regexp.MustCompile(volumeGroupNotFoundPattern)
	logicalVolumeNotFoundPattern = `Failed to find logical volume "(.*?)"`
	LogicalVolumeNotFoundPattern = regexp.MustCompile(logicalVolumeNotFoundPattern)
	deviceNotFoundPattern        = `Couldn't find device with uuid (.{6}-.{4}-.{4}-.{4}-.{4}-.{4}-.{6})`
	DeviceNotFoundPattern        = regexp.MustCompile(deviceNotFoundPattern)
	NotFoundPattern              = regexp.MustCompile(fmt.Sprintf(`%s|%s|%s`, volumeGroupNotFoundPattern, logicalVolumeNotFoundPattern, deviceNotFoundPattern))

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

	ConfigurationSectionNotCustomizableByProfilePattern = regexp.MustCompile(`Configuration section "(.*?)" is not customizable by a profile\.`)
)

// IsLVMError returns true if the error is an LVM error with a specific exit code and matches a specific pattern.
// The validExitCodes are the exit codes that are considered valid for the error.
// While lvm2go packages a lot of predefined patterns, it is possible to use a custom pattern.
//
// Example:
//
//	func IsLVMCustomError(err error) bool {
//		return IsLVMError(err, regexp.MustCompile(`custom error pattern`))
//	}
func IsLVMError(err error, pattern *regexp.Regexp) bool {
	if err == nil {
		return false
	}

	if stdErr, ok := AsLVMStdErr(err); ok {
		for _, line := range stdErr.Lines(true) {
			if pattern.Match(line) {
				return true
			}
		}
	}

	return false
}

func IsNotFound(err error) bool {
	return IsLVMError(err, NotFoundPattern)
}

func IsVolumeGroupNotFound(err error) bool {
	return IsLVMError(err, VolumeGroupNotFoundPattern)
}

func IsLogicalVolumeNotFound(err error) bool {
	return IsLVMError(err, LogicalVolumeNotFoundPattern)
}

func IsDeviceNotFound(err error) bool {
	return IsLVMError(err, DeviceNotFoundPattern)
}

func IsNoSuchCommand(err error) bool {
	return IsLVMError(err, NoSuchCommandPattern)
}

func IsMaximumLogicalVolumesReached(err error) bool {
	return IsLVMError(err, MaximumNumberOfLogicalVolumesPattern)
}

func IsMaximumPhysicalVolumesReached(err error) bool {
	return IsLVMError(err, MaximumNumberOfPhysicalVolumesPattern)
}

func IsVGImmutableDueToMissingPVs(err error) bool {
	return IsLVMError(err, CannotChangeVGWhilePVsAreMissingPattern)
}

func IsCouldNotFindDeviceWithUUID(err error) bool {
	return IsLVMError(err, CouldNotFindDeviceWithUUIDPattern)
}

func IsVGMissingPVs(err error) bool {
	return IsLVMError(err, VGMissingPVsPattern)
}

func VGMissingPVsDetails(err error) (vg string, pv string, lastWrittenTo string, ok bool) {
	submatches := VGMissingPVsPattern.FindStringSubmatch(err.Error())
	if submatches == nil {
		return "", "", "", false
	}
	return submatches[1], submatches[2], submatches[3], true
}

func IsPartialLVNeedsRepairOrRemove(err error) bool {
	return IsLVMError(err, PartialLVNeedsRepairOrRemovePattern)
}

func IsThereAreStillPartialLVs(err error) bool {
	return IsLVMError(err, ThereAreStillPartialLVsPattern)
}

func IsNoDataToMove(err error) bool {
	return IsLVMError(err, NoDataToMovePattern)
}

func IsNoFreeExtents(err error) bool {
	return IsLVMError(err, NoFreeExtentsPattern)
}

func IsConfigurationSectionNotCustomizableByProfile(err error) bool {
	return IsLVMError(err, ConfigurationSectionNotCustomizableByProfilePattern)
}
