package lvm2go

import (
	"context"
	"log/slog"
	"testing"
)

func Test_RawConfig(t *testing.T) {
	FailTestIfNotRoot(t)
	slog.SetDefault(slog.New(NewContextPropagatingSlogHandler(NewTestingHandler(t))))
	slog.SetLogLoggerLevel(slog.LevelDebug)
	ctx := context.Background()
	clnt := GetTestClient(ctx)

	ver, err := clnt.RawConfig(ctx, ConfigTypeFull)

	if err != nil {
		t.Fatalf("failed to get config: %v", err)
	}

	if len(ver) == 0 {
		t.Fatalf("RawConfig is empty")
	}
}
