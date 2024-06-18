package lvm2go

import (
	"encoding/json"
	"strconv"
)

type LogicalVolume struct {
	UUID     string `json:"lv_uuid"`
	Name     string `json:"lv_name"`
	FullName string `json:"lv_full_name"`

	Path  string `json:"lv_path"`
	Major uint64 `json:"lv_kernel_major"`
	Minor uint64 `json:"lv_kernel_minor"`

	Tags string `json:"lv_tags"`
	Attr string `json:"lv_attr"`
	Size uint64 `json:"lv_size"`

	Origin            string `json:"origin"`
	OriginSize        uint64 `json:"origin_size"`
	PoolLogicalVolume string `json:"pool_lv"`

	VolumeGroupName string `json:"vg_name"`

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
		"lv_name":      &lv.Name,
		"lv_full_name": &lv.FullName,
		"lv_path":      &lv.Path,
		"lv_tags":      &lv.Tags,
		"lv_attr":      &lv.Attr,
		"origin":       &lv.Origin,
		"pool_lv":      &lv.PoolLogicalVolume,
		"vg_name":      &lv.VolumeGroupName,
	} {
		if err := json.Unmarshal(raw[key], fieldPtr); err != nil {
			return err
		}
	}

	for key, fieldPtr := range map[string]*uint64{
		"lv_size":         &lv.Size,
		"lv_kernel_major": &lv.Major,
		"lv_kernel_minor": &lv.Minor,
		"origin_size":     &lv.OriginSize,
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

	return nil
}

func unmarshalAndConvertToUint64(raw map[string]json.RawMessage, key string, fieldPtr *uint64) error {
	if raw[key] == nil || len(raw[key]) == 0 {
		*fieldPtr = 0
		return nil
	}
	var str string
	if err := json.Unmarshal(raw[key], &str); err != nil {
		return err
	}
	if str == "" {
		*fieldPtr = 0
		return nil
	}
	val, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return err
	}
	*fieldPtr = val
	return nil
}

func unmarshalAndConvertToFloat64(raw map[string]json.RawMessage, key string, fieldPtr *float64) error {
	if raw[key] == nil || len(raw[key]) == 0 {
		*fieldPtr = 0
		return nil
	}
	var str string
	if err := json.Unmarshal(raw[key], &str); err != nil {
		return err
	}
	if str == "" {
		*fieldPtr = 0
		return nil
	}
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return err
	}
	*fieldPtr = val
	return nil
}
