package lvm2go

import (
	"context"
	"log/slog"
	"os"
	"sync"
)

var (
	isContainerized     bool
	detectContainerized sync.Once
)

func IsContainerized(ctx context.Context) bool {
	detectContainerized.Do(func() {
		if _, err := os.Stat("/.dockerenv"); err == nil {
			isContainerized = true
		} else if _, err := os.Stat("/.containerenv"); err == nil {
			isContainerized = true
		}
		if isContainerized {
			slog.DebugContext(ctx, "lvm2go is running in docker environment")
		}
	})
	return isContainerized
}
