package lvm2go

import (
	"encoding/json"
	"strconv"
)

func unmarshalAndConvertToSize(raw map[string]json.RawMessage, key string, fieldPtr *Size) error {
	if raw[key] == nil || len(raw[key]) == 0 {
		*fieldPtr = NewSize(0, UnitBytes)
		return nil
	}
	var str string
	if err := json.Unmarshal(raw[key], &str); err != nil {
		return err
	}
	if str == "" {
		*fieldPtr = NewSize(0, UnitBytes)
		return nil
	}

	size, err := ParseSize(str)
	if err != nil {
		return err
	}

	*fieldPtr = size

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
