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
	"context"
	"log/slog"
	"testing"

	. "github.com/jakobmoellerdev/lvm2go"
)

func Test_Version(t *testing.T) {
	SkipOrFailTestIfNotRoot(t)
	slog.SetDefault(slog.New(NewContextPropagatingSlogHandler(NewTestingHandler(t))))
	slog.SetLogLoggerLevel(slog.LevelDebug)
	ctx := context.Background()
	clnt := GetTestClient(ctx)

	ver, err := clnt.Version(ctx)

	if err != nil {
		t.Fatalf("failed to get version: %v", err)
	}

	if ver.LVMVersion == "" {
		t.Fatalf("LVM Version is empty")
	}

	if ver.LVMBuild.IsZero() {
		t.Fatalf("LVM Build Date is zero")
	}

	if ver.LibraryVersion == "" {
		t.Fatalf("Library Version is empty")
	}

	if ver.LibraryBuild.IsZero() {
		t.Fatalf("Library Build Date is zero")
	}

	if ver.DriverVersion == "" {
		t.Fatalf("Driver Version is empty")
	}

	if len(ver.ConfigurationFlags) == 0 {
		t.Fatalf("Configuration Flags is empty")
	}
}
