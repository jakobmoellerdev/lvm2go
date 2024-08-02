/*
 Copyright 2024 The lvm2go Authors.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

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

var _ slog.Handler = &ContextPropagatingSlogHandler{}

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

func (h TestingHandler) Enabled(_ context.Context, _ slog.Level) bool {
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

type ContextPropagatingSlogHandler struct {
	handler slog.Handler
}

// NewContextPropagatingSlogHandler returns a new slog.Handler that propagates context values as slog attributes.
// The handler is a wrapper around the provided handler.
func NewContextPropagatingSlogHandler(handler slog.Handler) slog.Handler {
	return &ContextPropagatingSlogHandler{
		handler: handler,
	}
}

func (h *ContextPropagatingSlogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *ContextPropagatingSlogHandler) Handle(ctx context.Context, record slog.Record) error {
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

func (h *ContextPropagatingSlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ContextPropagatingSlogHandler{h.handler.WithAttrs(attrs)}
}

func (h *ContextPropagatingSlogHandler) WithGroup(name string) slog.Handler {
	return &ContextPropagatingSlogHandler{h.handler.WithGroup(name)}
}
