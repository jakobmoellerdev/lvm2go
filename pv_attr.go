package lvm2go

import (
	"fmt"
	"strings"
)

type DuplicateAllocatableUsed rune

const (
	Duplicate   DuplicateAllocatableUsed = 'd'
	Allocatable DuplicateAllocatableUsed = 'a'
	Used        DuplicateAllocatableUsed = 'u'
)

type Missing rune

const (
	MissingTrue  Missing = 'm'
	MissingFalse Missing = '-'
)

type PVAttributes struct {
	DuplicateAllocatableUsed
	Exported
	Missing
}

func ParsePVAttributes(raw string) (PVAttributes, error) {
	if len(raw) != 3 {
		return PVAttributes{}, fmt.Errorf("%s is an invalid length vg_attr", raw)
	}
	return PVAttributes{
		DuplicateAllocatableUsed: DuplicateAllocatableUsed(raw[0]),
		Exported:                 Exported(raw[1]),
		Missing:                  Missing(raw[2]),
	}, nil
}

func (attr PVAttributes) String() string {
	var builder strings.Builder
	fields := []rune{
		rune(attr.DuplicateAllocatableUsed),
		rune(attr.Exported),
		rune(attr.Missing),
	}
	builder.Grow(len(fields))
	for _, r := range fields {
		builder.WriteRune(r)
	}
	return builder.String()
}

func (attr PVAttributes) MarshalText() ([]byte, error) {
	return []byte(attr.String()), nil
}
