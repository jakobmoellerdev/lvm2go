package lvm2go

import (
	"fmt"
	"strings"
)

type Select string

func (opt Select) ApplyToLVsOptions(opts *LVsOptions) {
	opts.Select = opt
}
func (opt Select) ApplyToVGsOptions(opts *VGsOptions) {
	opts.Select = opt
}
func (opt Select) ApplyToVGRemoveOptions(opts *VGRemoveOptions) {
	opts.Select = opt
}
func (opt Select) ApplyToLVRemoveOptions(opts *LVRemoveOptions) {
	opts.Select = opt
}

func NewMatchesAllSelector(fields map[string]string) Select {
	return NewSelector(AllFieldsMatch, Match, fields)
}

func NewMatchesAnySelector(fields map[string]string) Select {
	return NewSelector(AtLeastOneFieldMatches, Match, fields)
}

func NewSelector(
	lo LogicalAndGroupingOperator,
	co SelectionComparisonOperator,
	fields map[string]string,
) Select {
	var sb strings.Builder
	last := len(fields) - 1
	for field, value := range fields {
		last--
		sb.WriteRune('"')
		sb.WriteString(field)
		sb.WriteRune('"')
		sb.WriteString(string(co))
		sb.WriteRune('"')
		sb.WriteString(value)
		sb.WriteRune('"')
		if last > 0 {
			sb.WriteRune(' ')
			sb.WriteString(string(lo))
			sb.WriteRune(' ')
		}
	}
	return Select(sb.String())
}

func NewMatchesAllSelect(selects ...Select) Select {
	return NewCombinedSelect(AllFieldsMatch, selects...)
}

func NewMatchesAnySelect(selects ...Select) Select {
	return NewCombinedSelect(AtLeastOneFieldMatches, selects...)
}

func NotSelect(sel Select) Select {
	return Select(fmt.Sprintf("%s%s%s%s",
		string(LogicalNegation),
		string(LeftParenthesis),
		string(sel),
		string(RightParenthesis),
	))
}

func NewCombinedSelect(operator LogicalAndGroupingOperator, selects ...Select) Select {
	if len(selects) == 1 {
		return selects[0]
	}

	var sb strings.Builder
	last := len(selects) - 1
	for i, sel := range selects {
		sb.WriteString(string(LeftParenthesis))
		sb.WriteString(string(sel))
		sb.WriteString(string(RightParenthesis))
		if i < last {
			sb.WriteString(string(operator))
		}
	}
	return Select(sb.String())
}

type SelectionOperator string

type SelectionComparisonOperator SelectionOperator

const (
	MatchRegex    SelectionComparisonOperator = "=~"
	NotMatchRegex SelectionComparisonOperator = "!~"
	Match         SelectionComparisonOperator = "="
	NotMatch      SelectionComparisonOperator = "!="
	GreaterOrEq   SelectionComparisonOperator = ">="
	Greater       SelectionComparisonOperator = ">"
	LessOrEq      SelectionComparisonOperator = "<="
	Less          SelectionComparisonOperator = "<"
	Since         SelectionComparisonOperator = "since"
	After         SelectionComparisonOperator = "after"
	Until         SelectionComparisonOperator = "until"
	Before        SelectionComparisonOperator = "before"
)

type LogicalAndGroupingOperator SelectionOperator

const (
	AllFieldsMatch            LogicalAndGroupingOperator = "&&"
	AllFieldsMatchAlt         LogicalAndGroupingOperator = ","
	AtLeastOneFieldMatches    LogicalAndGroupingOperator = "||"
	AtLeastOneFieldMatchesAlt LogicalAndGroupingOperator = "#"
	LogicalNegation           LogicalAndGroupingOperator = "!"
	RightParenthesis          LogicalAndGroupingOperator = ")"
	LeftParenthesis           LogicalAndGroupingOperator = "("
	ListStart                 LogicalAndGroupingOperator = "["
	ListEnd                   LogicalAndGroupingOperator = "]"
	ListSubsetStart           LogicalAndGroupingOperator = "{"
	ListSubsetEnd             LogicalAndGroupingOperator = "}"
)
