package lvm2go

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

type contextKey string

var (
	fields contextKey = "slog_fields"
)

func WithValue(parent context.Context, key string, val any) context.Context {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	if v, ok := parent.Value(fields).(*sync.Map); ok {
		mapCopy := copySyncMap(v)
		mapCopy.Store(key, val)
		return context.WithValue(parent, fields, mapCopy)
	}
	v := &sync.Map{}
	v.Store(key, val)
	return context.WithValue(parent, fields, v)
}

func copySyncMap(m *sync.Map) *sync.Map {
	var cp sync.Map
	m.Range(func(k, v interface{}) bool {
		cp.Store(k, v)
		return true
	})
	return &cp
}

var _ slog.Handler = SlogHandler{}

type TestingHandler struct {
	tb         testing.TB
	extraAttrs []slog.Attr
	group      string
	format     string
}

func NewTestingHandler(tb testing.TB) slog.Handler {
	return TestingHandler{
		tb:     tb,
		format: time.RFC3339Nano,
		group:  tb.Name(),
	}
}

func (h TestingHandler) Enabled(_ context.Context, l slog.Level) bool {
	return !h.tb.Skipped()
}

func (h TestingHandler) Handle(_ context.Context, record slog.Record) error {
	if testing.Verbose() && record.Level < slog.LevelInfo {
		return nil
	}

	attributes := make([]string, 0, record.NumAttrs()+len(h.extraAttrs))
	for _, attr := range h.extraAttrs {
		attributes = append(attributes, fmt.Sprintf("%s=%v", attr.Key, attr.Value))
	}
	record.Attrs(func(attr slog.Attr) bool {
		attributes = append(attributes, fmt.Sprintf("%s=%v", attr.Key, attr.Value))
		return true
	})

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s %s", record.Time.Format(h.format), record.Level))

	if h.group != "" {
		sb.WriteString(fmt.Sprintf(" %s", h.group))
	}

	sb.WriteString(fmt.Sprintf(" %s %s", record.Message, strings.Join(attributes, " ")))

	if testing.Verbose() && h.printStack() {
		sb.WriteRune('\n')
		frames := runtime.CallersFrames([]uintptr{record.PC})
		frame, _ /* we know there is only one frame */ := frames.Next()
		sb.WriteString(fmt.Sprintf("%s:%d %s", frame.File, frame.Line, frame.Function))
	}

	h.tb.Helper()
	h.tb.Log(sb.String())

	return nil
}

func (h TestingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h.extraAttrs = append(h.extraAttrs, attrs...)
	return h
}

func (h TestingHandler) WithGroup(name string) slog.Handler {
	h.group = name
	return h
}

func (h TestingHandler) printStack() bool {
	return h.tb.Failed()
}

type SlogHandler struct {
	handler slog.Handler
}

func NewContextPropagatingSlogHandler(handler slog.Handler) slog.Handler {
	return SlogHandler{
		handler: handler,
	}
}

func (h SlogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h SlogHandler) Handle(ctx context.Context, record slog.Record) error {
	if v, ok := ctx.Value(fields).(*sync.Map); ok {
		v.Range(func(key, val any) bool {
			if keyString, ok := key.(string); ok {
				record.AddAttrs(slog.Any(keyString, val))
			}
			return true
		})
	}
	switch h.handler.(type) {
	case TestingHandler:
		h.handler.(TestingHandler).tb.Helper()
	}
	return h.handler.Handle(ctx, record)
}

func (h SlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return SlogHandler{h.handler.WithAttrs(attrs)}
}

func (h SlogHandler) WithGroup(name string) slog.Handler {
	return h.handler.WithGroup(name)
}
