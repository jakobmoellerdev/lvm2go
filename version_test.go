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
