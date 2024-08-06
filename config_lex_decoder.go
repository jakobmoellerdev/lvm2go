package lvm2go

import (
	"fmt"
	"io"
	"strconv"
)

type LexingDecoder interface {
	ConfigLexerReader
	Decode(v any) error
}

func NewLexingDecoder(reader io.Reader) LexingDecoder {
	return &configLexDecoder{
		ConfigLexerReader: NewConfigLexer(reader),
	}
}

type configLexDecoder struct {
	ConfigLexerReader
}

func (d *configLexDecoder) Decode(v any) error {
	fieldSpecs, err := readLVMStructTag(v)
	if err != nil {
		return err
	}

	fieldSpecsKeyed := make(map[string]map[string]lvmStructTagFieldSpec)
	for _, fieldSpec := range fieldSpecs {
		keyed, ok := fieldSpecsKeyed[fieldSpec.prefix]
		if !ok {
			fieldSpecsKeyed[fieldSpec.prefix] = make(map[string]lvmStructTagFieldSpec)
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

			switch valueInTree.Type {
			case ConfigTokenTypeString:
				fieldSpecsKeyed[section][key].SetString(string(valueInTree.Value))
			case ConfigTokenTypeInt64:
				if val, err := strconv.ParseInt(string(valueInTree.Value), 10, 64); err != nil {
					return fmt.Errorf("could not parse int64: %w", err)
				} else {
					fieldSpecsKeyed[section][key].SetInt(val)
				}
			default:
				return fmt.Errorf("unexpected value type %s", valueInTree.Type)
			}
		}
	}

	return nil
}
