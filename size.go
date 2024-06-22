package lvm2go

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidSizeGTZero = errors.New("invalid size specified, must be set")

var ErrInvalidUnit = errors.New("invalid unit specified")
var ErrCannotConvertSector = errors.New("cannot convert sector to other units")

type Unit rune

const (
	UnitBytes  Unit = 'b'
	UnitKiB    Unit = 'k'
	UnitMiB    Unit = 'm'
	UnitGiB    Unit = 'g'
	UnitTiB    Unit = 't'
	UnitPiB    Unit = 'p'
	UnitEiB    Unit = 'e'
	UnitSector Unit = 's'
	// UnitUnknown is used to represent the output unit when
	// LVs or VGs are queried without specifying a unit. (--nosuffix)
	UnitUnknown Unit = 'X'
)

var validUnits = []Unit{UnitBytes, UnitKiB, UnitMiB, UnitGiB, UnitTiB, UnitPiB, UnitEiB, UnitSector, UnitUnknown}

func IsValidUnit(unit Unit) bool {
	for _, valid := range validUnits {
		if valid == unit || strings.ToUpper(string(valid)) == string(unit) {
			return true
		}
	}
	if unicode.IsDigit(rune(unit)) {
		return true
	}
	return false
}

// Size is an input number that accepts an optional unit.
// Input units are always treated as base two values, regardless of capitalization, e.g.
// 'k' and 'K' both refer to 1024.
// The default input unit is specified by letter, followed by  |UNIT.
// UNIT represents other possible  input
// units: b is bytes, s is sectors of 512 bytes, k is KiB, m is MiB,
// g is GiB, t is TiB, p is PiB, e is EiB.
type Size struct {
	Val float64
	Unit
}

func (opt Size) IsEqualTo(other Size) (bool, error) {
	optBytes, err := opt.ToUnit(UnitBytes)
	if err != nil {
		return false, err
	}

	otherBytes, err := other.ToUnit(UnitBytes)
	if err != nil {
		return false, err
	}

	return optBytes == otherBytes, nil
}

var conversionTable = map[Unit]float64{
	UnitBytes: 0,
	UnitKiB:   1,
	UnitMiB:   2,
	UnitGiB:   3,
	UnitTiB:   4,
	UnitPiB:   5,
	UnitEiB:   6,
}

func (opt Size) ToUnit(unit Unit) (Size, error) {
	if opt.Unit == unit {
		return opt, nil
	}

	if !IsValidUnit(unit) || opt.Unit == UnitUnknown {
		return Size{}, ErrInvalidUnit
	}

	if opt.Unit == UnitSector || unit == UnitSector {
		return Size{}, ErrCannotConvertSector
	}

	var factor float64
	if conversionTable[opt.Unit] < conversionTable[unit] {
		factor = conversionTable[unit] - conversionTable[opt.Unit]
	} else {
		factor = conversionTable[opt.Unit] - conversionTable[unit]
	}

	return NewSize(opt.Val*math.Pow(1024, factor), unit), nil
}

func (opt Size) String() string {
	if opt.Unit == UnitUnknown {
		return strconv.FormatFloat(opt.Val, 'f', 2, 64)
	}
	return strconv.FormatFloat(opt.Val, 'f', 2, 64) + string(opt.Unit)
}

func MustParseSize(str string) Size {
	size, err := ParseSize(str)
	if err != nil {
		panic(err)
	}
	return size
}

func ParseSize(str string) (Size, error) {
	var unit Unit
	offset := 0
	if len(str) > 1 && !unicode.IsDigit(rune(str[len(str)-1])) {
		unit = Unit(unicode.ToLower(rune(str[len(str)-1])))
		offset++
	} else {
		unit = UnitUnknown
	}

	if !IsValidUnit(unit) {
		return Size{}, ErrInvalidUnit
	}

	fval, err := strconv.ParseFloat(str[:len(str)-offset], 64)
	if err != nil {
		return Size{}, err
	}

	return NewSize(fval, unit), nil
}

func NewSize(value float64, unit Unit) Size {
	return Size{
		Val:  value,
		Unit: unit,
	}
}

func (opt Size) Validate() error {
	if opt.Val <= 0 {
		return ErrInvalidSizeGTZero
	}

	if !IsValidUnit(opt.Unit) {
		return ErrInvalidUnit
	}

	return nil
}

func (opt Size) ApplyToLVCreateOptions(opts *LVCreateOptions) {
	opts.Size = opt
}

func (opt Size) ApplyToArgs(args Arguments) error {
	if err := opt.Validate(); err != nil {
		return err
	}

	args.AppendAll([]string{"--size", opt.String()})

	return nil
}
