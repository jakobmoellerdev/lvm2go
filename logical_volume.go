package lvm2go

import (
	"encoding/json"
	"errors"
	"fmt"
)

var ErrVolumeGroupNameRequired = errors.New("VolumeGroupName is required for a fully qualified logical volume")
var ErrLogicalVolumeNameRequired = errors.New("LogicalVolumeName is required for a fully qualified logical volume")

type LogicalVolume struct {
	UUID     string            `json:"lv_uuid"`
	Name     LogicalVolumeName `json:"lv_name"`
	FullName string            `json:"lv_full_name"`

	Path  string `json:"lv_path"`
	Major int64  `json:"lv_kernel_major"`
	Minor int64  `json:"lv_kernel_minor"`

	Tags Tags         `json:"lv_tags"`
	Attr LVAttributes `json:"lv_attr"`
	Size Size         `json:"lv_size"`

	Origin            string `json:"origin"`
	OriginSize        Size   `json:"origin_size"`
	PoolLogicalVolume string `json:"pool_lv"`

	VolumeGroupName VolumeGroupName `json:"vg_name"`

	DataPercent     float64 `json:"data_percent"`
	MetadataPercent float64 `json:"metadata_percent"`
}

func (lv *LogicalVolume) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	for key, fieldPtr := range map[string]*string{
		"lv_uuid":      &lv.UUID,
		"lv_name":      (*string)(&lv.Name),
		"lv_full_name": &lv.FullName,
		"lv_path":      &lv.Path,
		"origin":       &lv.Origin,
		"pool_lv":      &lv.PoolLogicalVolume,
		"vg_name":      (*string)(&lv.VolumeGroupName),
	} {
		if val, ok := raw[key]; !ok {
			continue
		} else if err := json.Unmarshal(val, fieldPtr); err != nil {
			return err
		}
	}

	for key, fieldPtr := range map[string]*Tags{
		"lv_tags": &lv.Tags,
	} {
		if err := unmarshalToStringAndParseCommaSeparatedStrings(raw, key, (*[]string)(fieldPtr)); err != nil {
			return err
		}
	}

	for key, fieldPtr := range map[string]*int64{
		"lv_kernel_major": &lv.Major,
		"lv_kernel_minor": &lv.Minor,
	} {
		if err := unmarshalToStringAndParseInt64(raw, key, fieldPtr); err != nil {
			return err
		}
	}

	for key, fieldPtr := range map[string]*float64{
		"data_percent":     &lv.DataPercent,
		"metadata_percent": &lv.MetadataPercent,
	} {
		if err := unmarshalToStringAndParseFloat64(raw, key, fieldPtr); err != nil {
			return err
		}
	}

	for key, fieldPtr := range map[string]*Size{
		"lv_size":     &lv.Size,
		"origin_size": &lv.OriginSize,
	} {
		if err := unmarshalToStringAndParse(raw, key, fieldPtr, ParseSizeLenient); err != nil {
			return err
		}
	}

	return unmarshalToStringAndParse(raw, "lv_attr", &lv.Attr, ParseLVAttributes)
}

func (lv *LogicalVolume) GetFQLogicalVolumeName() (*FQLogicalVolumeName, error) {
	fq, err := NewFQLogicalVolumeName(lv.VolumeGroupName, lv.Name)
	if err != nil {
		return nil, err
	}
	return fq, fq.Validate()
}

type LogicalVolumeName string

func (opt LogicalVolumeName) ApplyToLVRenameOptions(opts *LVRenameOptions) {
	opts.SetOldOrNew(opt)
}

func (opt LogicalVolumeName) ApplyToLVExtendOptions(opts *LVExtendOptions) {
	opts.LogicalVolumeName = opt
}

func (opt LogicalVolumeName) ApplyToLVCreateOptions(opts *LVCreateOptions) {
	opts.LogicalVolumeName = opt
}

func (opt LogicalVolumeName) ApplyToLVRemoveOptions(opts *LVRemoveOptions) {
	opts.LogicalVolumeName = opt
}

func (opt LogicalVolumeName) ApplyToLVResizeOptions(opts *LVResizeOptions) {
	opts.LogicalVolumeName = opt
}

func (opt LogicalVolumeName) ApplyToLVChangeOptions(opts *LVChangeOptions) {
	opts.LogicalVolumeName = opt
}

func (opt LogicalVolumeName) ApplyToLVsOptions(opts *LVsOptions) {
	opts.LogicalVolumeName = opt
}

func (opt LogicalVolumeName) ApplyToLVReduceOptions(opts *LVReduceOptions) {
	opts.LogicalVolumeName = opt
}

func (opt LogicalVolumeName) ApplyToPVMoveOptions(opts *PVMoveOptions) {
	opts.LogicalVolumeName = opt
}

type FQLogicalVolumeName struct {
	VolumeGroupName
	LogicalVolumeName
}

func (opt *FQLogicalVolumeName) ApplyToLVRemoveOptions(opts *LVRemoveOptions) {
	opts.VolumeGroupName, opts.LogicalVolumeName = opt.VolumeGroupName, opt.LogicalVolumeName
}

func (opt *FQLogicalVolumeName) ApplyToLVCreateOptions(opts *LVCreateOptions) {
	opts.VolumeGroupName, opts.LogicalVolumeName = opt.VolumeGroupName, opt.LogicalVolumeName
}

func (opt *FQLogicalVolumeName) ApplyToLVExtendOptions(opts *LVExtendOptions) {
	opts.VolumeGroupName, opts.LogicalVolumeName = opt.VolumeGroupName, opt.LogicalVolumeName
}

func (opt *FQLogicalVolumeName) ApplyToLVChangeOptions(opts *LVChangeOptions) {
	opts.VolumeGroupName, opts.LogicalVolumeName = opt.VolumeGroupName, opt.LogicalVolumeName
}

func (opt *FQLogicalVolumeName) ApplyToLVResizeOptions(opts *LVResizeOptions) {
	opts.VolumeGroupName, opts.LogicalVolumeName = opt.VolumeGroupName, opt.LogicalVolumeName
}

func (opt *FQLogicalVolumeName) ApplyToLVReduceOptions(opts *LVReduceOptions) {
	opts.VolumeGroupName, opts.LogicalVolumeName = opt.VolumeGroupName, opt.LogicalVolumeName
}

func (opt *FQLogicalVolumeName) ApplyToLVRenameOptions(opts *LVRenameOptions) {
	opts.VolumeGroupName = opt.VolumeGroupName
	opts.SetOldOrNew(opt.LogicalVolumeName)
}

func (opt *FQLogicalVolumeName) ApplyToLVsOptions(opts *LVsOptions) {
	opts.VolumeGroupName, opts.LogicalVolumeName = opt.VolumeGroupName, opt.LogicalVolumeName
}

func (opt *FQLogicalVolumeName) Split() (VolumeGroupName, LogicalVolumeName) {
	return opt.VolumeGroupName, opt.LogicalVolumeName
}

func (opt *FQLogicalVolumeName) Validate() error {
	if opt.VolumeGroupName == "" {
		return ErrVolumeGroupNameRequired
	}
	if opt.LogicalVolumeName == "" {
		return ErrLogicalVolumeNameRequired
	}
	return nil
}

func (opt *FQLogicalVolumeName) String() string {
	return fmt.Sprintf("%s/%s", opt.VolumeGroupName, opt.LogicalVolumeName)
}

func (opt *FQLogicalVolumeName) ApplyToArgs(args Arguments) error {
	if opt == nil {
		return nil
	}

	if err := opt.Validate(); err != nil {
		return err
	}

	args.AddOrReplace(fmt.Sprintf("%s/%s", opt.VolumeGroupName, opt.LogicalVolumeName))
	return nil
}

func MustNewFQLogicalVolumeName(vg VolumeGroupName, lv LogicalVolumeName) *FQLogicalVolumeName {
	fq, err := NewFQLogicalVolumeName(vg, lv)
	if err != nil {
		panic(err)
	}
	return fq
}

func NewFQLogicalVolumeName(vg VolumeGroupName, lv LogicalVolumeName) (*FQLogicalVolumeName, error) {
	fq := &FQLogicalVolumeName{vg, lv}
	return fq, fq.Validate()
}

var _ Argument = LogicalVolumeName("")
var _ Argument = (*FQLogicalVolumeName)(nil)

func (opt LogicalVolumeName) ApplyToArgs(args Arguments) error {
	if len(opt) == 0 {
		return nil
	}
	switch args.GetType() {
	case ArgsTypeLVRename:
		args.AddOrReplace(string(opt))
	default:
		args.AddOrReplace(fmt.Sprintf("--name=%s", string(opt)))
	}
	return nil
}
