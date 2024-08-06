package lvm2go

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

// ConfigLexerReader is an interface for reading tokens from a configuration file
// The lexer reads the configuration file and returns ConfigTokens that can be used to
// decode the configuration file into a struct or do other operations.
// Any returned Token is one of ConfigTokenType, for more details see the ConfigTokenType documentation.
type ConfigLexerReader interface {
	// Lex reads the configuration file and returns all tokens in the file or an error if one occurs
	// The lexer will return an EOF token when the end of the file is reached.
	// The lexer will return an Error token when an error occurs.
	// Lex can be used to read the entire configuration file in one operation and to decouple reading from parsing.
	Lex() (ConfigTokens, error)

	// Next returns the next set of tokens in the configuration file that can be read in a single operation
	// Note that using Next() will not fail out if an error occurs, it will return the ConfigTokenError in the tokens
	// as it is considered part of the read operation.
	// The lexer will return an EOF token when the end of the file is reached.
	// The lexer will return an Error token when an error occurs.
	// Next can be used to implement efficient parsers that only read the next token when needed.
	Next() ConfigTokens
}

func NewBufferedConfigLexer(dataStream io.Reader) ConfigLexerReader {
	return &ConfigLexer{
		current:     ConfigTokenTypeSOF,
		dataStream:  bufio.NewReaderSize(dataStream, 4096),
		lineBuffer:  bytes.NewBuffer(make([]byte, 256)),
		currentLine: 1,
	}
}

type ConfigTokenType rune

const (
	// ConfigTokenTypeComment represents comments such as
	// # This is a comment
	ConfigTokenTypeComment ConfigTokenType = iota

	ConfigTokenTypeCommentValue ConfigTokenType = iota

	// ConfigTokenTypeEndOfStatement represents the end of a statement
	// This can be a newline.
	ConfigTokenTypeEndOfStatement ConfigTokenType = iota

	// ConfigTokenTypeSection represents a section name
	// Example:
	// config { ← This is a section named "config"
	//     key = value
	// }
	ConfigTokenTypeSection ConfigTokenType = iota

	// ConfigTokenTypeSectionStart represents the start of a section
	// Example:
	// config { ← This is a section start token "{"
	ConfigTokenTypeSectionStart ConfigTokenType = iota

	// ConfigTokenTypeSectionEnd represents the end of a section
	// Example:
	// config { ← This is a section end token "}"
	ConfigTokenTypeSectionEnd ConfigTokenType = iota

	// ConfigTokenTypeString represents a string
	// Example:
	// key = "value" ← This is a string token "value"
	ConfigTokenTypeString ConfigTokenType = iota

	// ConfigTokenTypeIdentifier represents an identifier
	// Example:
	// key = value ← This is an identifier token "key"
	ConfigTokenTypeIdentifier ConfigTokenType = iota

	// ConfigTokenTypeAssignment represents an assignment
	// Example:
	// key = value ← This is an assignment token "="
	ConfigTokenTypeAssignment ConfigTokenType = iota

	// ConfigTokenTypeInt64 represents an int64
	// Example:
	// key = 1234 ← This is an int64 token "1234"
	ConfigTokenTypeInt64 ConfigTokenType = iota

	// ConfigTokenTypeSOF represents the start of the file
	ConfigTokenTypeSOF ConfigTokenType = iota

	// ConfigTokenTypeEOF represents the end of the file
	ConfigTokenTypeEOF ConfigTokenType = iota

	ConfigTokenTypeError ConfigTokenType = iota

	// configTokenTypeNotYetKnown represents a token that has not yet been lexed
	configTokenTypeNotYetKnown ConfigTokenType = iota
)

func (t ConfigTokenType) String() string {
	switch t {
	case ConfigTokenTypeComment:
		return "Comment"
	case ConfigTokenTypeCommentValue:
		return "CommentValue"
	case ConfigTokenTypeEndOfStatement:
		return "EndOfStatement"
	case ConfigTokenTypeSection:
		return "Section"
	case ConfigTokenTypeSectionStart:
		return "SectionStart"
	case ConfigTokenTypeSectionEnd:
		return "SectionEnd"
	case ConfigTokenTypeString:
		return "String"
	case ConfigTokenTypeIdentifier:
		return "Identifier"
	case ConfigTokenTypeAssignment:
		return "Assignment"
	case ConfigTokenTypeInt64:
		return "Int64"
	case ConfigTokenTypeSOF:
		return "SOF"
	case ConfigTokenTypeEOF:
		return "EOF"
	case ConfigTokenTypeError:
		return "Error"
	default:
		return "Unknown"
	}
}

type ConfigLexer struct {
	// dataStream is the stream of data to be lexed
	dataStream *bufio.Reader

	// lineBuffer is a buffer to store the current line being lexed in case of lookbehind
	lineBuffer *bytes.Buffer

	current     ConfigTokenType
	readCount   int
	currentLine int
}

type ConfigTokens []ConfigToken

func (t ConfigTokens) String() string {
	builder := strings.Builder{}
	for _, token := range t {
		builder.WriteString(token.String())
		builder.WriteRune('\n')
	}
	return builder.String()
}

func (t ConfigTokens) minimumSize() int {
	size := 0
	for _, token := range t {
		size += len(token.Value)
	}
	return size
}

type ConfigToken struct {
	Type        ConfigTokenType
	Value       []byte
	Err         error
	Line, Start int
}

func (t ConfigToken) String() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%d:%d\t", t.Line, t.Start))
	builder.WriteString(t.Type.String())
	builder.WriteRune(' ')
	if t.Err != nil {
		builder.WriteString(t.Err.Error())
	} else {
		builder.WriteString(fmt.Sprintf("%q", t.Value))
	}
	return builder.String()
}

var ConfigTokenEOF = ConfigToken{Type: ConfigTokenTypeEOF, Start: -1, Line: -1}

func ConfigTokenError(err error) ConfigToken {
	return ConfigToken{Type: ConfigTokenTypeError, Err: err, Start: -1, Line: -1}
}

func (l *ConfigLexer) Lex() (ConfigTokens, error) {
	tokens := make(ConfigTokens, 0, 4)
	for {
		tokensFromNext := l.Next()
		tokens = append(tokens, tokensFromNext...)

		// If the next token is an EOF or an error, return the tokens
		for _, next := range tokensFromNext {
			if next.Type == ConfigTokenTypeEOF {
				return tokens, nil
			}
			if next.Type == ConfigTokenTypeError {
				return tokens, next.Err
			}
		}
	}
}

// Next returns the next token in the stream
func (l *ConfigLexer) Next() ConfigTokens {
	l.lineBuffer.Reset()
	for {
		candidate, size, err := l.dataStream.ReadRune()
		if err == io.EOF {
			return ConfigTokens{ConfigTokenEOF}
		}
		if err != nil {
			return ConfigTokens{ConfigTokenError(err)}
		}

		l.readCount += size

		tokenType := l.RuneToTokenType(candidate)
		if tokenType == configTokenTypeNotYetKnown {
			l.lineBuffer.WriteRune(candidate)
			continue
		}

		if tokenType == ConfigTokenTypeEndOfStatement {
			l.lineBuffer.Reset()
			l.currentLine++
			return ConfigTokens{{
				Type:  ConfigTokenTypeEndOfStatement,
				Value: runeToUTF8(candidate),
				Start: l.readCount,
				Line:  l.currentLine,
			}}
		}

		loc := l.readCount
		tokens := ConfigTokens{}

		switch tokenType {
		case ConfigTokenTypeComment:
			tokens = l.handleComment(candidate, loc)
		case ConfigTokenTypeSectionEnd:
			tokens = append(tokens, ConfigToken{
				Type:  ConfigTokenTypeSectionEnd,
				Value: runeToUTF8(candidate),
				Start: l.readCount,
				Line:  l.currentLine,
			})
		case ConfigTokenTypeSectionStart:
			tokens = l.handleSectionStart(candidate, loc)
		case ConfigTokenTypeAssignment:
			tokens = l.handleAssignment(candidate, loc)
		default:
			return ConfigTokens{ConfigTokenError(fmt.Errorf("unexpected token type %v", tokenType))}
		}

		return tokens
	}
}

func (l *ConfigLexer) handleComment(candidate rune, loc int) ConfigTokens {
	comment, err := l.dataStream.ReadBytes('\n')
	l.readCount += len(comment)
	trimmedComment := bytes.TrimSpace(comment)

	tokens := ConfigTokens{
		{
			Type:  ConfigTokenTypeComment,
			Value: runeToUTF8(candidate),
			Start: loc,
			Line:  l.currentLine,
		},
		{
			Type:  ConfigTokenTypeCommentValue,
			Value: trimmedComment,
			Start: loc + len(comment) - len(trimmedComment),
			Line:  l.currentLine,
		},
		{
			Type:  ConfigTokenTypeEndOfStatement,
			Value: runeToUTF8('\n'),
			Start: loc + len(comment),
			Line:  l.currentLine,
		},
	}

	if err == io.EOF {
		tokens = append(tokens, ConfigTokenEOF)
	} else if err != nil {
		tokens = append(tokens, ConfigTokenError(err))
	}
	l.lineBuffer.Reset()
	l.currentLine++

	return tokens
}

func (l *ConfigLexer) handleSectionStart(candidate rune, loc int) ConfigTokens {
	section := l.lineBuffer.Bytes()
	sectionTrimmed := bytes.TrimSpace(section)

	tokens := ConfigTokens{
		{
			Type:  ConfigTokenTypeSection,
			Value: bytes.Clone(sectionTrimmed),
			Start: loc - len(section),
			Line:  l.currentLine,
		},
		{
			Type:  ConfigTokenTypeSectionStart,
			Value: runeToUTF8(candidate),
			Start: loc,
			Line:  l.currentLine,
		},
	}

	return tokens
}

func (l *ConfigLexer) handleAssignment(candidate rune, loc int) ConfigTokens {
	identifier := bytes.TrimSpace(l.lineBuffer.Bytes())
	tokens := ConfigTokens{
		{
			Type:  ConfigTokenTypeIdentifier,
			Value: bytes.Clone(identifier),
			Start: loc - len(identifier) - 1,
			Line:  l.currentLine,
		},
		{
			Type:  ConfigTokenTypeAssignment,
			Value: runeToUTF8(candidate),
			Start: loc,
			Line:  l.currentLine,
		},
	}

	restOfLine, err := l.dataStream.ReadBytes('\n')
	l.readCount += len(restOfLine)

	if postCommentIdx := bytes.IndexRune(restOfLine, '#'); postCommentIdx != -1 {
		comment := restOfLine[postCommentIdx:]
		restOfLineWOComment := restOfLine[:postCommentIdx]

		valueToken := l.createValueToken(restOfLineWOComment, loc)
		commentStart := loc + len(restOfLineWOComment) + 1
		commentTrimmed := bytes.TrimSpace(bytes.Trim(comment, "#"))

		tokens = append(tokens,
			valueToken,
			ConfigToken{
				Type:  ConfigTokenTypeComment,
				Value: runeToUTF8('#'),
				Start: commentStart,
				Line:  l.currentLine,
			},
			ConfigToken{
				Type:  ConfigTokenTypeCommentValue,
				Value: bytes.TrimSpace(commentTrimmed),
				Start: commentStart + len(comment) - len(commentTrimmed) - 1,
				Line:  l.currentLine,
			},
		)
	} else {
		valueToken := l.createValueToken(restOfLine, loc)
		tokens = append(tokens, valueToken)
	}

	tokens = append(tokens,
		ConfigToken{
			Type:  ConfigTokenTypeEndOfStatement,
			Value: runeToUTF8('\n'),
			Line:  l.currentLine,
			Start: l.readCount,
		},
	)

	l.lineBuffer.Reset()
	l.currentLine++

	if err == io.EOF {
		tokens = append(tokens, ConfigTokenEOF)
	} else if err != nil {
		tokens = append(tokens, ConfigTokenError(err))
	}

	return tokens
}

func (l *ConfigLexer) createValueToken(line []byte, loc int) ConfigToken {
	sQidx := bytes.IndexByte(line, '"')
	lQidx := bytes.LastIndexByte(line, '"')
	var valueToken ConfigToken
	if sQidx == -1 && lQidx == -1 {
		trimmedLine := bytes.TrimSpace(line)
		valueToken = ConfigToken{
			Type:  ConfigTokenTypeInt64,
			Value: trimmedLine,
			Start: loc + len(line) - len(trimmedLine),
			Line:  l.currentLine,
		}
	} else {
		trimmedLine := bytes.TrimSpace(line[sQidx+1 : lQidx])
		valueToken = ConfigToken{
			Type:  ConfigTokenTypeString,
			Value: trimmedLine,
			Start: loc + len(line) - len(trimmedLine) - 2,
			Line:  l.currentLine,
		}
	}
	return valueToken
}

func (l *ConfigLexer) RuneToTokenType(r rune) ConfigTokenType {
	switch r {
	case '{':
		return ConfigTokenTypeSectionStart
	case '}':
		return ConfigTokenTypeSectionEnd
	case '=':
		return ConfigTokenTypeAssignment
	case '\n':
		return ConfigTokenTypeEndOfStatement
	case '#':
		return ConfigTokenTypeComment
	default:
		return configTokenTypeNotYetKnown
	}
}

func runeToUTF8(r rune) []byte {
	return runesToUTF8([]rune{r})
}

func runesToUTF8(rs []rune) []byte {
	size := 0
	for _, r := range rs {
		size += utf8.RuneLen(r)
	}

	bs := make([]byte, size)

	count := 0
	for _, r := range rs {
		count += utf8.EncodeRune(bs[count:], r)
	}

	return bs
}
