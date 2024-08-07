package lvm2go

import (
	"bytes"
	"errors"
	"slices"
)

const LVMWarningPrefix = "WARNING: "

var stdErrNewLine = []byte("\n")

type LVMStdErr interface {
	error
	Bytes() []byte
	Lines(trimPrefix bool) [][]byte
	Warnings() []Warning
}

// AsLVMStdErr returns the LVMStdErr from the error if it exists and a bool indicating if LVMStdErr is present or not.
func AsLVMStdErr(err error) (LVMStdErr, bool) {
	var lvmStdErr LVMStdErr
	ok := errors.As(err, &lvmStdErr)
	return lvmStdErr, ok
}

func NewLVMStdErr(stderr []byte) LVMStdErr {
	if len(stderr) == 0 {
		return nil
	}
	lines := bytes.Split(stderr, stdErrNewLine)
	for i := range lines {
		lines[i] = bytes.TrimSpace(lines[i])
	}
	// Remove empty lines, e.g. due to double newlines
	lines = slices.DeleteFunc(lines, func(line []byte) bool {
		return len(line) == 0
	})
	// Sort and compact the lines, removing duplicates
	slices.SortStableFunc(lines, func(a, b []byte) int {
		return bytes.Compare(a, b)
	})
	lines = slices.CompactFunc(lines, func(a, b []byte) bool {
		return bytes.Equal(a, b)
	})
	return &stdErr{lines: lines}
}

type stdErr struct {
	lines [][]byte
}

func (e *stdErr) Error() string {
	return string(e.Bytes())
}

func (e *stdErr) Bytes() []byte {
	return bytes.Join(e.lines, stdErrNewLine)
}

func (e *stdErr) Lines(trimPrefix bool) [][]byte {
	if trimPrefix {
		trimmed := make([][]byte, len(e.lines))
		for i, line := range e.lines {
			trimmed[i] = bytes.TrimPrefix(line, []byte(LVMWarningPrefix))
		}
		return trimmed
	}
	return e.lines
}

func (e *stdErr) ExcludeWarnings() *stdErr {
	return &stdErr{lines: slices.DeleteFunc(e.lines, func(line []byte) bool {
		return bytes.HasPrefix(line, []byte(LVMWarningPrefix))
	})}
}

func (e *stdErr) Warnings() []Warning {
	warnings := make([]Warning, 0, len(e.lines))
	for _, line := range e.lines {
		if warning := NewWarning(line); warning != nil {
			warnings = append(warnings, warning)
		}
	}
	return warnings
}

type Warning interface {
	error
}

func NewWarning(raw []byte) Warning {
	if idx := bytes.LastIndex(raw, []byte(LVMWarningPrefix)); idx > 0 {
		return &warning{msg: raw[idx+len(LVMWarningPrefix):]}
	}
	return nil
}

type warning struct {
	msg []byte
}

func (w *warning) Error() string {
	return string(w.msg)
}
