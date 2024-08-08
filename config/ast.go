package config

import (
	"bufio"
	"fmt"
	"reflect"
	"strings"
)

type AST interface {
	Statements() []ASTStatement

	Parse(t Tokens)
	Tokens() Tokens
}

type ModifiableAST interface {
	AST
	Append(ASTStatement)
}

type ast struct {
	statements []ASTStatement
}

func (a *ast) Append(statement ASTStatement) {
	a.statements = append(a.statements, statement)
}

func (a *ast) Statements() []ASTStatement {
	return a.statements
}

type ASTStatement interface {
	Tokens() Tokens
}

type Newline interface {
	ASTStatement
}

type newline struct {
	*Token
}

func (n *newline) Tokens() Tokens {
	return Tokens{n.Token}
}

type Section interface {
	ASTStatement
	Name() string
	Children() []ASTStatement

	Append(ASTStatement)
}

type Comment interface {
	ASTStatement
	Value() string
}

type Assignment interface {
	ASTStatement
	Key() string
	Value() string
}

func NewAST() ModifiableAST {
	return &ast{}
}

type section struct {
	sectionToken *Token
	sectionStart *Token
	sectionEnd   *Token
	statements   []ASTStatement
}

func (s *section) Tokens() Tokens {
	tokens := make(Tokens, 0)
	tokens = append(tokens, s.sectionToken)
	tokens = append(tokens, s.sectionStart)
	tokens = append(tokens, NewNewline().Tokens()...)
	for _, statement := range s.statements {
		tokens = append(tokens, statement.Tokens()...)
	}
	tokens = append(tokens, s.sectionEnd)
	tokens = append(tokens, NewNewline().Tokens()...)
	return tokens
}

func (s *section) Append(statement ASTStatement) {
	s.statements = append(s.statements, statement)
}

func (s *section) Children() []ASTStatement {
	return s.statements
}

var _ Section = &section{}
var _ ASTStatement = &section{}

func (s *section) Name() string {
	return string(s.sectionToken.Value)
}

var _ Section = &section{}

type comment struct {
	indicator *Token
	value     *Token
}

func (c *comment) Tokens() Tokens {
	return append(Tokens{c.indicator, c.value}, NewNewline().Tokens()...)
}

func (c *comment) Value() string {
	return string(c.value.Value)
}

var _ Comment = &comment{}

type assignment struct {
	key        *Token
	assignment *Token
	value      *Token
}

func (a *assignment) Tokens() Tokens {
	return append(Tokens{a.key, a.assignment, a.value}, NewNewline().Tokens()...)
}

var _ Assignment = &assignment{}

func (a *assignment) Key() string {
	return string(a.key.Value)
}

func (a *assignment) Value() string {
	return string(a.value.Value)
}

func (a *ast) Parse(t Tokens) {
	a.statements = nil
	var currentSection *section
	add := func(statement ASTStatement) {
		if currentSection != nil {
			currentSection.statements = append(currentSection.statements, statement)
		} else {
			a.statements = append(a.statements, statement)
		}
	}
	for i, token := range t {
		if token.Type == TokenTypeSection {
			section := &section{sectionToken: token}
			a.statements = append(a.statements, section)
			currentSection = section
			continue
		}
		if token.Type == TokenTypeStartOfSection {
			currentSection.sectionStart = token
			continue
		}
		if token.Type == TokenTypeEndOfSection {
			currentSection.sectionEnd = token
			currentSection = nil
			continue
		}
		if token.Type == TokenTypeComment {
			comment := &comment{indicator: token}
			comment.value = t[i+1]
			add(comment)
			i++
			continue
		}
		if token.Type == TokenTypeAssignment {
			assignment := &assignment{assignment: token}
			assignment.key = t[i-1]
			assignment.value = t[i+1]
			add(assignment)
			i++
			continue
		}
	}
}

func (a *ast) Tokens() Tokens {
	tokens := make(Tokens, 0)
	for _, statement := range a.statements {
		tokens = append(tokens, statement.Tokens()...)
	}

	line := 1
	for _, token := range tokens {
		token.Line = line
		if token.Type == TokenTypeEndOfStatement {
			line++
		}
	}

	if len(tokens) > 0 {
		tokens = append(tokens, TokenEOF)
	}

	return tokens
}

func NewSection(name string) Section {
	return &section{
		sectionToken: &Token{Type: TokenTypeSection, Value: []byte(name)},
		sectionStart: &Token{Type: TokenTypeStartOfSection, Value: []byte{'{'}},
		sectionEnd:   &Token{Type: TokenTypeEndOfSection, Value: []byte{'}'}},
	}
}

type MultiLineComment interface {
	ASTStatement
}

type multiLineComment []Comment

func (m multiLineComment) Tokens() Tokens {
	tokens := make(Tokens, 0, len(m))
	for _, comment := range m {
		tokens = append(tokens, comment.Tokens()...)
	}
	return tokens
}

func NewMultiLineComment(value string) MultiLineComment {
	scanner := bufio.NewScanner(strings.NewReader(value))
	comments := make([]Comment, 0)
	for scanner.Scan() {
		comments = append(comments, NewComment(scanner.Text()))
	}
	return multiLineComment(comments)
}

func NewComment(value string) Comment {
	return &comment{
		indicator: &Token{Type: TokenTypeComment, Value: []byte{'#'}},
		value:     &Token{Type: TokenTypeCommentValue, Value: []byte(value)},
	}
}

func NewAssignmentFromString(key, value string) Assignment {
	return &assignment{
		key:        &Token{Type: TokenTypeIdentifier, Value: []byte(key)},
		assignment: &Token{Type: TokenTypeAssignment, Value: []byte{'='}},
		value:      &Token{Type: TokenTypeString, Value: []byte(value)},
	}
}

func NewAssignmentFromSpec(spec *LVMStructTagFieldMapping) Assignment {
	switch spec.Kind() {
	case reflect.Int64:
		return NewAssignmentInt64(spec.Name, spec.Value.Int())
	default:
		return NewAssignmentFromString(spec.Name, spec.Value.String())
	}
}

func NewAssignment(key string, v any) Assignment {
	switch v := v.(type) {
	case *LVMStructTagFieldMapping:
		return NewAssignmentFromSpec(v)
	case int64:
		return NewAssignmentInt64(key, v)
	case string:
		return NewAssignmentFromString(key, v)
	default:
		panic(fmt.Sprintf("unexpected type %T", v))
	}
}

func NewAssignmentInt64(key string, value int64) Assignment {
	return &assignment{
		key:        &Token{Type: TokenTypeIdentifier, Value: []byte(key)},
		assignment: &Token{Type: TokenTypeAssignment, Value: []byte{'='}},
		value:      &Token{Type: TokenTypeInt64, Value: []byte(fmt.Sprintf("%d", value))},
	}
}

func NewNewline() Newline {
	return &newline{&Token{Type: TokenTypeEndOfStatement, Value: []byte{'\n'}}}
}
