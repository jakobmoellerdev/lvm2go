package lvm2go

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
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
	ExtentPercentFree ExtentPercent = "FREE"
	// ExtentPercentOrigin determines percentage of the total size of the origin LV
	ExtentPercentOrigin ExtentPercent = "ORIGIN"
	// ExtentPercentPVS determines percentage of the total size of the specified PVs
	ExtentPercentPVS ExtentPercent = "PVS"
	// ExtentPercentVG determines percentage of the total size of the VG
	ExtentPercentVG ExtentPercent = "VG"
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
		percent := extents[pidx+1:]
		if percent == "" || !slices.Contains(percentCandidates, ExtentPercent(percent)) {
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

	args.AddOrReplace("--extents", fmt.Sprintf("%s%s",
		strconv.FormatUint(opt.Val, 10),
		map[bool]string{
			true:  fmt.Sprintf("%s%s", ExtentPercentSymbol, opt.ExtentPercent),
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

	if !slices.Contains(prefixCandidates, opt.SizePrefix) {
		return ErrInvalidSizePrefix
	}

	return nil
}

func (opt PrefixedExtents) ApplyToArgs(args Arguments) error {
	if err := opt.Validate(); err != nil {
		return err
	}

	args.AddOrReplace("--extents", fmt.Sprintf("%s%s%s",
		string(opt.SizePrefix),
		strconv.FormatUint(opt.Val, 10),
		map[bool]string{
			true:  fmt.Sprintf("%s%s", ExtentPercentSymbol, opt.ExtentPercent),
			false: "",
		}[len(opt.ExtentPercent) > 0],
	))

	return nil
}
