package config

import (
	"fmt"
	"io"
	"strconv"
)

type LexingDecoder interface {
	Decode(v any) error
}

type StructuredLexingDecoder interface {
	DecodeStructured(v any) error
}

type UnstructuredLexingDecoder interface {
	DecodeUnstructured(v any) error
}

func NewLexingConfigDecoder(reader io.Reader) LexingDecoder {
	return &lexDecoder{
		LexerReader: NewBufferedLexer(reader),
	}
}

type lexDecoder struct {
	LexerReader
}

var _ LexingDecoder = &lexDecoder{}

func (d *lexDecoder) Decode(v any) error {
	if isUnstructuredMap(v) {
		return d.DecodeUnstructured(v)
	}
	return d.DecodeStructured(v)
}

func (d *lexDecoder) DecodeUnstructured(v any) error {
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
		if node.Type == TokenTypeSection {
			section = string(node.Value)
			continue
		}
		if node.Type == TokenTypeEndOfSection {
			section = ""
			continue
		}
		if node.Type == TokenTypeAssignment {
			kidx := i - 1
			if kidx < 0 {
				return fmt.Errorf("expected identifier before assignment")
			}
			keyInTree := lexTree[i-1]
			if keyInTree.Type != TokenTypeIdentifier {
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
			case TokenTypeString:
				m[key] = string(valueInTree.Value)
			case TokenTypeInt64:
				if val, err := strconv.ParseInt(string(valueInTree.Value), 10, 64); err != nil {
					return fmt.Errorf("could not Parse int64: %w", err)
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

func (d *lexDecoder) DecodeStructured(v any) error {
	fieldSpecs, err := DecodeFieldMappings(v)
	if err != nil {
		return err
	}
	decoder := &structuredLexingDecoder{
		LexerReader:            d.LexerReader,
		structuredFieldMapping: fieldSpecs,
	}
	return decoder.Decode()
}

func NewLexingConfigDecoderWithFieldMapping(
	reader io.Reader,
	fieldSpecs LVMStructTagFieldMappings,
) *structuredLexingDecoder {
	return &structuredLexingDecoder{
		LexerReader:            NewBufferedLexer(reader),
		structuredFieldMapping: fieldSpecs,
		mapHints:               newHintsFromFieldSpecs(fieldSpecs),
	}
}

func newHintsFromFieldSpecs(mappings LVMStructTagFieldMappings) map[string]structuredDecodeHint {
	hints := make(map[string]structuredDecodeHint)
	for _, key := range mappings {
		hints[key.Name] = structuredDecodeHint{
			section: key.Prefix,
		}
	}
	return hints
}

type structuredLexingDecoder struct {
	LexerReader
	structuredFieldMapping LVMStructTagFieldMappings
	mapHints               map[string]structuredDecodeHint
}

type structuredDecodeHint struct {
	section string
}

func (d *structuredLexingDecoder) Decode() error {
	fieldSpecsKeyed := make(map[string]LVMStructTagFieldMappings)
	for _, fieldSpec := range d.structuredFieldMapping {
		keyed, ok := fieldSpecsKeyed[fieldSpec.Prefix]
		if !ok {
			fieldSpecsKeyed[fieldSpec.Prefix] = make(LVMStructTagFieldMappings)
			keyed = fieldSpecsKeyed[fieldSpec.Prefix]
		}
		keyed[fieldSpec.Name] = fieldSpec
	}

	lexTree, err := d.Lex()
	if err != nil {
		return err
	}

	var section string
	for i, node := range lexTree {
		if node.Type == TokenTypeSection {
			section = string(node.Value)
			continue
		}
		if node.Type == TokenTypeEndOfSection {
			section = ""
			continue
		}
		if node.Type == TokenTypeAssignment {
			kidx := i - 1
			if kidx < 0 {
				return fmt.Errorf("expected identifier before assignment")
			}
			keyInTree := lexTree[i-1]
			if keyInTree.Type != TokenTypeIdentifier {
				return fmt.Errorf("expected identifier before assignment, got %s", keyInTree.Type)
			}
			key := string(keyInTree.Value)

			vidx := i + 1
			if vidx >= len(lexTree) {
				return fmt.Errorf("expected value after assignment")
			}
			valueInTree := lexTree[i+1]

			if d.mapHints != nil {
				if hint, ok := d.mapHints[key]; ok {
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
			case TokenTypeString:
				field.SetString(string(valueInTree.Value))
			case TokenTypeInt64:
				if val, err := strconv.ParseInt(string(valueInTree.Value), 10, 64); err != nil {
					return fmt.Errorf("could not Parse int64: %w", err)
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
