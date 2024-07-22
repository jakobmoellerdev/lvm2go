package lvm2go

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidExtentsGTZero = errors.New("invalid extents specified, must be set")
var ErrInvalidCannotStartWithPercent = fmt.Errorf("invalid extents, cannot start with %q", ExtentPercentSymbol)
var ErrInvalidMultiplePercent = fmt.Errorf("multiple %q found", ExtentPercentSymbol)
var ErrInvalidPercentDefinition = fmt.Errorf("invalid percent definition, must be one of %v", percentCandidates)
var ErrNANExtents = errors.New("invalid extents specified, must be a valid integer number")

const ExtentPercentSymbol = "%"

type ExtentPercent string

const (
	// ExtentPercentFree determines percentage of remaining free space in the VG
	ExtentPercentFree ExtentPercent = ExtentPercentSymbol + "FREE"
	// ExtentPercentOrigin determines percentage of the total size of the origin LV
	ExtentPercentOrigin ExtentPercent = ExtentPercentSymbol + "ORIGIN"
	// ExtentPercentPVS determines percentage of the total size of the specified PVs
	ExtentPercentPVS ExtentPercent = ExtentPercentSymbol + "PVS"
	// ExtentPercentVG determines percentage of the total size of the VG
	ExtentPercentVG ExtentPercent = ExtentPercentSymbol + "VG"
)

var percentCandidates = []ExtentPercent{
	ExtentPercentFree,
	ExtentPercentOrigin,
	ExtentPercentPVS,
	ExtentPercentVG,
}

type Extents struct {
	Val uint64
	ExtentPercent
}

func NewExtents(val uint64, percent ExtentPercent) Extents {
	return Extents{
		Val:           val,
		ExtentPercent: percent,
	}
}

func MustParseExtents(extents string) Extents {
	e, err := ParseExtents(extents)
	if err != nil {
		panic(err)
	}
	return e
}

func ParseExtents(extents string) (Extents, error) {
	if extents == "" {
		return Extents{}, nil
	}

	e := Extents{}

	pidx := strings.Index(extents, ExtentPercentSymbol)

	if pidx == 0 {
		return e, ErrInvalidCannotStartWithPercent
	}

	lpidx := strings.LastIndex(extents, ExtentPercentSymbol)
	if pidx != lpidx {
		return e, ErrInvalidMultiplePercent
	}
	if pidx > 0 {
		percent := extents[pidx:]
		if !slices.Contains(percentCandidates, ExtentPercent(percent)) {
			return e, ErrInvalidPercentDefinition
		}

		e.ExtentPercent = ExtentPercent(percent)
		extents = extents[:pidx]
	}

	extentValue, err := strconv.ParseUint(extents, 10, 64)
	if err != nil {
		return e, fmt.Errorf("%q cannot be used: %w",
			extents, errors.Join(ErrNANExtents, err))
	}

	e.Val = extentValue

	return e, nil
}

func (opt Extents) ApplyToArgs(args Arguments) error {
	if err := opt.Validate(); err != nil {
		return err
	}
	if opt.Val == 0 {
		return nil
	}

	args.AddOrReplace(fmt.Sprintf("--extents=%s%s",
		strconv.FormatUint(opt.Val, 10),
		map[bool]string{
			true:  string(opt.ExtentPercent),
			false: "",
		}[len(opt.ExtentPercent) > 0],
	))
	return nil
}

func (opt Extents) Validate() error {
	if opt.Val <= 0 {
		return ErrInvalidExtentsGTZero
	}

	return nil
}

func (opt Extents) ToSize(extentSize uint64) Size {
	return NewSize(float64(opt.Val*extentSize), UnitBytes)
}

func (opt Extents) ApplyToLVCreateOptions(opts *LVCreateOptions) {
	opts.Extents = opt
}

type PrefixedExtents struct {
	SizePrefix
	Extents
}

func NewPrefixedExtents(prefix SizePrefix, extents Extents) PrefixedExtents {
	return PrefixedExtents{
		SizePrefix: prefix,
		Extents:    extents,
	}
}

func MustParsePrefixedExtents(str string) PrefixedExtents {
	opt, err := ParsePrefixedExtents(str)
	if err != nil {
		panic(err)
	}
	return opt
}

func ParsePrefixedExtents(str string) (PrefixedExtents, error) {
	prefix := SizePrefix(str[0])

	if unicode.IsDigit(rune(prefix)) {
		extents, err := ParseExtents(str)
		if err != nil {
			return PrefixedExtents{}, err
		}
		return NewPrefixedExtents(SizePrefixNone, extents), nil
	}

	if !slices.Contains(prefixCandidates, prefix) {
		return PrefixedExtents{}, ErrInvalidSizePrefix
	}

	extents, err := ParseExtents(str[1:])
	if err != nil {
		return PrefixedExtents{}, err
	}

	return NewPrefixedExtents(prefix, extents), nil
}

func (opt PrefixedExtents) Validate() error {
	if err := opt.Extents.Validate(); err != nil {
		return err
	}

	if opt.SizePrefix != SizePrefixNone && !slices.Contains(prefixCandidates, opt.SizePrefix) {
		return ErrInvalidSizePrefix
	}

	return nil
}

func (opt PrefixedExtents) ApplyToArgs(args Arguments) error {
	if err := opt.Validate(); err != nil {
		return err
	}
	if opt.Val == 0 {
		return nil
	}

	args.AddOrReplace(fmt.Sprintf("--extents=%s%s%s",
		map[bool]string{
			true:  string(opt.SizePrefix),
			false: "",
		}[opt.SizePrefix != SizePrefixNone],
		strconv.FormatUint(opt.Val, 10),
		map[bool]string{
			true:  string(opt.ExtentPercent),
			false: "",
		}[len(opt.ExtentPercent) > 0],
	))

	return nil
}

func (opt PrefixedExtents) ApplyToLVExtendOptions(opts *LVExtendOptions) {
	opts.PrefixedExtents = opt
}
