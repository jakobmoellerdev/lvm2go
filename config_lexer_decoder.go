package lvm2go

import (
	"fmt"
	"io"
	"strconv"
)

type LexingConfigDecoder interface {
	ConfigLexerReader
	Decode(v any) error
}

type StructuredLexingConfigDecoder interface {
	DecodeStructured(v any) error
}

type UnstructuredLexingConfigDecoder interface {
	DecodeUnstructured(v any) error
}

func NewLexingConfigDecoder(reader io.Reader) LexingConfigDecoder {
	return &configLexDecoder{
		ConfigLexerReader: NewConfigLexer(reader),
	}
}

type configLexDecoder struct {
	ConfigLexerReader
}

func (d *configLexDecoder) Decode(v any) error {
	if isUnstructuredMap(v) {
		return d.DecodeUnstructured(v)
	}
	return d.DecodeStructured(v)
}

func isUnstructuredMap(v any) bool {
	switch v.(type) {
	case map[string]interface{}, *map[string]interface{}:
		return true
	}
	return false
}

func (d *configLexDecoder) DecodeUnstructured(v any) error {
	lexTree, err := d.Lex()
	if err != nil {
		return err
	}

	m, ok := v.(map[string]interface{})
	if !ok {
		mptr, ok := v.(*map[string]interface{})
		if !ok {
			return fmt.Errorf("expected map[string]interface{} or *map[string]interface{}, got %T", v)
		}
		m = *mptr
	}

	var section string
	for i, node := range lexTree {
		if node.Type == ConfigTokenTypeSection {
			section = string(node.Value)
			continue
		}
		if node.Type == ConfigTokenTypeSectionEnd {
			section = ""
			continue
		}
		if node.Type == ConfigTokenTypeAssignment {
			kidx := i - 1
			if kidx < 0 {
				return fmt.Errorf("expected identifier before assignment")
			}
			keyInTree := lexTree[i-1]
			if keyInTree.Type != ConfigTokenTypeIdentifier {
				return fmt.Errorf("expected identifier before assignment, got %s", keyInTree.Type)
			}
			key := string(keyInTree.Value)

			vidx := i + 1
			if vidx >= len(lexTree) {
				return fmt.Errorf("expected value after assignment")
			}
			valueInTree := lexTree[i+1]

			if section != "" {
				key = section + "/" + key
			}

			switch valueInTree.Type {
			case ConfigTokenTypeString:
				m[key] = string(valueInTree.Value)
			case ConfigTokenTypeInt64:
				if val, err := strconv.ParseInt(string(valueInTree.Value), 10, 64); err != nil {
					return fmt.Errorf("could not parse int64: %w", err)
				} else {
					m[key] = val
				}
			default:
				return fmt.Errorf("unexpected value type %s", valueInTree.Type)
			}
		}
	}
	return nil
}

func (d *configLexDecoder) DecodeStructured(v any) error {
	fieldSpecs, err := readLVMStructTag(v)
	if err != nil {
		return err
	}
	decoder := &structuredConfigLexDecoder{
		ConfigLexerReader:      d.ConfigLexerReader,
		StructuredFieldMapping: fieldSpecs,
	}
	return decoder.Decode()
}

func newLexingConfigDecoderWithFieldMapping(
	reader io.Reader,
	fieldSpecs map[string]LVMStructTagFieldMapping,
) *structuredConfigLexDecoder {
	return &structuredConfigLexDecoder{
		ConfigLexerReader:      NewConfigLexer(reader),
		StructuredFieldMapping: fieldSpecs,
		MapHints:               newHintsFromFieldSpecs(fieldSpecs),
	}
}

func newHintsFromFieldSpecs(keys map[string]LVMStructTagFieldMapping) map[string]structuredDecodeHint {
	hints := make(map[string]structuredDecodeHint)
	for _, key := range keys {
		hints[key.name] = structuredDecodeHint{
			section: key.prefix,
		}
	}
	return hints
}

type structuredConfigLexDecoder struct {
	ConfigLexerReader
	StructuredFieldMapping map[string]LVMStructTagFieldMapping
	MapHints               map[string]structuredDecodeHint
}

type structuredDecodeHint struct {
	section string
}

func (d *structuredConfigLexDecoder) Decode() error {
	fieldSpecsKeyed := make(map[string]map[string]LVMStructTagFieldMapping)
	for _, fieldSpec := range d.StructuredFieldMapping {
		keyed, ok := fieldSpecsKeyed[fieldSpec.prefix]
		if !ok {
			fieldSpecsKeyed[fieldSpec.prefix] = make(map[string]LVMStructTagFieldMapping)
			keyed = fieldSpecsKeyed[fieldSpec.prefix]
		}
		keyed[fieldSpec.name] = fieldSpec
	}

	lexTree, err := d.Lex()
	if err != nil {
		return err
	}

	var section string
	for i, node := range lexTree {
		if node.Type == ConfigTokenTypeSection {
			section = string(node.Value)
			continue
		}
		if node.Type == ConfigTokenTypeSectionEnd {
			section = ""
			continue
		}
		if node.Type == ConfigTokenTypeAssignment {
			kidx := i - 1
			if kidx < 0 {
				return fmt.Errorf("expected identifier before assignment")
			}
			keyInTree := lexTree[i-1]
			if keyInTree.Type != ConfigTokenTypeIdentifier {
				return fmt.Errorf("expected identifier before assignment, got %s", keyInTree.Type)
			}
			key := string(keyInTree.Value)

			vidx := i + 1
			if vidx >= len(lexTree) {
				return fmt.Errorf("expected value after assignment")
			}
			valueInTree := lexTree[i+1]

			if d.MapHints != nil {
				if hint, ok := d.MapHints[key]; ok {
					if hint.section != "" {
						section = hint.section
					}
				}
			}

			section, ok := fieldSpecsKeyed[section]
			if !ok {
				continue
			}
			field, ok := section[key]
			if !ok {
				continue
			}

			switch valueInTree.Type {
			case ConfigTokenTypeString:
				field.SetString(string(valueInTree.Value))
			case ConfigTokenTypeInt64:
				if val, err := strconv.ParseInt(string(valueInTree.Value), 10, 64); err != nil {
					return fmt.Errorf("could not parse int64: %w", err)
				} else {
					field.SetInt(val)
				}
			default:
				return fmt.Errorf("unexpected value type %s", valueInTree.Type)
			}
		}
	}
	return nil
}
