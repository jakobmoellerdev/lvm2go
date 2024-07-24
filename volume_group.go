package lvm2go

import (
	"encoding/json"
)

type VolumeGroup struct {
	UUID string          `json:"vg_uuid"`
	Name VolumeGroupName `json:"vg_name"`

	Size Size `json:"vg_size"`
	Free Size `json:"vg_free"`
}

func (vg *VolumeGroup) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	for key, fieldPtr := range map[string]*string{
		"vg_uuid": &vg.UUID,
		"vg_name": (*string)(&vg.Name),
	} {
		if err := json.Unmarshal(raw[key], fieldPtr); err != nil {
			return err
		}
	}

	for key, fieldPtr := range map[string]*Size{
		"vg_size": &vg.Size,
		"vg_free": &vg.Free,
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
func (opt VolumeGroupName) ApplyToVGRenameOptions(opts *VGRenameOptions) {
	opts.SetOldOrNew(opt)
}
func (opt VolumeGroupName) ApplyToLVRemoveOptions(opts *LVRemoveOptions) {
	opts.VolumeGroupName = opt
}

func (opt VolumeGroupName) ApplyToLVResizeOptions(opts *LVResizeOptions) {
	opts.VolumeGroupName = opt
}

func (opt VolumeGroupName) ApplyToArgs(args Arguments) error {
	args.AddOrReplace(string(opt))
	return nil
}
