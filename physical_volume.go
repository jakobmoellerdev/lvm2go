package lvm2go

import (
	"encoding/json"
)

type PhysicalVolume struct {
	UUID string             `json:"pv_uuid"`
	Name PhysicalVolumeName `json:"pv_name"`
}

func (pv *PhysicalVolume) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	for key, fieldPtr := range map[string]*string{
		"pv_uuid": &pv.UUID,
		"pv_name": (*string)(&pv.Name),
	} {
		if err := json.Unmarshal(raw[key], fieldPtr); err != nil {
			return err
		}
	}

	return nil
}

type PhysicalVolumeName string

var _ Argument = PhysicalVolumeName("")

func (opt PhysicalVolumeName) ApplyToArgs(args Arguments) error {
	args.AppendAll([]string{string(opt)})
	return nil
}

func (opt PhysicalVolumeName) ApplyToVGCreateOptions(opts *VGCreateOptions) {
	opts.PhysicalVolumeNames = append(opts.PhysicalVolumeNames, opt)
}

type PhysicalVolumeNames []PhysicalVolumeName

func PhysicalVolumeNamesFrom(names ...string) PhysicalVolumeNames {
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
	args.AppendAll(raw)
	return nil
}

func (opt PhysicalVolumeNames) ApplyToVGCreateOptions(opts *VGCreateOptions) {
	opts.PhysicalVolumeNames = append(opts.PhysicalVolumeNames, opt...)
}
