package lvm2go

import (
	"errors"
	"fmt"
	"strings"
)

type VolumeType rune

var (
	ErrPartialActivation                         = errors.New("found partial activation of physical volumes, one or more physical volumes are setup incorrectly")
	ErrUnknownVolumeHealth                       = errors.New("unknown volume health reported, verification on the host system is required")
	ErrWriteCacheError                           = errors.New("write cache error signifies that dm-writecache reports an error")
	ErrThinPoolFailed                            = errors.New("thin pool encounters serious failures and hence no further I/O is permitted at all")
	ErrThinPoolOutOfDataSpace                    = errors.New("thin pool is out of data space, no further data can be written to the thin pool without extension")
	ErrThinPoolMetadataReadOnly                  = errors.New("metadata read only signifies that thin pool encounters certain types of failures, but it's still possible to do data reads. However, no metadata changes are allowed")
	ErrThinVolumeFailed                          = errors.New("the underlying thin pool entered a failed state and no further I/O is permitted")
	ErrRAIDRefreshNeeded                         = errors.New("RAID volume requires a refresh, one or more Physical Volumes have suffered a write error. This could be due to temporary failure of the Physical Volume or an indication it is failing. The device should be refreshed or replaced")
	ErrRAIDMismatchesExist                       = errors.New("RAID volume has portions of the array that are not coherent. Inconsistencies are detected by initiating a check RAID logical volume. The scrubbing operations, \"check\" and \"repair\", can be performed on a RAID volume via the \"lvchange\" command")
	ErrRAIDReshaping                             = errors.New("RAID volume is currently reshaping. Reshaping signifies a RAID Logical Volume is either undergoing a stripe addition/removal, a stripe size or RAID algorithm change")
	ErrRAIDReshapeRemoved                        = errors.New("RAID volume signifies freed raid images after reshaping")
	ErrRAIDWriteMostly                           = errors.New("RAID volume is marked as write-mostly. this signifies the devices in a RAID 1 logical volume have been marked write-mostly. This means that reading from this device will be avoided, and other devices will be preferred for reading (unless no other devices are available). This minimizes the I/O to the specified device")
	ErrLogicalVolumeSuspended                    = errors.New("logical volume is in a suspended state, no I/O is permitted")
	ErrInvalidSnapshot                           = errors.New("logical volume is an invalid snapshot, no I/O is permitted")
	ErrSnapshotMergeFailed                       = errors.New("snapshot merge failed, no I/O is permitted")
	ErrMappedDevicePresentWithInactiveTables     = errors.New("mapped device present with inactive tables, no I/O is permitted")
	ErrMappedDevicePresentWithoutTables          = errors.New("mapped device present without tables, no I/O is permitted")
	ErrThinPoolCheckNeeded                       = errors.New("a thin pool check is needed")
	ErrUnknownVolumeState                        = errors.New("unknown volume state, verification on the host system is required")
	ErrHistoricalVolumeState                     = errors.New("historical volume state (volume no longer exists but is kept around in logs), verification on the host system is required")
	ErrLogicalVolumeUnderlyingDeviceStateUnknown = errors.New("logical volume underlying device state is unknown, verification on the host system is required")
)

const (
	VolumeTypeMirrored                   VolumeType = 'm'
	VolumeTypeMirroredNoInitialSync      VolumeType = 'M'
	VolumeTypeOrigin                     VolumeType = 'o'
	VolumeTypeOriginWithMergingSnapshot  VolumeType = 'O'
	VolumeTypeRAID                       VolumeType = 'r'
	VolumeTypeRAIDNoInitialSync          VolumeType = 'R'
	VolumeTypeSnapshot                   VolumeType = 's'
	VolumeTypeMergingSnapshot            VolumeType = 'S'
	VolumeTypePVMove                     VolumeType = 'p'
	VolumeTypeVirtual                    VolumeType = 'v'
	VolumeTypeMirrorOrRAIDImage          VolumeType = 'i'
	VolumeTypeMirrorOrRAIDImageOutOfSync VolumeType = 'I'
	VolumeTypeMirrorLogDevice            VolumeType = 'l'
	VolumeTypeUnderConversion            VolumeType = 'c'
	VolumeTypeThinVolume                 VolumeType = 'V'
	VolumeTypeThinPool                   VolumeType = 't'
	VolumeTypeThinPoolData               VolumeType = 'T'
	VolumeTypeThinPoolMetadata           VolumeType = 'e'
	VolumeTypeNone                       VolumeType = '-'
)

type LVPermissions rune

const (
	LVPermissionsWriteable                             LVPermissions = 'w'
	LVPermissionsReadOnly                              LVPermissions = 'r'
	LVPermissionsReadOnlyActivationOfNonReadOnlyVolume LVPermissions = 'R'
	LVPermissionsNone                                  LVPermissions = '-'
)

type LVAllocationPolicyAttr rune

const (
	LVAllocationPolicyAttrAnywhere         LVAllocationPolicyAttr = 'a'
	LVAllocationPolicyAttrAnywhereLocked   LVAllocationPolicyAttr = 'A'
	LVAllocationPolicyAttrContiguous       LVAllocationPolicyAttr = 'c'
	LVAllocationPolicyAttrContiguousLocked LVAllocationPolicyAttr = 'C'
	LVAllocationPolicyAttrInherited        LVAllocationPolicyAttr = 'i'
	LVAllocationPolicyAttrInheritedLocked  LVAllocationPolicyAttr = 'I'
	LVAllocationPolicyAttrCling            LVAllocationPolicyAttr = 'l'
	LVAllocationPolicyAttrClingLocked      LVAllocationPolicyAttr = 'L'
	LVAllocationPolicyAttrNormal           LVAllocationPolicyAttr = 'n'
	LVAllocationPolicyAttrNormalLocked     LVAllocationPolicyAttr = 'N'
	LVAllocationPolicyAttrNone                                    = '-'
)

type Minor rune

const (
	MinorTrue  Minor = 'm'
	MinorFalse Minor = '-'
)

type State rune

const (
	StateActive                                State = 'a'
	StateSuspended                             State = 's'
	StateInvalidSnapshot                       State = 'I'
	StateSuspendedSnapshot                     State = 'S'
	StateSnapshotMergeFailed                   State = 'm'
	StateSuspendedSnapshotMergeFailed          State = 'M'
	StateMappedDevicePresentWithoutTables      State = 'd'
	StateMappedDevicePresentWithInactiveTables State = 'i'
	StateNone                                  State = '-'
	StateHistorical                            State = 'h'
	StateThinPoolCheckNeeded                   State = 'c'
	StateSuspendedThinPoolCheckNeeded          State = 'C'
	StateUnknown                               State = 'X'
)

type Open rune

const (
	OpenTrue    Open = 'o'
	OpenFalse   Open = '-'
	OpenUnknown Open = 'X'
)

type OpenTarget rune

const (
	OpenTargetMirror   = 'm'
	OpenTargetRaid     = 'r'
	OpenTargetSnapshot = 's'
	OpenTargetThin     = 't'
	OpenTargetUnknown  = 'u'
	OpenTargetVirtual  = 'v'
)

type ZeroAttr rune

const (
	ZeroAttrTrue  ZeroAttr = 'z'
	ZeroAttrFalse ZeroAttr = '-'
)

type VolumeHealth rune

const (
	VolumeHealthPartialActivation        = 'p'
	VolumeHealthUnknown                  = 'X'
	VolumeHealthOK                       = '-'
	VolumeHealthRAIDRefreshNeeded        = 'r'
	VolumeHealthRAIDMismatchesExist      = 'm'
	VolumeHealthRAIDWriteMostly          = 'w'
	VolumeHealthRAIDReshaping            = 's'
	VolumeHealthRAIDReshapeRemoved       = 'R'
	VolumeHealthThinFailed               = 'F'
	VolumeHealthThinPoolOutOfDataSpace   = 'D'
	VolumeHealthThinPoolMetadataReadOnly = 'M'
	VolumeHealthWriteCacheError          = 'E'
)

type SkipActivation rune

const (
	SkipActivationTrue  SkipActivation = 'k'
	SkipActivationFalse SkipActivation = '-'
)

// LVAttributes has mapped lv_attr information, see https://linux.die.net/man/8/lvs
// It is a complete parsing of the entire attribute byte flags that is attached to each LV.
// This is useful when attaching logic to the state of an LV as the state of an LV can be determined
// from the Attributes, e.g. for determining whether an LV is considered a Thin-Pool or not.
type LVAttributes struct {
	VolumeType
	LVPermissions
	LVAllocationPolicyAttr
	Minor
	State
	Open
	OpenTarget
	ZeroAttr
	VolumeHealth
	SkipActivation
}

func ParseLVAttributes(raw string) (LVAttributes, error) {
	if len(raw) != 10 {
		return LVAttributes{}, fmt.Errorf("%s is an invalid length lv_attr", raw)
	}
	return LVAttributes{
		VolumeType(raw[0]),
		LVPermissions(raw[1]),
		LVAllocationPolicyAttr(raw[2]),
		Minor(raw[3]),
		State(raw[4]),
		Open(raw[5]),
		OpenTarget(raw[6]),
		ZeroAttr(raw[7]),
		VolumeHealth(raw[8]),
		SkipActivation(raw[9]),
	}, nil
}

func (attr LVAttributes) String() string {
	var builder strings.Builder
	fields := []rune{
		rune(attr.VolumeType),
		rune(attr.LVPermissions),
		rune(attr.LVAllocationPolicyAttr),
		rune(attr.Minor),
		rune(attr.State),
		rune(attr.Open),
		rune(attr.OpenTarget),
		rune(attr.ZeroAttr),
		rune(attr.VolumeHealth),
		rune(attr.SkipActivation),
	}
	builder.Grow(len(fields))
	for _, r := range fields {
		builder.WriteRune(r)
	}
	return builder.String()
}

func (attr LVAttributes) MarshalText() ([]byte, error) {
	return []byte(attr.String()), nil
}

// VerifyHealth checks the health of the logical volume based on the attributes, mainly
// bit 9 (volume health indicator) based on bit 1 (volume type indicator)
// All failed known states are reported with an error message.
func (attr LVAttributes) VerifyHealth() error {
	if attr.VolumeHealth == VolumeHealthPartialActivation {
		return ErrPartialActivation
	}
	if attr.VolumeHealth == VolumeHealthUnknown {
		return ErrUnknownVolumeHealth
	}
	if attr.VolumeHealth == VolumeHealthWriteCacheError {
		return ErrWriteCacheError
	}

	if attr.VolumeType == VolumeTypeThinPool {
		switch attr.VolumeHealth {
		case VolumeHealthThinFailed:
			return ErrThinPoolFailed
		case VolumeHealthThinPoolOutOfDataSpace:
			return ErrThinPoolOutOfDataSpace
		case VolumeHealthThinPoolMetadataReadOnly:
			return ErrThinPoolMetadataReadOnly
		}
	}

	if attr.VolumeType == VolumeTypeThinVolume {
		switch attr.VolumeHealth {
		case VolumeHealthThinFailed:
			return ErrThinVolumeFailed
		}
	}

	if attr.VolumeType == VolumeTypeRAID || attr.VolumeType == VolumeTypeRAIDNoInitialSync {
		switch attr.VolumeHealth {
		case VolumeHealthRAIDRefreshNeeded:
			return ErrRAIDRefreshNeeded
		case VolumeHealthRAIDMismatchesExist:
			return ErrRAIDMismatchesExist
		case VolumeHealthRAIDReshaping:
			return ErrRAIDReshaping
		case VolumeHealthRAIDReshapeRemoved:
			return ErrRAIDReshapeRemoved
		case VolumeHealthRAIDWriteMostly:
			return ErrRAIDWriteMostly
		}
	}

	switch attr.State {
	case StateSuspended, StateSuspendedSnapshot:
		return ErrLogicalVolumeSuspended
	case StateInvalidSnapshot:
		return ErrInvalidSnapshot
	case StateSnapshotMergeFailed, StateSuspendedSnapshotMergeFailed:
		return ErrSnapshotMergeFailed
	case StateMappedDevicePresentWithInactiveTables:
		return ErrMappedDevicePresentWithInactiveTables
	case StateMappedDevicePresentWithoutTables:
		return ErrMappedDevicePresentWithoutTables
	case StateThinPoolCheckNeeded, StateSuspendedThinPoolCheckNeeded:
		return ErrThinPoolCheckNeeded
	case StateUnknown:
		return ErrUnknownVolumeState
	case StateHistorical:
		return ErrHistoricalVolumeState
	}

	switch attr.Open {
	case OpenUnknown:
		return ErrLogicalVolumeUnderlyingDeviceStateUnknown
	}

	return nil
}
