package lvm2go

import (
	"errors"
	"fmt"
	"log/slog"
	"math"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidSizeGEZero = errors.New("invalid size specified, must be greater than or equal to zero")

var ErrInvalidUnit = errors.New("invalid unit specified")
var ErrInvalidSizePrefix = fmt.Errorf("invalid size prefix specified, must be one of %v", prefixCandidates)

type SizePrefix rune

const (
	SizePrefixNone  SizePrefix = 0
	SizePrefixMinus SizePrefix = '-'
	SizePrefixPlus  SizePrefix = '+'
)

var prefixCandidates = []SizePrefix{
	SizePrefixMinus,
	SizePrefixPlus,
}

const (
	sizeArg                = "--size"
	metadataSizeArg        = "--metadatasize"
	poolMetadataSizeArg    = "--poolmetadatasize"
	virtualSizeArg         = "--virtualsize"
	chunkSizeArg           = "--chunksize"
	dataAlignmentArg       = "--dataalignment"
	dataAlignmentOffsetArg = "--dataalignmentoffset"
)

type Unit rune

func (unit Unit) String() string {
	if unit == UnitUnknown {
		return ""
	}
	return string(unit)
}

func (unit Unit) MarshalText() ([]byte, error) {
	return []byte(unit.String()), nil
}

func (unit Unit) ApplyToArgs(args Arguments) error {
	if unit == 0 {
		return nil
	}

	if err := unit.Validate(); err != nil {
		return err
	}

	args.AddOrReplace(fmt.Sprintf("--units=%s", unit.String()))

	return nil
}

func (unit Unit) Validate() error {
	var ok bool
	for _, valid := range validUnits {
		if valid == unit || strings.ToUpper(string(valid)) == string(unit) {
			ok = true
		}
	}
	if !ok {
		return fmt.Errorf("%w: %s is not a valid unit", ErrInvalidUnit, unit)
	}
	return nil
}

func (unit Unit) ApplyToLVsOptions(opts *LVsOptions) {
	opts.Unit = unit
}
func (unit Unit) ApplyToVGsOptions(opts *VGsOptions) {
	opts.Unit = unit
}
func (unit Unit) ApplyToPVsOptions(opts *PVsOptions) {
	opts.Unit = unit
}

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
	UnitUnknown Unit = 0
)

var validUnits = []Unit{UnitBytes, UnitKiB, UnitMiB, UnitGiB, UnitTiB, UnitPiB, UnitEiB, UnitSector, UnitUnknown}

var InvalidSize = Size{Val: -1, Unit: UnitUnknown}
var InvalidPrefixedSize = PrefixedSize{SizePrefix: SizePrefixNone, Size: InvalidSize}

// IsUnitOrDigit returns true if the unit is a valid unit.
// a valid unit is defined as a unit that is a member of validUnits.
// if the unit is not part of a valid unit, IsUnitOrDigit checks if the unit is a digit.
func IsUnitOrDigit(unit Unit) bool {
	for _, valid := range validUnits {
		if valid == unit || strings.ToUpper(string(valid)) == string(unit) {
			return true
		}
	}
	return unicode.IsDigit(rune(unit))
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

func (opt Size) Virtual() VirtualSize {
	return VirtualSize(opt)
}

func (opt Size) ToPoolMetadata() PoolMetadataSize {
	return PoolMetadataSize(opt)
}

func (opt Size) MarshalText() ([]byte, error) {
	return []byte(opt.String()), nil
}

func (opt Size) ToExtents(extentSize uint64, percent ExtentPercent) (Extents, error) {
	bytes, err := opt.ToUnit(UnitBytes)
	if err != nil {
		return Extents{}, err
	}
	extents := uint64(math.Ceil(bytes.Val / float64(extentSize)))
	return NewExtents(extents, percent), nil
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

	if !IsUnitOrDigit(unit) {
		return InvalidSize, fmt.Errorf("%w: %s is neither a valid unit nor a digit", ErrInvalidUnit, unit)
	}
	if opt.Unit == UnitUnknown {
		return InvalidSize, fmt.Errorf(
			"%w: %q cannot be converted to %q, because the unit is unknown - a valid unit is required for conversion (if you meant to use bytes, specify the unit explicitly)",
			ErrInvalidUnit,
			opt,
			unit,
		)
	}

	return NewSize(convert(opt.Val, opt.Unit, unit), unit), nil
}

func (opt Size) String() string {
	var precision int
	if opt.Unit != UnitBytes {
		precision = 2
	}
	val := strconv.FormatFloat(opt.Val, 'f', precision, 64)
	if opt.Unit == UnitUnknown {
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

func ParseSizeLenient(str string) (Size, error) {
	eval := strings.TrimSpace(str)
	if len(eval) == 0 || eval == "0" {
		return Size{Val: 0, Unit: UnitBytes}, nil
	}
	return ParseSize(str)
}

func ParseSize(str string) (Size, error) {
	if len(str) == 0 {
		return Size{Val: 0, Unit: UnitUnknown}, nil
	}

	var unit Unit
	offset := 0
	if len(str) > 1 && !unicode.IsDigit(rune(str[len(str)-1])) {
		unit = Unit(unicode.ToLower(rune(str[len(str)-1])))
		offset++
	} else {
		unit = UnitUnknown
	}

	if !IsUnitOrDigit(unit) {
		return InvalidSize, fmt.Errorf("%w: %s is neither a valid unit nor a digit", ErrInvalidUnit, unit)
	}

	if prefix := str[0]; prefix == '<' || prefix == '>' {
		slog.Warn("size string starts with '<' or '>', this is not supported by lvm2go without losing precision, specify the unit explicitly if possible", slog.String("size", str))
		str = str[1:]
	}

	fval, err := strconv.ParseFloat(str[:len(str)-offset], 64)
	if err != nil {
		return InvalidSize, fmt.Errorf("the value of the size cannot be parsed: %w", err)
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
	if opt.Val < 0 {
		return ErrInvalidSizeGEZero
	}

	if !IsUnitOrDigit(opt.Unit) {
		return ErrInvalidUnit
	}

	return nil
}

func (opt Size) ApplyToLVCreateOptions(opts *LVCreateOptions) {
	opts.Size = opt
}

func (opt Size) ApplyToLVResizeOptions(opts *LVResizeOptions) {
	opts.Size = opt
}

func (opt Size) ApplyToArgs(args Arguments) error {
	return opt.applyToArgs(sizeArg, args)
}

func (opt Size) applyToArgs(arg string, args Arguments) error {
	if err := opt.Validate(); err != nil {
		return err
	}

	args.AddOrReplace(fmt.Sprintf("%s=%s", arg, opt.String()))

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
			return InvalidPrefixedSize, err
		}
		return PrefixedSize{Size: size}, nil
	}

	prefix := SizePrefix(str[0])

	if unicode.IsDigit(rune(prefix)) {
		size, err := ParseSize(str)
		if err != nil {
			return InvalidPrefixedSize, err
		}
		return NewPrefixedSize(SizePrefixNone, size), nil
	}

	if !slices.Contains(prefixCandidates, prefix) {
		return InvalidPrefixedSize, ErrInvalidSizePrefix
	}

	size, err := ParseSize(str[1:])
	if err != nil {
		return InvalidPrefixedSize, err
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

	if opt.SizePrefix != 0 && !slices.Contains(prefixCandidates, opt.SizePrefix) {
		return ErrInvalidSizePrefix
	}

	return nil
}

func (opt PrefixedSize) ApplyToArgs(args Arguments) error {
	return opt.applyToArgs(sizeArg, args)
}

func (opt PrefixedSize) applyToArgs(arg string, args Arguments) error {
	if err := opt.Validate(); err != nil {
		return err
	}
	if opt.Val == 0 {
		return nil
	}

	var sizeBuilder strings.Builder
	if opt.SizePrefix != 0 {
		sizeBuilder.WriteRune(rune(opt.SizePrefix))
	}
	sizeBuilder.WriteString(opt.Size.String())

	args.AddOrReplace(fmt.Sprintf("%s=%s", arg, sizeBuilder.String()))

	return nil
}

func (opt PrefixedSize) ApplyToLVResizeOptions(opts *LVResizeOptions) {
	opts.PrefixedSize = opt
}

func (opt PrefixedSize) ApplyToLVExtendOptions(opts *LVExtendOptions) {
	opts.PrefixedSize = opt
}

type PoolMetadataPrefixedSize PrefixedSize

func (opt PoolMetadataPrefixedSize) ApplyToArgs(args Arguments) error {
	return PrefixedSize(opt).applyToArgs(poolMetadataSizeArg, args)
}

type PoolMetadataSize Size

func (opt PoolMetadataSize) ApplyToArgs(args Arguments) error {
	return Size(opt).applyToArgs(poolMetadataSizeArg, args)
}

type VirtualSize Size

func (opt VirtualSize) ApplyToLVCreateOptions(opts *LVCreateOptions) {
	opts.VirtualSize = opt
}

func (opt VirtualSize) ApplyToArgs(args Arguments) error {
	return Size(opt).applyToArgs(virtualSizeArg, args)
}

type VirtualPrefixedSize PrefixedSize

func (opt VirtualPrefixedSize) ApplyToArgs(args Arguments) error {
	return PrefixedSize(opt).applyToArgs(virtualSizeArg, args)
}

type ChunkSize Size

func (opt ChunkSize) ApplyToLVCreateOptions(opts *LVCreateOptions) {
	opts.ChunkSize = opt
}

func (opt ChunkSize) ApplyToArgs(args Arguments) error {
	return Size(opt).applyToArgs(chunkSizeArg, args)
}

type DataAlignment Size

func (opt DataAlignment) ApplyToArgs(args Arguments) error {
	return Size(opt).applyToArgs(dataAlignmentArg, args)
}

func (opt DataAlignment) ApplyToVGCreateOptions(opts *VGCreateOptions) {
	opts.DataAlignment = opt
}

func (opt DataAlignment) ApplyToPVCreateOptions(opts *PVCreateOptions) {
	opts.DataAlignment = opt
}

type DataAlignmentOffset Size

func (opt DataAlignmentOffset) ApplyToArgs(args Arguments) error {
	return Size(opt).applyToArgs(dataAlignmentOffsetArg, args)
}

func (opt DataAlignmentOffset) ApplyToVGCreateOptions(opts *VGCreateOptions) {
	opts.DataAlignmentOffset = opt
}

func (opt DataAlignmentOffset) ApplyToPVCreateOptions(opts *PVCreateOptions) {
	opts.DataAlignmentOffset = opt
}

type MetadataSize Size

func (opt MetadataSize) ApplyToArgs(args Arguments) error {
	return Size(opt).applyToArgs(metadataSizeArg, args)
}

func (opt MetadataSize) ApplyToVGCreateOptions(opts *VGCreateOptions) {
	opts.MetadataSize = opt
}

func (opt MetadataSize) ApplyToPVCreateOptions(opts *PVCreateOptions) {
	opts.MetadataSize = opt
}
