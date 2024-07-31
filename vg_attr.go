package lvm2go

import (
	"fmt"
	"strings"
)

type VGPermissions rune

const (
	VGPermissionsWriteable VGPermissions = 'w'
	VGPermissionsReadOnly  VGPermissions = 'r'
	VGPermissionsNone      VGPermissions = '-'
)

type Resizeable rune

const (
	ResizeableTrue  Resizeable = 'z'
	ResizeableFalse Resizeable = '-'
)

type Exported rune

const (
	ExportedTrue  Exported = 'x'
	ExportedFalse Exported = '-'
)

type PartialAttr rune

const (
	PartialAttrTrue  PartialAttr = 'p'
	PartialAttrFalse PartialAttr = '-'
)

type VGAllocationPolicyAttr rune

const (
	VGAllocationPolicyAttrAnywhere   VGAllocationPolicyAttr = 'a'
	VGAllocationPolicyAttrContiguous VGAllocationPolicyAttr = 'c'
	VGAllocationPolicyAttrCling      VGAllocationPolicyAttr = 'l'
	VGAllocationPolicyAttrNormal     VGAllocationPolicyAttr = 'n'
	VGAllocationPolicyAttrNone       VGAllocationPolicyAttr = '-'
)

type ClusteredOrShared rune

const (
	ClusteredOrSharedTrue  ClusteredOrShared = 'c'
	ClusteredOrSharedFalse ClusteredOrShared = '-'
)

type VGAttributes struct {
	VGPermissions
	Resizeable
	Exported
	PartialAttr
	VGAllocationPolicyAttr
	ClusteredOrShared
}

func ParseVGAttributes(raw string) (VGAttributes, error) {
	if len(raw) != 6 {
		return VGAttributes{}, fmt.Errorf("%s is an invalid length vg_attr", raw)
	}
	return VGAttributes{
		VGPermissions:          VGPermissions(raw[0]),
		Resizeable:             Resizeable(raw[1]),
		Exported:               Exported(raw[2]),
		PartialAttr:            PartialAttr(raw[3]),
		VGAllocationPolicyAttr: VGAllocationPolicyAttr(raw[4]),
		ClusteredOrShared:      ClusteredOrShared(raw[5]),
	}, nil
}

func (attr VGAttributes) String() string {
	var builder strings.Builder
	fields := []rune{
		rune(attr.VGPermissions),
		rune(attr.Resizeable),
		rune(attr.Exported),
		rune(attr.PartialAttr),
		rune(attr.VGAllocationPolicyAttr),
		rune(attr.ClusteredOrShared),
	}
	builder.Grow(len(fields))
	for _, r := range fields {
		builder.WriteRune(r)
	}
	return builder.String()
}

func (attr VGAttributes) MarshalText() ([]byte, error) {
	return []byte(attr.String()), nil
}
