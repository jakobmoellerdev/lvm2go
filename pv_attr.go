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
