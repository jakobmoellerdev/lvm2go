package lvm2go

import (
	"encoding/json"
	"strconv"
	"strings"
)

func unmarshalAndConvertToStrings(raw map[string]json.RawMessage, key string, fieldPtr *[]string) error {
	if raw[key] == nil || len(raw[key]) == 0 {
		*fieldPtr = nil
		return nil
	}

	var str string
	if err := json.Unmarshal(raw[key], &str); err != nil {
		return err
	}

	*fieldPtr = strings.Split(str, ",")

	return nil
}

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

func unmarshalAndConvertToInt64(raw map[string]json.RawMessage, key string, fieldPtr *int64) error {
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
	val, err := strconv.ParseInt(str, 10, 64)
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

func unmarshalAndConvertToLVAttributes(raw map[string]json.RawMessage, key string, fieldPtr *LVAttributes) error {
	if raw[key] == nil || len(raw[key]) == 0 {
		*fieldPtr = LVAttributes{}
		return nil
	}
	var str string
	if err := json.Unmarshal(raw[key], &str); err != nil {
		return err
	}
	if str == "" {
		*fieldPtr = LVAttributes{}
		return nil
	}
	attrs, err := ParsedLVAttributes(str)
	if err != nil {
		return err
	}
	*fieldPtr = attrs
	return nil
}
