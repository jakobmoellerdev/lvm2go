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
	"errors"
	"fmt"
	"reflect"
	"testing"
)

type ExtentTestCase struct {
	InputToParse string
	expected     Extents
	err          error
}

type PrefixedExtentTestCase struct {
	InputToParse string
	expected     PrefixedExtents
	err          error
}

var DefaultExtentTestCases = []ExtentTestCase{
	{"", Extents{}, nil},
	{"0", Extents{}, nil},
	{"1", Extents{Val: 1}, nil},
	{"1%FREE", Extents{Val: 1, ExtentPercent: ExtentPercentFree}, nil},
	{"1%VG", Extents{Val: 1, ExtentPercent: ExtentPercentVG}, nil},
	{"1%PVS", Extents{Val: 1, ExtentPercent: ExtentPercentPVS}, nil},
	{"%PVS", Extents{}, ErrInvalidCannotStartWithPercent},
	{"1%PVS%", Extents{}, ErrInvalidMultiplePercent},
	{"1%BLA", Extents{}, ErrInvalidPercentDefinition},
	{"x%FREE", Extents{}, ErrNANExtents},
}

var PrefixedExtentTestCases = []PrefixedExtentTestCase{
	{".1%FREE", PrefixedExtents{}, ErrInvalidSizePrefix},
}

func init() {
	for _, tc := range DefaultExtentTestCases {
		PrefixedExtentTestCases = append(
			PrefixedExtentTestCases,
			PrefixedExtentTestCase{
				InputToParse: fmt.Sprintf(
					"%s%s",
					string(SizePrefixMinus),
					tc.InputToParse,
				),
				expected: PrefixedExtents{
					SizePrefix: SizePrefixMinus,
					Extents:    tc.expected,
				},
				err: tc.err,
			},
			PrefixedExtentTestCase{
				InputToParse: fmt.Sprintf(
					"%s%s",
					string(SizePrefixPlus),
					tc.InputToParse,
				),
				expected: PrefixedExtents{
					SizePrefix: SizePrefixPlus,
					Extents:    tc.expected,
				},
				err: tc.err,
			},
		)
	}
}

func Test_Extents(t *testing.T) {
	t.Parallel()
	for _, tc := range DefaultExtentTestCases {
		t.Run(tc.InputToParse, func(t *testing.T) {
			actual, err := ParseExtents(tc.InputToParse)
			if err != nil {
				if tc.err == nil {
					t.Errorf("unexpected error: %v", err)
				}
				if !errors.Is(err, tc.err) {
					t.Errorf("differing error: %v", err)
				}
			} else if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("unexpected extents: %v", actual)
			}
		})

		t.Run("MustParse_"+tc.InputToParse, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if tc.err == nil {
						t.Errorf("unexpected panic: %v", r)
					}
					if !errors.Is(r.(error), tc.err) {
						t.Errorf("differing panic: %v", r)
					}
				}
			}()
			actual := MustParseExtents(tc.InputToParse)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("unexpected extents: %v", actual)
			}
		})
	}

	for _, tc := range PrefixedExtentTestCases {
		t.Run(tc.InputToParse, func(t *testing.T) {
			actual, err := ParsePrefixedExtents(tc.InputToParse)
			if err != nil {
				if tc.err == nil {
					t.Errorf("unexpected error: %v", err)
				}
				if !errors.Is(err, tc.err) {
					t.Errorf("differing error: %v", err)
				}
			} else if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("unexpected extents: %v", actual)
			}
		})

		t.Run("MustParse_"+tc.InputToParse, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if tc.err == nil {
						t.Errorf("unexpected panic: %v", r)
					}
					if !errors.Is(r.(error), tc.err) {
						t.Errorf("differing panic: %v", r)
					}
				}
			}()
			actual := MustParsePrefixedExtents(tc.InputToParse)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("unexpected extents: %v", actual)
			}
		})
	}

	t.Run("NewExtents", func(t *testing.T) {
		actual := NewExtents(1, ExtentPercentOrigin)
		expected := Extents{Val: 1, ExtentPercent: ExtentPercentOrigin}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("unexpected extents: %v", actual)
		}
		if err := actual.Validate(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("NewPrefixedExtents", func(t *testing.T) {
		actual := NewPrefixedExtents(SizePrefixPlus, NewExtents(1, ExtentPercentOrigin))
		expected := PrefixedExtents{SizePrefix: SizePrefixPlus, Extents: NewExtents(1, ExtentPercentOrigin)}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("unexpected extents: %v", actual)
		}
		if err := actual.Validate(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("ApplyToArgs", func(t *testing.T) {
		args := NewArgs(ArgsTypeLVCreate)
		extents := NewExtents(1, ExtentPercentOrigin)
		if err := extents.ApplyToArgs(args); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if act, exp := args.GetRaw(), "--extents=1%ORIGIN"; len(act) != 1 || act[0] != exp {
			t.Errorf("unexpected args: expected %s but got %s(%v)", exp, act, len(act))
		}
	})

	t.Run("PrefixedApplyToArgs", func(t *testing.T) {
		args := NewArgs(ArgsTypeLVCreate)
		extents := NewPrefixedExtents(SizePrefixPlus, NewExtents(1, ExtentPercentOrigin))
		if err := extents.ApplyToArgs(args); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if act, exp := args.GetRaw(), "--extents=+1%ORIGIN"; len(act) != 1 || act[0] != exp {
			t.Errorf("unexpected args: expected %s but got %s(%v)", exp, act, len(act))
		}
	})

	t.Run("Validate", func(t *testing.T) {
		if err := NewExtents(0, ExtentPercentOrigin).Validate(); !errors.Is(err, ErrInvalidExtentsGTZero) {
			t.Errorf("unexpected error: %v", err)
		}
		if err := NewPrefixedExtents(SizePrefix('x'), NewExtents(1, ExtentPercentOrigin)).Validate(); !errors.Is(err, ErrInvalidSizePrefix) {
			t.Errorf("unexpected error: %v", err)
		}
		if err := NewPrefixedExtents(SizePrefixPlus, NewExtents(0, ExtentPercentOrigin)).Validate(); !errors.Is(err, ErrInvalidExtentsGTZero) {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
