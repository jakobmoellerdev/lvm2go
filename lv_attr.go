package lvm2go

import (
	"fmt"
	"strings"
)

type VolumeType rune

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

type Permissions rune

const (
	PermissionsWriteable                             Permissions = 'w'
	PermissionsReadOnly                              Permissions = 'r'
	PermissionsReadOnlyActivationOfNonReadOnlyVolume Permissions = 'R'
	PermissionsNone                                  Permissions = '-'
)

type AllocationPolicyAttr rune

const (
	AllocationPolicyAttrAnywhere         AllocationPolicyAttr = 'a'
	AllocationPolicyAttrAnywhereLocked   AllocationPolicyAttr = 'A'
	AllocationPolicyAttrContiguous       AllocationPolicyAttr = 'c'
	AllocationPolicyAttrContiguousLocked AllocationPolicyAttr = 'C'
	AllocationPolicyAttrInherited        AllocationPolicyAttr = 'i'
	AllocationPolicyAttrInheritedLocked  AllocationPolicyAttr = 'I'
	AllocationPolicyAttrCling            AllocationPolicyAttr = 'l'
	AllocationPolicyAttrClingLocked      AllocationPolicyAttr = 'L'
	AllocationPolicyAttrNormal           AllocationPolicyAttr = 'n'
	AllocationPolicyAttrNormalLocked     AllocationPolicyAttr = 'N'
	AllocationPolicyAttrNone                                  = '-'
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
	StateUnknown                               State = 'X'
	StateCheckNeeded                           State = 'c'
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

type PartialAttr rune

const (
	PartialAttrTrue  = 'p'
	PartialAttrFalse = '-'
)

// LVAttributes has mapped lv_attr information, see https://linux.die.net/man/8/lvs
// It is a complete parsing of the entire attribute byte flags that is attached to each LV.
// This is useful when attaching logic to the state of an LV as the state of an LV can be determined
// from the Attributes, e.g. for determining whether an LV is considered a Thin-Pool or not.
type LVAttributes struct {
	VolumeType
	Permissions
	AllocationPolicyAttr
	Minor
	State
	Open
	OpenTarget
	ZeroAttr
	PartialAttr
}

func ParsedLVAttributes(raw string) (LVAttributes, error) {
	if len(raw) != 10 {
		return LVAttributes{}, fmt.Errorf("%s is an invalid length lv_attr", raw)
	}
	return LVAttributes{
		VolumeType(raw[0]),
		Permissions(raw[1]),
		AllocationPolicyAttr(raw[2]),
		Minor(raw[3]),
		State(raw[4]),
		Open(raw[5]),
		OpenTarget(raw[6]),
		ZeroAttr(raw[7]),
		PartialAttr(raw[8]),
	}, nil
}

func (l LVAttributes) String() string {
	var builder strings.Builder
	builder.Grow(9)
	builder.WriteRune(rune(l.VolumeType))
	builder.WriteRune(rune(l.Permissions))
	builder.WriteRune(rune(l.AllocationPolicyAttr))
	builder.WriteRune(rune(l.Minor))
	builder.WriteRune(rune(l.State))
	builder.WriteRune(rune(l.Open))
	builder.WriteRune(rune(l.OpenTarget))
	builder.WriteRune(rune(l.ZeroAttr))
	builder.WriteRune(rune(l.PartialAttr))
	return builder.String()
}

func (l LVAttributes) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}
