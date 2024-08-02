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
	"strconv"
	"strings"
)

func unmarshalToStringAndParseCommaSeparatedStrings(raw map[string]json.RawMessage, key string, fieldPtr *[]string) error {
	return unmarshalToStringAndParse(raw, key, fieldPtr, func(str string) ([]string, error) {
		if len(str) > 0 {
			return strings.Split(str, ","), nil
		}
		return nil, nil
	})
}

func unmarshalToStringAndParseInt64(raw map[string]json.RawMessage, key string, fieldPtr *int64) error {
	return unmarshalToStringAndParse(raw, key, fieldPtr, func(str string) (int64, error) {
		return strconv.ParseInt(str, 10, 64)
	})
}

func unmarshalToStringAndParseFloat64(raw map[string]json.RawMessage, key string, fieldPtr *float64) error {
	return unmarshalToStringAndParse(raw, key, fieldPtr, func(str string) (float64, error) {
		return strconv.ParseFloat(str, 64)
	})
}

func unmarshalToStringAndParse[T any](
	raw map[string]json.RawMessage,
	key string,
	fieldPtr *T,
	parse func(str string) (T, error),
) error {
	if raw[key] == nil || len(raw[key]) == 0 {
		*fieldPtr = *new(T)
		return nil
	}
	var str string
	if err := json.Unmarshal(raw[key], &str); err != nil {
		return err
	}
	if str == "" {
		*fieldPtr = *new(T)
		return nil
	}
	attrs, err := parse(str)
	if err != nil {
		return err
	}
	*fieldPtr = attrs
	return nil
}
