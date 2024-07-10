package lvm2go

import (
	"errors"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidSizeGTZero = errors.New("invalid size specified, must be set")

var ErrInvalidUnit = errors.New("invalid unit specified")
var ErrInvalidSizePrefix = errors.New("invalid size prefix specified")

type SizePrefix rune

const (
	SizePrefixMinus SizePrefix = '-'
	SizePrefixPlus  SizePrefix = '+'
)

var prefixCandidates = []SizePrefix{
	SizePrefixMinus,
	SizePrefixPlus,
}

type Unit rune

const (
	conversionFactor      = 1024
	UnitBytes        Unit = 'b'
	UnitKiB          Unit = 'k'
	UnitMiB          Unit = 'm'
	UnitGiB          Unit = 'g'
	UnitTiB          Unit = 't'
	UnitPiB          Unit = 'p'
	UnitEiB          Unit = 'e'
	UnitSector       Unit = 's'
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

// Size is an InputToParse number that accepts an optional unit.
// InputToParse units are always treated as base two values, regardless of capitalization, e.g.
// 'k' and 'K' both refer to 1024.
// The default InputToParse unit is specified by letter, followed by  |UNIT.
// UNIT represents other possible  InputToParse
// units: b is bytes, s is sectors of 512 bytes, k is KiB, m is MiB,
// g is GiB, t is TiB, p is PiB, e is EiB.
type Size struct {
	Val float64
	Unit
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

func convert(val float64, a, b Unit) float64 {
	if a == UnitUnknown || b == UnitUnknown {
		return val
	}

	if a == UnitSector {
		val *= 512
		a = UnitBytes
	}

	toSectorAtEnd := false
	if b == UnitSector {
		toSectorAtEnd = true
		b = UnitBytes
	}

	if conversionTable[a] < conversionTable[b] {
		val /= math.Pow(conversionFactor, conversionTable[b]-conversionTable[a])
	} else {
		val *= math.Pow(conversionFactor, conversionTable[a]-conversionTable[b])
	}

	if toSectorAtEnd {
		val /= 512
	}

	return val
}

func (opt Size) IsEqualTo(other Size) (bool, error) {
	if opt.Unit == other.Unit {
		return opt.Val == other.Val, nil
	}

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

func (opt Size) ToUnit(unit Unit) (Size, error) {
	if opt.Unit == unit {
		return opt, nil
	}

	if !IsValidUnit(unit) || opt.Unit == UnitUnknown {
		return Size{}, ErrInvalidUnit
	}

	return NewSize(convert(opt.Val, opt.Unit, unit), unit), nil
}

func (opt Size) unsafeToUnit(unit Unit) Size {
	if opt.Unit == unit {
		return opt
	}

	if !IsValidUnit(unit) || opt.Unit == UnitUnknown {
		return opt
	}

	return NewSize(convert(opt.Val, opt.Unit, unit), unit)
}

func (opt Size) String() string {
	var precision int
	if opt.Unit != UnitBytes {
		precision = 2
	}
	val := strconv.FormatFloat(opt.Val, 'f', precision, 64)
	if opt.Unit == UnitUnknown || opt.Unit == 0 {
		return val
	}
	return val + string(opt.Unit)
}

func MustParseSize(str string) Size {
	size, err := ParseSize(str)
	if err != nil {
		panic(err)
	}
	return size
}

func ParseSize(str string) (Size, error) {
	if len(str) == 0 {
		return Size{Unit: UnitUnknown}, nil
	}

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

	args.AddOrReplaceAll([]string{"--size", opt.String()})

	return nil
}

type PrefixedSize struct {
	SizePrefix
	Size
}

func MustParsePrefixedSize(str string) PrefixedSize {
	opt, err := ParsePrefixedSize(str)
	if err != nil {
		panic(err)
	}
	return opt
}

func ParsePrefixedSize(str string) (PrefixedSize, error) {
	if len(str) == 0 {
		size, err := ParseSize(str)
		if err != nil {
			return PrefixedSize{}, err
		}
		return PrefixedSize{Size: size}, nil
	}

	prefix := SizePrefix(str[0])
	if !slices.Contains(prefixCandidates, prefix) {
		return PrefixedSize{}, ErrInvalidSizePrefix
	}

	size, err := ParseSize(str[1:])
	if err != nil {
		return PrefixedSize{}, err
	}

	return NewPrefixedSize(prefix, size), nil
}

func NewPrefixedSize(prefix SizePrefix, size Size) PrefixedSize {
	return PrefixedSize{
		SizePrefix: prefix,
		Size:       size,
	}
}

func (opt PrefixedSize) Validate() error {
	if err := opt.Size.Validate(); err != nil {
		return err
	}

	if !slices.Contains(prefixCandidates, opt.SizePrefix) {
		return ErrInvalidSizePrefix
	}

	return nil
}

func (opt PrefixedSize) ApplyToArgs(args Arguments) error {
	if err := opt.Validate(); err != nil {
		return err
	}

	args.AddOrReplaceAll([]string{
		"--size",
		fmt.Sprintf("%s%s", string(opt.SizePrefix), opt.String()),
	})

	return nil
}
