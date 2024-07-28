package lvm2go_test

import (
	"context"
	"log/slog"
	"testing"

	. "github.com/jakobmoellerdev/lvm2go"
)

func Test_RawConfig(t *testing.T) {
	t.Parallel()
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
