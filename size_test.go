package lvm2go_test

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"

	. "github.com/jakobmoellerdev/lvm2go"
)

type SizeTestCase struct {
	InputToParse string
	expected     Size
	err          error
}

type PrefixedSizeTestCase struct {
	InputToParse string
	expected     PrefixedSize
	err          error
}

var DefaultSizeTestCases = []SizeTestCase{
	{"", Size{Unit: UnitUnknown}, nil},
	{"0", Size{Unit: UnitUnknown}, nil},
	{"1", Size{Val: 1, Unit: UnitUnknown}, nil},
	{"1B", Size{Val: 1, Unit: UnitBytes}, nil},
	{"1b", Size{Val: 1, Unit: UnitBytes}, nil},
	{"1K", Size{Val: 1, Unit: UnitKiB}, nil},
	{"1k", Size{Val: 1, Unit: UnitKiB}, nil},
	{"1M", Size{Val: 1, Unit: UnitMiB}, nil},
	{"1m", Size{Val: 1, Unit: UnitMiB}, nil},
	{"1G", Size{Val: 1, Unit: UnitGiB}, nil},
	{"1g", Size{Val: 1, Unit: UnitGiB}, nil},
	{"1T", Size{Val: 1, Unit: UnitTiB}, nil},
	{"1t", Size{Val: 1, Unit: UnitTiB}, nil},
	{"1P", Size{Val: 1, Unit: UnitPiB}, nil},
	{"1p", Size{Val: 1, Unit: UnitPiB}, nil},
	{"1E", Size{Val: 1, Unit: UnitEiB}, nil},
	{"1e", Size{Val: 1, Unit: UnitEiB}, nil},
	{"1S", Size{Val: 1, Unit: UnitSector}, nil},
	{"1s", Size{Val: 1, Unit: UnitSector}, nil},
	{"1%", Size{}, ErrInvalidUnit},
	{"xs", Size{}, strconv.ErrSyntax},
}

var PrefixedSizeTestCases = []PrefixedSizeTestCase{
	{"", PrefixedSize{Size: Size{Unit: UnitUnknown}}, nil},
	{".1B", PrefixedSize{}, ErrInvalidSizePrefix},
	{"+1%", PrefixedSize{}, ErrInvalidUnit},
}

func init() {
	for _, tc := range DefaultSizeTestCases {
		PrefixedSizeTestCases = append(
			PrefixedSizeTestCases,
			PrefixedSizeTestCase{
				InputToParse: fmt.Sprintf(
					"%s%s",
					string(SizePrefixMinus),
					tc.InputToParse,
				),
				expected: PrefixedSize{
					SizePrefix: SizePrefixMinus,
					Size:       tc.expected,
				},
				err: tc.err,
			},
			PrefixedSizeTestCase{
				InputToParse: fmt.Sprintf(
					"%s%s",
					string(SizePrefixPlus),
					tc.InputToParse,
				),
				expected: PrefixedSize{
					SizePrefix: SizePrefixPlus,
					Size:       tc.expected,
				},
				err: tc.err,
			},
		)
	}
}

func Test_Size(t *testing.T) {
	for _, tc := range DefaultSizeTestCases {
		t.Run(tc.InputToParse, func(t *testing.T) {
			actual, err := ParseSize(tc.InputToParse)
			if err != nil {
				if tc.err == nil {
					t.Errorf("unexpected error: %v", err)
				}
				if !errors.Is(err, tc.err) {
					t.Errorf("differing error: %v", err)
				}
			} else if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("unexpected size: %v (expected %v)", actual, tc.expected)
			}
		})

		t.Run("MustParse_"+tc.InputToParse, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if tc.err == nil {
						t.Errorf("unexpected panic: %v", r)
					}
				}
			}()
			actual := MustParseSize(tc.InputToParse)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("unexpected size: %v (expected %v)", actual, tc.expected)
			}
		})
	}

	for _, tc := range PrefixedSizeTestCases {
		t.Run(tc.InputToParse, func(t *testing.T) {
			actual, err := ParsePrefixedSize(tc.InputToParse)
			if err != nil {
				if tc.err == nil {
					t.Errorf("unexpected error: %v", err)
				}
				if !errors.Is(err, tc.err) {
					t.Errorf("differing error: %v", err)
				}
			} else if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("unexpected prefixed size: %v (expected %v)", actual, tc.expected)
			}
		})

		t.Run("MustParse_"+tc.InputToParse, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if tc.err == nil {
						t.Errorf("unexpected panic: %v", r)
					}
				}
			}()
			actual := MustParsePrefixedSize(tc.InputToParse)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("unexpected prefixed size: %v (expected %v)", actual, tc.expected)
			}
		})
	}

	t.Run("Equality", func(t *testing.T) {
		invalidSize := Size{Val: 1, Unit: '/'}
		for _, tc := range []struct {
			a, b Size
			eq   bool
			err  error
		}{
			{MustParseSize("1B"), MustParseSize("1B"), true, nil},
			{MustParseSize("1B"), MustParseSize("1b"), true, nil},
			{MustParseSize("1024B"), MustParseSize("1K"), true, nil},
			{MustParseSize("1K"), MustParseSize("1M"), false, nil},
			{invalidSize, MustParseSize("1M"), false, ErrInvalidUnit},
			{MustParseSize("1M"), invalidSize, false, ErrInvalidUnit},
		} {
			t.Run(fmt.Sprintf("%s==%s", tc.a, tc.b), func(t *testing.T) {
				if equals, err := tc.a.IsEqualTo(tc.b); err != nil {
					if tc.err == nil {
						t.Errorf("unexpected error: %v", err)
					}
					if !errors.Is(err, tc.err) {
						t.Errorf("differing error: %v", err)
					}
				} else if equals != tc.eq {
					t.Errorf("sizes should be equal")
				}
			})
		}
	})

	t.Run("Validation", func(t *testing.T) {
		for _, tc := range []struct {
			size Size
			err  error
		}{
			{Size{Val: 0, Unit: UnitUnknown}, ErrInvalidSizeGEZero},
			{Size{Val: 1, Unit: '/'}, ErrInvalidUnit},
			{Size{Val: 1, Unit: UnitUnknown}, nil},
		} {
			t.Run(tc.size.String(), func(t *testing.T) {
				if err := tc.size.Validate(); err != nil {
					if tc.err == nil {
						t.Errorf("unexpected error: %v", err)
					}
					if !errors.Is(err, tc.err) {
						t.Errorf("differing error: %v", err)
					}
				}
			})
		}

		for _, tc := range []struct {
			size PrefixedSize
			err  error
		}{
			{NewPrefixedSize(SizePrefixPlus, NewSize(1, UnitBytes)), nil},
			{NewPrefixedSize(SizePrefixMinus, NewSize(1, UnitBytes)), nil},
			{NewPrefixedSize(SizePrefixMinus, Size{Val: 1, Unit: '/'}), ErrInvalidUnit},
			{PrefixedSize{SizePrefix: 'x', Size: NewSize(1, UnitBytes)}, ErrInvalidSizePrefix},
		} {
			t.Run(tc.size.String(), func(t *testing.T) {
				if err := tc.size.Validate(); err != nil {
					if tc.err == nil {
						t.Errorf("unexpected error: %v", err)
					}
					if !errors.Is(err, tc.err) {
						t.Errorf("differing error: %v", err)
					}
				}
			})
		}
	})

	t.Run("convert", func(t *testing.T) {
		for _, tc := range []struct {
			val   float64
			a, b  Unit
			exp   float64
			error error
		}{
			{1, UnitKiB, UnitBytes, 1024, nil},
			{2, UnitKiB, UnitBytes, 2048, nil},
			{1, UnitMiB, UnitBytes, 1048576, nil},
			{1, UnitGiB, UnitBytes, 1073741824, nil},
			{1, UnitTiB, UnitBytes, 1099511627776, nil},
			{1, UnitPiB, UnitBytes, 1125899906842624, nil},
			{1, UnitEiB, UnitBytes, 1152921504606846976, nil},
			{1, UnitBytes, UnitKiB, 0.0009765625, nil},
			{1, UnitBytes, UnitMiB, 0.00000095367431640625, nil},
			{1, UnitBytes, UnitGiB, 0.0000000009313225746154785156, nil},
			{1, UnitBytes, UnitTiB, 0.0000000000009094947017729282379, nil},
			{1, UnitBytes, UnitPiB, 0.0000000000000008881784197001252323, nil},
			{1, UnitBytes, UnitEiB, 0.0000000000000000008673617379884035, nil},
			{1, UnitSector, UnitBytes, 512, nil},
			{1, UnitBytes, UnitSector, 0.001953125, nil},
			{1, UnitSector, UnitKiB, 0.5, nil},
			{2, UnitSector, UnitKiB, 1, nil},
			{1, UnitGiB, UnitKiB, 1048576, nil},
			{1, UnitGiB, UnitUnknown, 1, nil},
			{1, UnitUnknown, UnitGiB, -1, ErrInvalidUnit},
		} {
			t.Run(fmt.Sprintf("%.2f%s->%s", tc.val, tc.a, tc.b), func(t *testing.T) {
				a, b := NewSize(tc.val, tc.a), NewSize(tc.exp, tc.b)
				actual, err := a.ToUnit(tc.b)
				if !errors.Is(err, tc.error) {
					t.Errorf("unexpected error: %v", err)
				}

				if err == nil && !reflect.DeepEqual(actual, b) {
					t.Errorf("unexpected size: %v (expected %v)", actual, b)
				}
			})
		}
	})

	t.Run("ApplyToArgs", func(t *testing.T) {
		type applyToArgsTestCase struct {
			size     Size
			expected []string
			err      error
		}

		tcs := []applyToArgsTestCase{
			{Size{Val: -1, Unit: UnitUnknown}, nil, ErrInvalidSizeGEZero},
			{Size{Val: 0, Unit: UnitUnknown}, strings.Split("--size=0.00", " "), nil},
			{Size{Val: 1, Unit: UnitUnknown}, strings.Split("--size=1.00", " "), nil},
			{Size{Val: 1, Unit: UnitGiB}, strings.Split("--size=1.00g", " "), nil},
			{Size{Val: 1, Unit: UnitBytes}, strings.Split("--size=1b", " "), nil},
			{Size{Val: 1.555, Unit: UnitGiB}, strings.Split("--size=1.55g", " "), nil},
		}

		for _, tc := range tcs {
			t.Run(tc.size.String(), func(t *testing.T) {
				args := NewArgs(ArgsTypeGeneric)
				if err := tc.size.ApplyToArgs(args); err != nil {
					if tc.err == nil {
						t.Errorf("unexpected error: %v", err)
					}
					if !errors.Is(err, tc.err) {
						t.Errorf("differing error: %v", err)
					}
				} else if !reflect.DeepEqual(args.GetRaw(), tc.expected) {
					t.Errorf("unexpected args: %v", args.GetRaw())
				}
			})
		}

		for _, tc := range []struct {
			size     PrefixedSize
			expected []string
			err      error
		}{
			{PrefixedSize{SizePrefix: SizePrefixPlus, Size: Size{Val: 0, Unit: UnitUnknown}}, nil, ErrInvalidSizeGEZero},
			{PrefixedSize{SizePrefix: SizePrefixPlus, Size: Size{Val: 1, Unit: UnitUnknown}}, strings.Split("--size=+1.00", " "), nil},
			{PrefixedSize{SizePrefix: SizePrefixPlus, Size: Size{Val: 1, Unit: UnitBytes}}, strings.Split("--size=+1b", " "), nil},
			{PrefixedSize{SizePrefix: SizePrefixPlus, Size: Size{Val: 1, Unit: UnitGiB}}, strings.Split("--size=+1.00g", " "), nil},
			{PrefixedSize{SizePrefix: SizePrefixMinus, Size: Size{Val: 1, Unit: UnitGiB}}, strings.Split("--size=-1.00g", " "), nil},
		} {
			t.Run(tc.size.String(), func(t *testing.T) {
				args := NewArgs(ArgsTypeGeneric)
				if err := tc.size.ApplyToArgs(args); err != nil {
					if tc.err == nil {
						t.Errorf("unexpected error: %v", err)
					}
					if !errors.Is(err, tc.err) {
						t.Errorf("differing error: %v", err)
					}
				} else if !reflect.DeepEqual(args.GetRaw(), tc.expected) {
					t.Errorf("unexpected args: %v", args.GetRaw())
				}
			})
		}
	})
}
