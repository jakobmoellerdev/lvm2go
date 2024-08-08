package config

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	lvmreflect "github.com/jakobmoellerdev/lvm2go/reflect"
)

const LVMConfigStructTag = "lvm"
const LVMProfileExtension = ".profile"

type LVMStructTagFieldMappings map[string]*LVMStructTagFieldMapping

type LVMStructTagFieldMapping struct {
	Prefix string
	Name   string
	reflect.Value
}

func (f LVMStructTagFieldMapping) String() string {
	switch f.Kind() {
	case reflect.Int64:
		return fmt.Sprintf("%s = %d", f.Name, f.Int())
	default:
		return fmt.Sprintf("%s = %q", f.Name, f.Value.String())
	}
}

func DecodeFieldMappings(v any) (LVMStructTagFieldMappings, error) {
	fields, typeAccessor, valueAccessor, err := lvmreflect.AccessStructOrPointerToStruct(v)
	if err != nil {
		return nil, err
	}

	tagOrIgnore := func(tag reflect.StructTag) (string, bool) {
		return tag.Get(LVMConfigStructTag), tag.Get(LVMConfigStructTag) == "-"
	}

	fieldSpecs := make(LVMStructTagFieldMappings)
	for i := range fields {
		outerField := typeAccessor(i)
		prefix, ignore := tagOrIgnore(outerField.Tag)
		if ignore {
			continue
		}
		fields, typeAccessor, valueAccessor, err := lvmreflect.AccessStructOrPointerToStruct(valueAccessor(i))
		if err != nil {
			return nil, err
		}
		for j := range fields {
			innerField := typeAccessor(j)
			name, ignore := tagOrIgnore(innerField.Tag)
			if ignore {
				continue
			}
			fieldSpecs[name] = &LVMStructTagFieldMapping{
				prefix,
				name,
				valueAccessor(j),
			}
		}
	}
	return fieldSpecs, nil
}

func FieldMappingsToTokens(mappings LVMStructTagFieldMappings) Tokens {
	fieldSpecsKeyed := make(map[string]LVMStructTagFieldMappings)

	sectionKeys := make([]string, 0, len(fieldSpecsKeyed))
	fieldKeys := make([]string, 0, len(fieldSpecsKeyed))

	for _, fieldSpec := range mappings {
		keyed, ok := fieldSpecsKeyed[fieldSpec.Prefix]
		if !ok {
			fieldSpecsKeyed[fieldSpec.Prefix] = make(LVMStructTagFieldMappings)
			keyed = fieldSpecsKeyed[fieldSpec.Prefix]
		}
		keyed[fieldSpec.Name] = fieldSpec
		sectionKeys = append(sectionKeys, fieldSpec.Prefix)
		fieldKeys = append(fieldKeys, fieldSpec.Name)
	}

	// sections should appear only once
	sectionKeys = slices.Compact(sectionKeys)

	// Sort the keys to ensure a deterministic output
	slices.SortStableFunc(sectionKeys, func(i, j string) int {
		return strings.Compare(i, j)
	})
	slices.SortStableFunc(fieldKeys, func(i, j string) int {
		return strings.Compare(i, j)
	})

	ast := NewAST()
	for _, section := range sectionKeys {
		astSection := NewSection(section)
		ast.Append(astSection)
		for _, field := range fieldKeys {
			astSection.Append(NewAssignmentFromSpec(fieldSpecsKeyed[section][field]))
		}
	}
	return ast.Tokens()
}
