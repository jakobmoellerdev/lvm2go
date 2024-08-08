package config

import (
	"bytes"
	"fmt"
	"io"
	"slices"
	"strings"
)

type LexingEncoder interface {
	Encode(v any) error
}

type StructuredLexingEncoder interface {
	EncodeStructured(v any) error
}

type UnstructuredLexingEncoder interface {
	EncodeUnstructured(v any) error
}

func NewLexingEncoder(writer io.Writer) LexingEncoder {
	return &lexEncoder{
		Writer: writer,
	}
}

type lexEncoder struct {
	Writer io.Writer
}

func (c *lexEncoder) Encode(v any) error {
	switch v := v.(type) {
	case Tokens:
		return c.writeTokens(v)
	}
	if isUnstructuredMap(v) {
		return c.EncodeUnstructured(v)
	}
	return c.EncodeStructured(v)
}

func (c *lexEncoder) EncodeStructured(v any) error {
	fieldSpecs, err := DecodeFieldMappings(v)
	if err != nil {
		return err
	}
	tokens := FieldMappingsToTokens(fieldSpecs)
	return c.writeTokens(tokens)
}

func (c *lexEncoder) EncodeUnstructured(v any) error {
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
		for _, key := range fieldKeys {
			value := mbySection[section][key]
			astSection.Append(NewAssignment(key, value))
		}
	}

	return c.writeTokens(ast.Tokens())
}

func (c *lexEncoder) writeTokens(tokens Tokens) error {
	data, err := TokensToBytes(tokens)
	if err != nil {
		return fmt.Errorf("failed to write tokens into byte representation: %w", err)
	}

	if _, err := io.Copy(c.Writer, bytes.NewReader(data)); err != nil {
		return fmt.Errorf("failed to write tokens: %w", err)
	}
	return nil
}

var _ LexingEncoder = &lexEncoder{}

func TokensToBytes(tokens Tokens) ([]byte, error) {
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
		case TokenTypeComment:
			writeTabOrSpaceIfInSection(token.Line)
			buf.Write(token.Value)
			buf.WriteRune(' ') // readability: Insert a space after the comment identifier
		case TokenTypeCommentValue:
			buf.Write(token.Value)
		case TokenTypeEndOfStatement:
			buf.Write(token.Value)
		case TokenTypeSection:
			buf.Write(token.Value)
			buf.WriteRune(' ') // readability: Insert a space after the section identifier
		case TokenTypeStartOfSection:
			buf.Write(token.Value)
			inSection = true
		case TokenTypeEndOfSection:
			buf.Write(token.Value)
			inSection = false
		case TokenTypeString:
			buf.WriteRune('"')
			buf.Write(token.Value)
			buf.WriteRune('"')
		case TokenTypeIdentifier:
			writeTabOrSpaceIfInSection(token.Line)
			buf.Write(token.Value)
			buf.WriteRune(' ') // readability: Insert a space after the identifier
		case TokenTypeAssignment:
			buf.Write(token.Value)
			buf.WriteRune(' ') // readability: Insert a space after the assignment
		case TokenTypeInt64:
			buf.Write(token.Value)
		case TokenTypeSOF:
			continue
		case TokenTypeEOF:
			break
		case TokenTypeError:
			return nil, token.Err
		case TokenTypeNotYetKnown:
			return nil, fmt.Errorf("unexpected token type %v", token.Type)
		}
	}
	return buf.Bytes(), nil
}
