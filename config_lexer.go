package lvm2go

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

type ConfigLexerReader interface {
	Lex() (ConfigTokens, error)
	Next() ConfigTokens
}

func NewConfigLexer(dataStream io.Reader) ConfigLexerReader {
	return &configLexer{
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

type configLexer struct {
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
		builder.Write(t.Value)
	}
	return builder.String()
}

var ConfigTokenEOF = ConfigToken{Type: ConfigTokenTypeEOF, Start: -1, Line: -1}

func ConfigTokenError(err error) ConfigToken {
	return ConfigToken{Type: ConfigTokenTypeError, Err: err, Start: -1, Line: -1}
}

func (l *configLexer) Lex() (ConfigTokens, error) {
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
func (l *configLexer) Next() ConfigTokens {
	l.lineBuffer.Reset()
	for {
		candidate, size, err := l.dataStream.ReadRune()
		if err == io.EOF {
			return ConfigTokens{ConfigTokenEOF}
		} else if err != nil {
			return ConfigTokens{ConfigTokenError(err)}
		}

		l.readCount += size

		tokenType := l.RuneToTokenType(candidate)
		// If the token is already known from the rune, return it
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

		tokens := make(ConfigTokens, 0, 4)

		loc := l.readCount

		switch tokenType {
		case ConfigTokenTypeComment:
			comment, err := l.dataStream.ReadBytes('\n')
			l.readCount += len(comment)

			trimmedComment := bytes.TrimSpace(comment)
			tokens = append(
				tokens,
				ConfigToken{
					Type:  ConfigTokenTypeComment,
					Value: runeToUTF8(candidate),
					Start: loc,
					Line:  l.currentLine,
				},
				ConfigToken{
					Type:  ConfigTokenTypeCommentValue,
					Value: trimmedComment,
					Start: loc + len(comment) - len(trimmedComment),
					Line:  l.currentLine,
				},
				ConfigToken{
					Type:  ConfigTokenTypeEndOfStatement,
					Value: runeToUTF8('\n'),
					Start: loc + len(comment),
					Line:  l.currentLine,
				},
			)

			if err == io.EOF {
				return append(tokens, ConfigTokenEOF)
			} else if err != nil {
				return append(tokens, ConfigTokenError(err))
			}
			l.lineBuffer.Reset()
			l.currentLine++
		case ConfigTokenTypeSectionEnd:
			tokens = append(
				tokens,
				ConfigToken{
					Type:  ConfigTokenTypeSectionEnd,
					Value: runeToUTF8(candidate),
					Start: l.readCount,
					Line:  l.currentLine,
				},
			)
		case ConfigTokenTypeSectionStart:
			section := l.lineBuffer.Bytes()
			sectionTrimmed := bytes.TrimSpace(section)
			tokens = append(
				tokens,
				ConfigToken{
					Type:  ConfigTokenTypeSection,
					Value: bytes.Clone(sectionTrimmed),
					Start: l.readCount - len(section),
					Line:  l.currentLine,
				},
				ConfigToken{
					Type:  ConfigTokenTypeSectionStart,
					Value: runeToUTF8(candidate),
					Start: l.readCount,
					Line:  l.currentLine,
				},
			)
		case ConfigTokenTypeAssignment:
			identifier := bytes.TrimSpace(l.lineBuffer.Bytes())
			tokens = append(
				tokens,
				ConfigToken{
					Type:  ConfigTokenTypeIdentifier,
					Value: bytes.Clone(identifier),
					Start: l.readCount - len(identifier) - 1,
					Line:  l.currentLine,
				},
				ConfigToken{
					Type:  ConfigTokenTypeAssignment,
					Value: runeToUTF8(candidate),
					Start: l.readCount,
					Line:  l.currentLine,
				},
			)

			restOfLine, err := l.dataStream.ReadBytes('\n')

			// If the rest of the line contains a comment, split it
			if postCommentIdx := bytes.IndexRune(restOfLine, '#'); postCommentIdx != -1 {
				comment := restOfLine[postCommentIdx:]
				restOfLineWOComment := restOfLine[:postCommentIdx]

				var valueToken ConfigToken
				sQidx := bytes.IndexByte(restOfLineWOComment, '"')
				lQidx := bytes.LastIndexByte(restOfLineWOComment, '"')
				if sQidx == -1 && lQidx == -1 {
					trimmedRestOfLine := bytes.TrimSpace(restOfLineWOComment)
					// If the comment is inside a string, ignore it
					valueToken = ConfigToken{
						Type:  ConfigTokenTypeInt64,
						Value: trimmedRestOfLine,
						Start: l.readCount + len(restOfLineWOComment) - len(trimmedRestOfLine),
						Line:  l.currentLine,
					}
				} else {
					trimmedRestOfLine := bytes.TrimSpace(restOfLineWOComment[sQidx:lQidx])
					valueToken = ConfigToken{
						Type:  ConfigTokenTypeString,
						Value: trimmedRestOfLine,
						Start: l.readCount + len(restOfLineWOComment) - len(trimmedRestOfLine) - 2,
						Line:  l.currentLine,
					}
				}

				commentStart := l.readCount + len(restOfLineWOComment) + 1
				commentTrimmed := bytes.TrimSpace(bytes.Trim(comment, "#"))
				tokens = append(
					tokens,
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
				var valueToken ConfigToken
				sQidx := bytes.IndexByte(restOfLine, '"')
				lQidx := bytes.LastIndexByte(restOfLine, '"')
				if sQidx == -1 && lQidx == -1 {
					trimmedRestOfLine := bytes.TrimSpace(restOfLine)
					valueToken = ConfigToken{
						Type:  ConfigTokenTypeInt64,
						Value: trimmedRestOfLine,
						Start: l.readCount + len(restOfLine) - len(trimmedRestOfLine),
						Line:  l.currentLine,
					}
				} else {
					trimmedRestOfLine := bytes.TrimSpace(restOfLine[sQidx+1 : lQidx])
					valueToken = ConfigToken{
						Type:  ConfigTokenTypeString,
						Value: trimmedRestOfLine,
						Start: l.readCount + len(restOfLine) - len(trimmedRestOfLine) - 2,
						Line:  l.currentLine,
					}
				}
				tokens = append(tokens, valueToken)
			}

			l.readCount += len(restOfLine)

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
				return append(tokens, ConfigTokenEOF)
			} else if err != nil {
				return append(tokens, ConfigTokenError(err))
			}
		default:
			err := fmt.Errorf("unexpected token type %v", tokenType)
			return append(tokens, ConfigTokenError(err))
		}

		return tokens
	}
}

func (l *configLexer) RuneToTokenType(r rune) ConfigTokenType {
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
