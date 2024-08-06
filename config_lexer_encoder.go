package lvm2go

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"
)

type configLexEncoder struct {
	Writer io.Writer
}

func NewLexingConfigEncoder(writer io.Writer) LexingConfigEncoder {
	return &configLexEncoder{
		Writer: writer,
	}
}

func (c *configLexEncoder) Encode(v any) error {
	if isUnstructuredMap(v) {
		return c.EncodeUnstructured(v)
	}
	return c.EncodeStructured(v)
}

func (c *configLexEncoder) EncodeStructured(v any) error {
	fieldSpecs, err := readLVMStructTag(v)
	if err != nil {
		return err
	}
	fieldSpecsKeyed := make(map[string]map[string]LVMStructTagFieldMapping)
	for _, fieldSpec := range fieldSpecs {
		keyed, ok := fieldSpecsKeyed[fieldSpec.prefix]
		if !ok {
			fieldSpecsKeyed[fieldSpec.prefix] = make(map[string]LVMStructTagFieldMapping)
			keyed = fieldSpecsKeyed[fieldSpec.prefix]
		}
		keyed[fieldSpec.name] = fieldSpec
	}

	line := 1
	tokens := make(ConfigTokens, 0, len(fieldSpecs))
	for section, fields := range fieldSpecsKeyed {
		tokens = append(tokens, ConfigToken{
			Type:  ConfigTokenTypeSection,
			Value: []byte(section),
			Line:  line,
		}, ConfigToken{
			Type:  ConfigTokenTypeSectionStart,
			Value: []byte{'{'},
			Line:  line,
		}, ConfigToken{
			Type:  ConfigTokenTypeEndOfStatement,
			Value: []byte{'\n'},
			Line:  line,
		})
		line++

		for _, fieldSpec := range fields {
			tokens = append(tokens, ConfigToken{
				Type:  ConfigTokenTypeIdentifier,
				Value: []byte(fieldSpec.name),
				Line:  line,
			}, ConfigToken{
				Type:  ConfigTokenTypeAssignment,
				Value: []byte{'='},
				Line:  line,
			})

			switch fieldSpec.Kind() {
			case reflect.Int64:
				tokens = append(tokens, ConfigToken{
					Type:  ConfigTokenTypeInt64,
					Value: []byte(fmt.Sprintf("%d", fieldSpec.Value.Int())),
					Line:  line,
				})
			default:
				tokens = append(tokens, ConfigToken{
					Type:  ConfigTokenTypeString,
					Value: []byte(fieldSpec.Value.String()),
					Line:  line,
				})
			}

			tokens = append(tokens, ConfigToken{
				Type:  ConfigTokenTypeEndOfStatement,
				Value: []byte{'\n'},
				Line:  line,
			})
			line++
		}

		tokens = append(tokens, ConfigToken{
			Type:  ConfigTokenTypeSectionEnd,
			Value: []byte{'}'},
			Line:  line,
		}, ConfigToken{
			Type:  ConfigTokenTypeEndOfStatement,
			Value: []byte{'\n'},
			Line:  line,
		})
		line++
	}

	return c.writeTokens(tokens)
}

func (c *configLexEncoder) EncodeUnstructured(v any) error {
	m, ok := v.(map[string]interface{})
	if !ok {
		mptr, ok := v.(*map[string]interface{})
		if !ok {
			return fmt.Errorf("expected map[string]interface{} or *map[string]interface{}, got %T", v)
		}
		m = *mptr
	}

	sectionKeys := make([]string, 0, len(m))
	fieldKeys := make([]string, 0, len(m))
	mbySection := make(map[string]map[string]interface{})
	for k, v := range m {
		splitKey := strings.Split(k, "/")
		section, key := splitKey[0], splitKey[1]

		if section == "" {
			return fmt.Errorf("expected section prefix in key %q", k)
		}
		if key == "" {
			return fmt.Errorf("expected key suffix in key %q", k)
		}

		sectionMap, ok := mbySection[section]
		if !ok {
			mbySection[section] = make(map[string]interface{})
			sectionMap = mbySection[section]
			sectionKeys = append(sectionKeys, section)
		}

		sectionMap[key] = v
		fieldKeys = append(fieldKeys, key)
	}

	// Sort the keys to ensure a deterministic output
	sort.Strings(sectionKeys)
	sort.Strings(fieldKeys)

	tokens := make(ConfigTokens, 0, len(m))

	line := 1
	for _, section := range sectionKeys {
		tokens = append(tokens, ConfigToken{
			Type:  ConfigTokenTypeSection,
			Value: []byte(section),
			Line:  line,
		}, ConfigToken{
			Type:  ConfigTokenTypeSectionStart,
			Value: []byte{'{'},
			Line:  line,
		}, ConfigToken{
			Type:  ConfigTokenTypeEndOfStatement,
			Value: []byte{'\n'},
			Line:  line,
		})
		line++
		for _, key := range fieldKeys {
			value := mbySection[section][key]

			tokens = append(tokens, ConfigToken{
				Type:  ConfigTokenTypeIdentifier,
				Value: []byte(key),
				Line:  line,
			}, ConfigToken{
				Type:  ConfigTokenTypeAssignment,
				Value: []byte{'='},
				Line:  line,
			})

			switch value := value.(type) {
			case int64:
				tokens = append(tokens, ConfigToken{
					Type:  ConfigTokenTypeInt64,
					Value: []byte(fmt.Sprintf("%d", value)),
					Line:  line,
				})
			default:
				tokens = append(tokens, ConfigToken{
					Type:  ConfigTokenTypeString,
					Value: []byte(fmt.Sprintf("%v", value)),
					Line:  line,
				})
			}

			tokens = append(tokens, ConfigToken{
				Type:  ConfigTokenTypeEndOfStatement,
				Value: []byte{'\n'},
				Line:  line,
			})
			line++
		}

		tokens = append(tokens, ConfigToken{
			Type:  ConfigTokenTypeSectionEnd,
			Value: []byte{'}'},
			Line:  line,
		}, ConfigToken{
			Type:  ConfigTokenTypeEndOfStatement,
			Value: []byte{'\n'},
			Line:  line,
		})
		line++
	}

	return c.writeTokens(tokens)
}

func (c *configLexEncoder) writeTokens(tokens ConfigTokens) error {
	data, err := ConfigTokensToBytes(tokens)
	if err != nil {
		return fmt.Errorf("failed to write tokens into byte representation: %w", err)
	}

	if _, err := io.Copy(c.Writer, bytes.NewReader(data)); err != nil {
		return fmt.Errorf("failed to write tokens: %w", err)
	}
	return nil
}

var _ LexingConfigEncoder = &configLexEncoder{}

func ConfigTokensToBytes(tokens ConfigTokens) ([]byte, error) {
	// We can estimate a good buffer size by requesting 15% more than the minimum size for all values
	// This way we are accommodating for the fact that we might need to add spaces or tabs for indentation.
	expectedSize := int(float32(tokens.minimumSize()) * 1.15)
	buf := bytes.NewBuffer(make([]byte, 0, expectedSize))

	// writeTabOrSpaceIfInSection writes a tab if the line is in a section and has not been indented yet
	// otherwise it writes a space
	// we can use this for correct indentation of the configuration file
	inSection := false
	linesIndented := map[int]struct{}{}
	writeTabOrSpaceIfInSection := func(line int) {
		if inSection {
			if _, ok := linesIndented[line]; !ok {
				buf.WriteRune('\t')
				linesIndented[line] = struct{}{}
			} else {
				buf.WriteRune(' ')
			}
		}
	}

	for _, token := range tokens {
		switch token.Type {
		case ConfigTokenTypeComment:
			writeTabOrSpaceIfInSection(token.Line)
			buf.Write(token.Value)
			buf.WriteRune(' ') // readability: Add a space after the comment identifier
		case ConfigTokenTypeCommentValue:
			buf.Write(token.Value)
		case ConfigTokenTypeEndOfStatement:
			buf.Write(token.Value)
		case ConfigTokenTypeSection:
			buf.Write(token.Value)
			buf.WriteRune(' ') // readability: Add a space after the section identifier
		case ConfigTokenTypeSectionStart:
			buf.Write(token.Value)
			inSection = true
		case ConfigTokenTypeSectionEnd:
			buf.Write(token.Value)
			inSection = false
		case ConfigTokenTypeString:
			buf.WriteRune('"')
			buf.Write(token.Value)
			buf.WriteRune('"')
		case ConfigTokenTypeIdentifier:
			writeTabOrSpaceIfInSection(token.Line)
			buf.Write(token.Value)
			buf.WriteRune(' ') // readability: Add a space after the identifier
		case ConfigTokenTypeAssignment:
			buf.Write(token.Value)
			buf.WriteRune(' ') // readability: Add a space after the assignment
		case ConfigTokenTypeInt64:
			buf.Write(token.Value)
		case ConfigTokenTypeSOF:
			continue
		case ConfigTokenTypeEOF:
			break
		case ConfigTokenTypeError:
			return nil, token.Err
		case configTokenTypeNotYetKnown:
			return nil, fmt.Errorf("unexpected token type %v", token.Type)
		}
	}
	return buf.Bytes(), nil
}
