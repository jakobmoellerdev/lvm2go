package lvm2go

import (
	"fmt"
	"reflect"
)

type LVMStructTagFieldMapping struct {
	prefix string
	name   string
	reflect.Value
}

func (f LVMStructTagFieldMapping) String() string {
	switch f.Kind() {
	case reflect.Int64:
		return fmt.Sprintf("%s = %d", f.name, f.Int())
	default:
		return fmt.Sprintf("%s = %q", f.name, f.Value.String())
	}
}

func readLVMStructTag(v any) (map[string]LVMStructTagFieldMapping, error) {
	fields, typeAccessor, valueAccessor, err := accessStructOrPointerToStruct(v)
	if err != nil {
		return nil, err
	}

	tagOrIgnore := func(tag reflect.StructTag) (string, bool) {
		return tag.Get(LVMConfigStructTag), tag.Get(LVMConfigStructTag) == "-"
	}

	fieldSpecs := make(map[string]LVMStructTagFieldMapping)
	for i := range fields {
		outerField := typeAccessor(i)
		prefix, ignore := tagOrIgnore(outerField.Tag)
		if ignore {
			continue
		}
		fields, typeAccessor, valueAccessor, err := accessStructOrPointerToStruct(valueAccessor(i))
		if err != nil {
			return nil, err
		}
		for j := range fields {
			innerField := typeAccessor(j)
			name, ignore := tagOrIgnore(innerField.Tag)
			if ignore {
				continue
			}
			fieldSpecs[name] = LVMStructTagFieldMapping{
				prefix,
				name,
				valueAccessor(j),
			}
		}
	}
	return fieldSpecs, nil
}