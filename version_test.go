package lvm2go

import (
	"context"
	"log/slog"
	"testing"
)

func Test_Version(t *testing.T) {
	FailTestIfNotRoot(t)
	slog.SetDefault(slog.New(NewContextPropagatingSlogHandler(NewTestingHandler(t))))
	slog.SetLogLoggerLevel(slog.LevelDebug)
	ctx := context.Background()
	clnt := GetTestClient(ctx)

	ver, err := clnt.Version(ctx)

	if err != nil {
		t.Fatalf("failed to get version: %v", err)
	}

	t.Logf("LVM Version: %s", ver.LVMVersion)
	t.Logf("LVM Build Date: %s", ver.LVMBuild)
	t.Logf("Library Version: %s", ver.LibraryVersion)
	t.Logf("Library Build Date: %s", ver.LibraryBuild)
	t.Logf("Driver Version: %s", ver.DriverVersion)
	t.Logf("Configuration: %s", ver.ConfigurationFlags)
}
