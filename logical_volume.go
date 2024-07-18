package lvm2go

import (
	"encoding/json"
	"errors"
	"strings"
)

var ErrVolumeGroupNameRequired = errors.New("VolumeGroupName is required for a fully qualified logical volume")
var ErrLogicalVolumeNameRequired = errors.New("LogicalVolumeName is required for a fully qualified logical volume")

type LogicalVolume struct {
	UUID     string            `json:"lv_uuid"`
	Name     LogicalVolumeName `json:"lv_name"`
	FullName string            `json:"lv_full_name"`

	Path  string `json:"lv_path"`
	Major uint64 `json:"lv_kernel_major"`
	Minor uint64 `json:"lv_kernel_minor"`

	Tags string `json:"lv_tags"`
	Attr string `json:"lv_attr"`
	Size Size   `json:"lv_size"`

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
		"lv_tags":      &lv.Tags,
		"lv_attr":      &lv.Attr,
		"origin":       &lv.Origin,
		"pool_lv":      &lv.PoolLogicalVolume,
		"vg_name":      (*string)(&lv.VolumeGroupName),
	} {
		if err := json.Unmarshal(raw[key], fieldPtr); err != nil {
			return err
		}
	}

	for key, fieldPtr := range map[string]*uint64{
		"lv_kernel_major": &lv.Major,
		"lv_kernel_minor": &lv.Minor,
	} {
		if err := unmarshalAndConvertToUint64(raw, key, fieldPtr); err != nil {
			return err
		}
	}

	for key, fieldPtr := range map[string]*float64{
		"data_percent":     &lv.DataPercent,
		"metadata_percent": &lv.MetadataPercent,
	} {
		if err := unmarshalAndConvertToFloat64(raw, key, fieldPtr); err != nil {
			return err
		}
	}

	for key, fieldPtr := range map[string]*Size{
		"lv_size":     &lv.Size,
		"origin_size": &lv.OriginSize,
	} {
		if err := unmarshalAndConvertToSize(raw, key, fieldPtr); err != nil {
			return err
		}
	}

	return nil
}

type LogicalVolumeName string

type FQLogicalVolumeName string

func (opt FQLogicalVolumeName) ApplyToLVRemoveOptions(opts *LVRemoveOptions) {
	opts.VolumeGroupName, opts.LogicalVolumeName = opt.Split()
}

func (opt FQLogicalVolumeName) ApplyToLVCreateOptions(opts *LVCreateOptions) {
	opts.VolumeGroupName, opts.LogicalVolumeName = opt.Split()
}

func (opt FQLogicalVolumeName) ApplyToLVExtendOptions(opts *LVExtendOptions) {
	opts.VolumeGroupName, opts.LogicalVolumeName = opt.Split()
}

func (opt FQLogicalVolumeName) ApplyToLVChangeOptions(opts *LVChangeOptions) {
	opts.VolumeGroupName, opts.LogicalVolumeName = opt.Split()
}

func (opt FQLogicalVolumeName) ApplyToLVResizeOptions(opts *LVResizeOptions) {
	opts.VolumeGroupName, opts.LogicalVolumeName = opt.Split()
}

func (opt FQLogicalVolumeName) ApplyToLVReduceOptions(opts *LVReduceOptions) {
	opts.VolumeGroupName, opts.LogicalVolumeName = opt.Split()
}

func (opt FQLogicalVolumeName) ApplyToLVRenameOptions(opts *LVRenameOptions) {
	vgname, lvname := opt.Split()
	opts.VolumeGroupName = vgname
	opts.SetOldOrNew(lvname)
}

func (opt FQLogicalVolumeName) Split() (VolumeGroupName, LogicalVolumeName) {
	split := strings.Split(string(opt), "/")
	return VolumeGroupName(split[0]), LogicalVolumeName(split[1])
}

func (opt FQLogicalVolumeName) ApplyToArgs(args Arguments) error {
	args.AddOrReplaceAll([]string{string(opt)})
	return nil
}

func NewFQLogicalVolumeName(vg VolumeGroupName, lv LogicalVolumeName) (FQLogicalVolumeName, error) {
	if vg == "" {
		return "", ErrVolumeGroupNameRequired
	}
	if lv == "" {
		return "", ErrLogicalVolumeNameRequired
	}
	return FQLogicalVolumeName(string(vg) + "/" + string(lv)), nil
}

var _ Argument = LogicalVolumeName("")
var _ Argument = FQLogicalVolumeName("")

func (opt LogicalVolumeName) ApplyToLVCreateOptions(opts *LVCreateOptions) {
	opts.LogicalVolumeName = opt
}

func (opt LogicalVolumeName) ApplyToLVRemoveOptions(opts *LVRemoveOptions) {
	opts.LogicalVolumeName = opt
}

func (opt LogicalVolumeName) ApplyToArgs(args Arguments) error {
	switch args.GetType() {
	case ArgsTypeLVRename:
		args.AddOrReplaceAll([]string{string(opt)})
	}

	args.AddOrReplaceAll([]string{"--name", string(opt)})
	return nil
}
