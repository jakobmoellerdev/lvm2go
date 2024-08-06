package config

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/jakobmoellerdev/lvm2go/util"
)

// LexerReader is an interface for reading tokens from a configuration file
// The lexer reads the configuration file and returns ConfigTokens that can be used to
// decode the configuration file into a struct or do other operations.
// Any returned Token is one of TokenType, for more details see the TokenType documentation.
type LexerReader interface {
	// Lex reads the configuration file and returns all tokens in the file or an error if one occurs
	// The lexer will return an EOF token when the end of the file is reached.
	// The lexer will return an Error token when an error occurs.
	// Lex can be used to read the entire configuration file in one operation and to decouple reading from parsing.
	Lex() (Tokens, error)

	// Next returns the next set of tokens in the configuration file that can be read in a single operation
	// Note that using Next() will not fail out if an error occurs, it will return the ConfigTokenError in the tokens
	// as it is considered part of the read operation.
	// The lexer will return an EOF token when the end of the file is reached.
	// The lexer will return an Error token when an error occurs.
	// Next can be used to implement efficient parsers that only read the next token when needed.
	Next() Tokens
}

func NewBufferedLexer(dataStream io.Reader) LexerReader {
	return &Lexer{
		current:     TokenTypeSOF,
		dataStream:  bufio.NewReaderSize(dataStream, 4096),
		lineBuffer:  bytes.NewBuffer(make([]byte, 256)),
		currentLine: 1,
	}
}

type TokenType rune

const (
	// TokenTypeComment represents comments such as
	// # This is a comment
	TokenTypeComment TokenType = iota

	TokenTypeCommentValue TokenType = iota

	// TokenTypeEndOfStatement represents the end of a statement
	// This can be a newline.
	TokenTypeEndOfStatement TokenType = iota

	// TokenTypeSection represents a section name
	// Example:
	// config { ← This is a section named "config"
	//     key = value
	// }
	TokenTypeSection TokenType = iota

	// TokenTypeStartOfSection represents the start of a section
	// Example:
	// config { ← This is a section start token "{"
	TokenTypeStartOfSection TokenType = iota

	// TokenTypeEndOfSection represents the end of a section
	// Example:
	// config { ← This is a section end token "}"
	TokenTypeEndOfSection TokenType = iota

	// TokenTypeString represents a string
	// Example:
	// key = "value" ← This is a string token "value"
	TokenTypeString TokenType = iota

	// TokenTypeIdentifier represents an identifier
	// Example:
	// key = value ← This is an identifier token "key"
	TokenTypeIdentifier TokenType = iota

	// TokenTypeAssignment represents an assignment
	// Example:
	// key = value ← This is an assignment token "="
	TokenTypeAssignment TokenType = iota

	// TokenTypeInt64 represents an int64
	// Example:
	// key = 1234 ← This is an int64 token "1234"
	TokenTypeInt64 TokenType = iota

	// TokenTypeSOF represents the start of the file
	TokenTypeSOF TokenType = iota

	// TokenTypeEOF represents the end of the file
	TokenTypeEOF TokenType = iota

	TokenTypeError TokenType = iota

	// TokenTypeNotYetKnown represents a token that has not yet been lexed
	TokenTypeNotYetKnown TokenType = iota
)

func (t TokenType) String() string {
	switch t {
	case TokenTypeComment:
		return "Comment"
	case TokenTypeCommentValue:
		return "CommentValue"
	case TokenTypeEndOfStatement:
		return "EndOfStatement"
	case TokenTypeSection:
		return "Section"
	case TokenTypeStartOfSection:
		return "SectionStart"
	case TokenTypeEndOfSection:
		return "SectionEnd"
	case TokenTypeString:
		return "String"
	case TokenTypeIdentifier:
		return "Identifier"
	case TokenTypeAssignment:
		return "Assignment"
	case TokenTypeInt64:
		return "Int64"
	case TokenTypeSOF:
		return "SOF"
	case TokenTypeEOF:
		return "EOF"
	case TokenTypeError:
		return "Error"
	default:
		return "Unknown"
	}
}

type Lexer struct {
	// dataStream is the stream of data to be lexed
	dataStream *bufio.Reader

	// lineBuffer is a buffer to store the current line being lexed in case of lookbehind
	lineBuffer *bytes.Buffer

	current     TokenType
	readCount   int
	currentLine int
}

type Tokens []*Token

func (t Tokens) String() string {
	builder := strings.Builder{}
	for _, token := range t {
		builder.WriteString(token.String())
		builder.WriteRune('\n')
	}
	return builder.String()
}

func (t Tokens) minimumSize() int {
	size := 0
	for _, token := range t {
		size += len(token.Value)
	}
	return size
}

type configTokensByIdentifier map[string]Tokens

func AssignmentsWithSections(t Tokens) configTokensByIdentifier {
	sectionIndex := -1
	assignments := make(map[string]Tokens)
	for i, token := range t {
		if token.Type == TokenTypeSection {
			sectionIndex = i
			continue
		}

		if token.Type != TokenTypeAssignment {
			continue
		}

		assignments[string(t[i-1].Value)] = Tokens{
			t[sectionIndex],
			t[i-1],
			token,
			t[i+1],
		}
	}
	return assignments
}

func (a configTokensByIdentifier) OverrideWith(other configTokensByIdentifier) (notFound Tokens) {
	for key, value := range other {
		v, ok := a[key]
		if !ok {
			notFound = append(notFound, value...)
		} else {
			v[3].Value = value[3].Value
		}
	}
	return
}

func AppendAssignmentsAtEndOfSections(into Tokens, toAdd Tokens) Tokens {
	section := ""
	tokens := Tokens{}
	for i, token := range into {
		tokens = append(tokens, token)
		if token.Type == TokenTypeSection {
			section = string(token.Value)
			continue
		}
		if token.Type == TokenTypeEndOfSection {
			candidates := Tokens{}
			for j, token := range toAdd {
				if token.Type != TokenTypeAssignment {
					continue
				}
				inSection := section != ""
				isID := toAdd[j-1].Type == TokenTypeIdentifier
				isSection := toAdd[j-2].Type == TokenTypeSection
				if inSection && isID && isSection && section == string(toAdd[j-2].Value) {
					candidates = append(candidates,
						&Token{
							Type:  TokenTypeComment,
							Value: runeToUTF8('#'),
						},
						&Token{
							Type:  TokenTypeCommentValue,
							Value: []byte(generateLVMConfigEditComment()),
						},
						&Token{
							Type:  TokenTypeEndOfStatement,
							Value: runeToUTF8('\n'),
						},
						toAdd[j-1], token, toAdd[j+1],
						&Token{
							Type:  TokenTypeEndOfStatement,
							Value: runeToUTF8('\n'),
						})
				}
			}

			tokens = append(tokens[:i], append(candidates, tokens[i:]...)...)
		}
	}
	return tokens
}

func (t Tokens) InSection(section string) Tokens {
	tokensInSection := Tokens{}
	for _, token := range t {
		if token.Type == TokenTypeSection {
			if inSection := string(token.Value) == section; inSection {
				continue
			}
		}
		if token.Type == TokenTypeStartOfSection {
			continue
		}
		if token.Type == TokenTypeEndOfSection {
			break
		}
		tokensInSection = append(tokensInSection, token)
	}
	return tokensInSection
}

type Token struct {
	Type        TokenType
	Value       []byte
	Err         error
	Line, Start int
}

func (t Token) String() string {
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

var TokenEOF = &Token{Type: TokenTypeEOF, Start: -1, Line: -1}

func TokenError(err error) *Token {
	return &Token{Type: TokenTypeError, Err: err, Start: -1, Line: -1}
}

func (l *Lexer) Lex() (Tokens, error) {
	tokens := make(Tokens, 0, 4)
	for {
		tokensFromNext := l.Next()
		tokens = append(tokens, tokensFromNext...)

		// If the next token is an EOF or an error, return the tokens
		for _, next := range tokensFromNext {
			if next.Type == TokenTypeEOF {
				return tokens, nil
			}
			if next.Type == TokenTypeError {
				return tokens, next.Err
			}
		}
	}
}

// Next returns the next token in the stream
func (l *Lexer) Next() Tokens {
	l.lineBuffer.Reset()
	for {
		candidate, size, err := l.dataStream.ReadRune()
		if err == io.EOF {
			return Tokens{TokenEOF}
		}
		if err != nil {
			return Tokens{TokenError(err)}
		}

		l.readCount += size

		tokenType := l.RuneToTokenType(candidate)
		if tokenType == TokenTypeNotYetKnown {
			l.lineBuffer.WriteRune(candidate)
			continue
		}

		if tokenType == TokenTypeEndOfStatement {
			l.lineBuffer.Reset()
			l.currentLine++
			return Tokens{{
				Type:  TokenTypeEndOfStatement,
				Value: runeToUTF8(candidate),
				Start: l.readCount,
				Line:  l.currentLine,
			}}
		}

		loc := l.readCount
		tokens := Tokens{}

		switch tokenType {
		case TokenTypeComment:
			tokens = l.newComment(candidate, loc)
		case TokenTypeEndOfSection:
			tokens = append(tokens, &Token{
				Type:  TokenTypeEndOfSection,
				Value: runeToUTF8(candidate),
				Start: l.readCount,
				Line:  l.currentLine,
			})
		case TokenTypeStartOfSection:
			tokens = l.newSectionStart(candidate, loc)
		case TokenTypeAssignment:
			tokens = l.newAssignment(candidate, loc)
		default:
			return Tokens{TokenError(fmt.Errorf("unexpected token type %v", tokenType))}
		}

		return tokens
	}
}

func (l *Lexer) newComment(candidate rune, loc int) Tokens {
	comment, err := l.dataStream.ReadBytes('\n')
	l.readCount += len(comment)
	trimmedComment := bytes.TrimSpace(comment)

	tokens := Tokens{
		{
			Type:  TokenTypeComment,
			Value: runeToUTF8(candidate),
			Start: loc,
			Line:  l.currentLine,
		},
		{
			Type:  TokenTypeCommentValue,
			Value: trimmedComment,
			Start: loc + len(comment) - len(trimmedComment),
			Line:  l.currentLine,
		},
		{
			Type:  TokenTypeEndOfStatement,
			Value: runeToUTF8('\n'),
			Start: loc + len(comment),
			Line:  l.currentLine,
		},
	}

	if err == io.EOF {
		tokens = append(tokens, TokenEOF)
	} else if err != nil {
		tokens = append(tokens, TokenError(err))
	}
	l.lineBuffer.Reset()
	l.currentLine++

	return tokens
}

func (l *Lexer) newSectionStart(candidate rune, loc int) Tokens {
	section := l.lineBuffer.Bytes()
	sectionTrimmed := bytes.TrimSpace(section)

	tokens := Tokens{
		{
			Type:  TokenTypeSection,
			Value: bytes.Clone(sectionTrimmed),
			Start: loc - len(section),
			Line:  l.currentLine,
		},
		{
			Type:  TokenTypeStartOfSection,
			Value: runeToUTF8(candidate),
			Start: loc,
			Line:  l.currentLine,
		},
	}

	return tokens
}

func (l *Lexer) newAssignment(candidate rune, loc int) Tokens {
	identifier := bytes.TrimSpace(l.lineBuffer.Bytes())
	tokens := Tokens{
		{
			Type:  TokenTypeIdentifier,
			Value: bytes.Clone(identifier),
			Start: loc - len(identifier) - 1,
			Line:  l.currentLine,
		},
		{
			Type:  TokenTypeAssignment,
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
			&Token{
				Type:  TokenTypeComment,
				Value: runeToUTF8('#'),
				Start: commentStart,
				Line:  l.currentLine,
			},
			&Token{
				Type:  TokenTypeCommentValue,
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
		&Token{
			Type:  TokenTypeEndOfStatement,
			Value: runeToUTF8('\n'),
			Line:  l.currentLine,
			Start: l.readCount,
		},
	)

	l.lineBuffer.Reset()
	l.currentLine++

	if err == io.EOF {
		tokens = append(tokens, TokenEOF)
	} else if err != nil {
		tokens = append(tokens, TokenError(err))
	}

	return tokens
}

func (l *Lexer) createValueToken(line []byte, loc int) *Token {
	sQidx := bytes.IndexByte(line, '"')
	lQidx := bytes.LastIndexByte(line, '"')
	var valueToken Token
	if sQidx == -1 && lQidx == -1 {
		trimmedLine := bytes.TrimSpace(line)
		valueToken = Token{
			Type:  TokenTypeInt64,
			Value: trimmedLine,
			Start: loc + len(line) - len(trimmedLine),
			Line:  l.currentLine,
		}
	} else {
		trimmedLine := bytes.TrimSpace(line[sQidx+1 : lQidx])
		valueToken = Token{
			Type:  TokenTypeString,
			Value: trimmedLine,
			Start: loc + len(line) - len(trimmedLine) - 2,
			Line:  l.currentLine,
		}
	}
	return &valueToken
}

func (l *Lexer) RuneToTokenType(r rune) TokenType {
	switch r {
	case '{':
		return TokenTypeStartOfSection
	case '}':
		return TokenTypeEndOfSection
	case '=':
		return TokenTypeAssignment
	case '\n':
		return TokenTypeEndOfStatement
	case '#':
		return TokenTypeComment
	default:
		return TokenTypeNotYetKnown
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

// generateLVMConfigEditComment generates a comment to be added to the configuration file
// This comment is used to indicate that the field was edited by the client.
func generateLVMConfigEditComment() string {
	return fmt.Sprintf(`This field was edited by %s at %s`, util.ModuleID(), time.Now().Format(time.RFC3339))
}
