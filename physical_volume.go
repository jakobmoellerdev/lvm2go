/*
 Copyright 2024 The lvm2go Authors.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package lvm2go

import (
	"encoding/json"
	"errors"
)

var ErrPhysicalVolumeNameRequired = errors.New("PhysicalVolumeName is required")

// PhysicalVolumeNameUnknown is a placeholder for an unknown physical volume name used by lvm2
// in case of failure to retrieve the name.
const PhysicalVolumeNameUnknown = PhysicalVolumeName("[unknown]")

type PhysicalVolume struct {
	UUID         string             `json:"pv_uuid"`
	Name         PhysicalVolumeName `json:"pv_name"`
	DevSize      Size               `json:"dev_size"`
	Attr         PVAttributes       `json:"pv_attr"`
	Major        int64              `json:"pv_major"`
	Minor        int64              `json:"pv_minor"`
	MdaFree      Size               `json:"pv_mda_free"`
	MdaSize      Size               `json:"pv_mda_size"`
	PeStart      Size               `json:"pe_start"`
	Size         Size               `json:"pv_size"`
	Free         Size               `json:"pv_free"`
	Used         Size               `json:"pv_used"`
	MdaCount     int64              `json:"pv_mda_count"`
	MdaUsedCount int64              `json:"pv_mda_used_count"`
	Tags         Tags               `json:"pv_tags"`
	VGName       VolumeGroupName    `json:"vg_name"`
	DeviceID     string             `json:"pv_device_id"`
	DeviceIDType string             `json:"pv_device_id_type"`
}

func (pv *PhysicalVolume) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	for key, fieldPtr := range map[string]*string{
		"pv_uuid":      &pv.UUID,
		"pv_name":      (*string)(&pv.Name),
		"vg_name":      (*string)(&pv.VGName),
		"pv_device_id": &pv.DeviceID,
	} {
		if val, ok := raw[key]; !ok {
			continue
		} else if err := json.Unmarshal(val, fieldPtr); err != nil {
			return err
		}
	}

	for key, fieldPtr := range map[string]*Tags{
		"pv_tags": &pv.Tags,
	} {
		if err := unmarshalToStringAndParseCommaSeparatedStrings(raw, key, (*[]string)(fieldPtr)); err != nil {
			return err
		}
	}

	for key, fieldPtr := range map[string]*Size{
		"dev_size":    &pv.DevSize,
		"pv_size":     &pv.Size,
		"pv_free":     &pv.Free,
		"pv_used":     &pv.Used,
		"pv_mda_free": &pv.MdaFree,
		"pv_mda_size": &pv.MdaSize,
		"pe_start":    &pv.PeStart,
	} {
		if err := unmarshalToStringAndParse(raw, key, fieldPtr, ParseSizeLenient); err != nil {
			return err
		}
	}

	for key, fieldPtr := range map[string]*int64{
		"pv_major":          &pv.Major,
		"pv_minor":          &pv.Minor,
		"pv_mda_count":      &pv.MdaCount,
		"pv_mda_used_count": &pv.MdaUsedCount,
	} {
		if err := unmarshalToStringAndParseInt64(raw, key, fieldPtr); err != nil {
			return err
		}
	}

	return unmarshalToStringAndParse(raw, "pv_attr", &pv.Attr, ParsePVAttributes)
}

type PhysicalVolumeName string

var _ Argument = PhysicalVolumeName("")

func (opt PhysicalVolumeName) ApplyToArgs(args Arguments) error {
	args.AddOrReplaceAll([]string{string(opt)})
	return nil
}

func (opt PhysicalVolumeName) ApplyToVGCreateOptions(opts *VGCreateOptions) {
	opts.PhysicalVolumeNames = append(opts.PhysicalVolumeNames, opt)
}
func (opt PhysicalVolumeName) ApplyToVGExtendOptions(opts *VGExtendOptions) {
	opts.PhysicalVolumeNames = append(opts.PhysicalVolumeNames, opt)
}
func (opt PhysicalVolumeName) ApplyToVGReduceOptions(opts *VGReduceOptions) {
	opts.PhysicalVolumeNames = append(opts.PhysicalVolumeNames, opt)
}
func (opt PhysicalVolumeName) ApplyToPVChangeOptions(opts *PVChangeOptions) {
	opts.PhysicalVolumeName = opt
}
func (opt PhysicalVolumeName) ApplyToPVRemoveOptions(opts *PVRemoveOptions) {
	opts.PhysicalVolumeName = opt
}
func (opt PhysicalVolumeName) ApplyToPVMoveOptions(opts *PVMoveOptions) {
	opts.SetOldOrNew(opt)
}

type PhysicalVolumeNames []PhysicalVolumeName

func (opt PhysicalVolumeNames) ApplyToVGReduceOptions(opts *VGReduceOptions) {
	for _, name := range opt {
		name.ApplyToVGReduceOptions(opts)
	}
}

func (opt PhysicalVolumeNames) ApplyToVGExtendOptions(opts *VGExtendOptions) {
	for _, name := range opt {
		name.ApplyToVGExtendOptions(opts)
	}
}

func (opt PhysicalVolumeNames) ApplyToPVMoveOptions(opts *PVMoveOptions) {
	if opts.From == "" {
		opts.From = opt[0]
		opts.To = append(opts.To, opt[1:]...)
	} else {
		opts.To = append(opts.To, opt...)
	}
}

func PhysicalVolumesFrom(names ...string) PhysicalVolumeNames {
	opts := make(PhysicalVolumeNames, len(names))
	for i, v := range names {
		opts[i] = PhysicalVolumeName(v)
	}
	return opts
}

var _ Argument = PhysicalVolumeNames{}

func (opt PhysicalVolumeNames) ApplyToArgs(args Arguments) error {
	raw := make([]string, len(opt))
	for i, v := range opt {
		raw[i] = string(v)
	}
	args.AddOrReplaceAll(raw)
	return nil
}

func (opt PhysicalVolumeNames) ApplyToVGCreateOptions(opts *VGCreateOptions) {
	opts.PhysicalVolumeNames = append(opts.PhysicalVolumeNames, opt...)
}
