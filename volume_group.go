package lvm2go

import (
	"encoding/json"
)

type VolumeGroup struct {
	UUID         string          `json:"vg_uuid"`
	Name         VolumeGroupName `json:"vg_name"`
	SysID        string          `json:"vg_sysid"`
	LockType     string          `json:"vg_lock_type"`
	LockArgs     string          `json:"vg_lock_args"`
	VGAttributes string          `json:"vg_attr"`
	Tags         Tags            `json:"vg_tags"`

	ExtentSize     Size  `json:"vg_extent_size"`
	ExtentCount    int64 `json:"vg_extent_count"`
	SeqNo          int64 `json:"vg_seqno"`
	Size           Size  `json:"vg_size"`
	Free           Size  `json:"vg_free"`
	PvCount        int64 `json:"pv_count"`
	MissingPVCount int64 `json:"missing_pv_count"`
	MaxPv          int64 `json:"max_pv"`
	LvCount        int64 `json:"lv_count"`
	MaxLv          int64 `json:"max_lv"`
	SnapCount      int64 `json:"snap_count"`
	MDACount       int64 `json:"mda_count"`
	MDAUsedCount   int64 `json:"mda_used_count"`
	MDACopies      int64 `json:"mda_copies"`
}

func (vg *VolumeGroup) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	for key, fieldPtr := range map[string]*string{
		"vg_uuid":      &vg.UUID,
		"vg_name":      (*string)(&vg.Name),
		"vg_sysid":     &vg.SysID,
		"vg_lock_type": &vg.LockType,
		"vg_lock_args": &vg.LockArgs,
		"vg_attr":      &vg.VGAttributes,
	} {
		if val, ok := raw[key]; !ok {
			continue
		} else if err := json.Unmarshal(val, fieldPtr); err != nil {
			return err
		}
	}

	for key, fieldPtr := range map[string]*Tags{
		"vg_tags": &vg.Tags,
	} {
		if err := unmarshalAndConvertToStrings(raw, key, (*[]string)(fieldPtr)); err != nil {
			return err
		}
	}

	for key, fieldPtr := range map[string]*int64{
		"vg_extent_count":  &vg.ExtentCount,
		"pv_count":         &vg.PvCount,
		"missing_pv_count": &vg.MissingPVCount,
		"max_pv":           &vg.MaxPv,
		"lv_count":         &vg.LvCount,
		"max_lv":           &vg.MaxLv,
		"snap_count":       &vg.SnapCount,
		"mda_count":        &vg.MDACount,
		"mda_used_count":   &vg.MDAUsedCount,
		"mda_copies":       &vg.MDACopies,
		"vg_seqno":         &vg.SeqNo,
	} {
		if err := unmarshalAndConvertToInt64(raw, key, fieldPtr); err != nil {
			return err
		}
	}

	for key, fieldPtr := range map[string]*Size{
		"vg_size":        &vg.Size,
		"vg_free":        &vg.Free,
		"vg_extent_size": &vg.ExtentSize,
	} {
		if err := unmarshalAndConvertToSize(raw, key, fieldPtr); err != nil {
			return err
		}
	}

	return nil
}

type VolumeGroupName string

var _ Argument = VolumeGroupName("")

func (opt VolumeGroupName) ApplyToLVsOptions(opts *LVsOptions) {
	opts.VolumeGroupName = opt
}
func (opt VolumeGroupName) ApplyToLVRenameOptions(opts *LVRenameOptions) {
	opts.VolumeGroupName = opt
}
func (opt VolumeGroupName) ApplyToLVChangeOptions(opts *LVChangeOptions) {
	opts.VolumeGroupName = opt
}

func (opt VolumeGroupName) ApplyToLVExtendOptions(opts *LVExtendOptions) {
	opts.VolumeGroupName = opt
}
func (opt VolumeGroupName) ApplyToVGsOptions(opts *VGsOptions) {
	opts.VolumeGroupName = opt
}
func (opt VolumeGroupName) ApplyToVGCreateOptions(opts *VGCreateOptions) {
	opts.VolumeGroupName = opt
}
func (opt VolumeGroupName) ApplyToLVCreateOptions(opts *LVCreateOptions) {
	opts.VolumeGroupName = opt
}
func (opt VolumeGroupName) ApplyToVGRemoveOptions(opts *VGRemoveOptions) {
	opts.VolumeGroupName = opt
}
func (opt VolumeGroupName) ApplyToVGExtendOptions(opts *VGExtendOptions) {
	opts.VolumeGroupName = opt
}
func (opt VolumeGroupName) ApplyToVGRenameOptions(opts *VGRenameOptions) {
	opts.SetOldOrNew(opt)
}
func (opt VolumeGroupName) ApplyToVGChangeOptions(opts *VGChangeOptions) {
	opts.VolumeGroupName = opt
}
func (opt VolumeGroupName) ApplyToVGReduceOptions(opts *VGReduceOptions) {
	opts.VolumeGroupName = opt
}
func (opt VolumeGroupName) ApplyToLVRemoveOptions(opts *LVRemoveOptions) {
	opts.VolumeGroupName = opt
}
func (opt VolumeGroupName) ApplyToLVResizeOptions(opts *LVResizeOptions) {
	opts.VolumeGroupName = opt
}

func (opt VolumeGroupName) ApplyToArgs(args Arguments) error {
	if len(opt) > 0 {
		args.AddOrReplace(string(opt))
	}
	return nil
}
