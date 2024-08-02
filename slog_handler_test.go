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

package lvm2go_test

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"

	. "github.com/jakobmoellerdev/lvm2go"
)

func TestNewContextPropagatingSlogHandler(t *testing.T) {
	SkipOrFailTestIfNotRoot(t)
	stdout := &bytes.Buffer{}
	var loggingHandler slog.Handler
	loggingHandler = slog.NewJSONHandler(stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})
	loggingHandler = NewContextPropagatingSlogHandler(loggingHandler)
	slog.SetDefault(slog.New(loggingHandler))
	ctx := context.Background()

	lvm := NewClient()
	if _, err := lvm.Version(ctx); err != nil {
		t.Errorf("Error: %v", err)
	}

	lineReader := bufio.NewScanner(stdout)
	var lines []map[string]any
	for lineReader.Scan() {
		line := make(map[string]any)
		if err := json.NewDecoder(bytes.NewReader(lineReader.Bytes())).Decode(&line); err != nil {
			t.Errorf("Error: %v", err)
		}
		lines = append(lines, line)
	}
	if err := lineReader.Err(); err != nil {
		t.Errorf("Error: %v", err)
	}
	if len(lines) == 0 {
		t.Errorf("Expected output in logger, got nothing")
	}

	foundcommand := false
	for _, line := range lines {
		if line["command"] != nil && line["msg"] != nil {
			foundcommand = true
			break
		}
	}
	if !foundcommand {
		t.Errorf("Expected command in logger output, got nothing")
	}
}
