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
	"fmt"
	"reflect"
)

// accessStructOrPointerToStruct returns the number of fields in the struct,
// a function to access the fields, and a function to access the values of the fields.
// If the value is a pointer, its reference will be used or initialized if it is nil.
// If the value is not a struct or a pointer to a struct, an error will be returned.
// The fieldAccessor function returns the field at the given index.
// The valueAccessor function returns the value at the given index.
// The fieldAccessor and valueAccessor functions are safe to use in a loop and will not panic for idx < fieldNum.
// The valueAccessor function will dereference pointers if necessary and initialize them if they are nil.
// The valueAccessor function will panic if idx >= fieldNum.
func accessStructOrPointerToStruct(v interface{}) (
	fieldNum int,
	fieldAccessor func(idx int) reflect.StructField,
	valueAccessor func(idx int) reflect.Value,
	err error,
) {
	var value reflect.Value
	switch v := v.(type) {
	case reflect.Value:
		value = v
	default:
		if value = reflect.ValueOf(v); !value.IsValid() {
			return 0, nil, nil, fmt.Errorf("invalid value")
		}
	}

	value = initPointerIfNeeded(value)
	t := value.Type()

	if t.Kind() != reflect.Struct {
		return 0, nil, nil, fmt.Errorf("expected struct or pointer to struct, got %s", t.Kind())
	}

	fieldNum = t.NumField()
	fieldAccessor = func(idx int) reflect.StructField {
		return t.Field(idx)
	}
	valueAccessor = func(idx int) reflect.Value {
		return initPointerIfNeeded(value.Field(idx))
	}
	return
}

func initPointerIfNeeded(val reflect.Value) reflect.Value {
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			// Initialize the nil pointer to a new instance of the struct
			val.Set(reflect.New(val.Type().Elem()))
		}
		val = val.Elem()
	}
	return val
}
