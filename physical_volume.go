package lvm2go

import (
	"encoding/json"
)

type PhysicalVolume struct {
	UUID         string             `json:"pv_uuid"`
	Name         PhysicalVolumeName `json:"pv_name"`
	DevSize      Size               `json:"dev_size"`
	Major        uint64             `json:"pv_major"`
	Minor        uint64             `json:"pv_minor"`
	MdaFree      Size               `json:"pv_mda_free"`
	MdaSize      Size               `json:"pv_mda_size"`
	PeStart      Size               `json:"pe_start"`
	Size         Size               `json:"pv_size"`
	Free         Size               `json:"pv_free"`
	Used         Size               `json:"pv_used"`
	MdaCount     uint64             `json:"pv_mda_count"`
	MdaUsedCount uint64             `json:"pv_mda_used_count"`
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
		if err := json.Unmarshal(raw[key], fieldPtr); err != nil {
			return err
		}
	}

	for key, fieldPtr := range map[string]*Tags{
		"pv_tags": &pv.Tags,
	} {
		if err := unmarshalAndConvertToStrings(raw, key, (*[]string)(fieldPtr)); err != nil {
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
		if err := unmarshalAndConvertToSize(raw, key, fieldPtr); err != nil {
			return err
		}
	}

	for key, fieldPtr := range map[string]*uint64{
		"pv_major":          &pv.Major,
		"pv_minor":          &pv.Minor,
		"pv_mda_count":      &pv.MdaCount,
		"pv_mda_used_count": &pv.MdaUsedCount,
	} {
		if err := unmarshalAndConvertToUint64(raw, key, fieldPtr); err != nil {
			return err
		}
	}

	return nil
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

type PhysicalVolumeNames []PhysicalVolumeName

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
