package main

import (
	"context"
	"log/slog"
	"os"

	. "github.com/jakobmoellerdev/lvm2go"
)

func main() {
	if os.Geteuid() != 0 {
		panic("panicking because it requires root privileges to setup its environment.")
	}

	ctx := context.Background()
	lvm := NewClient()

	deviceSize := MustParseSize("1G")

	losetup, err := NewLoopbackDevice(deviceSize)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer func() {
		if err := losetup.Close(); err != nil {
			slog.Error(err.Error())
		}
	}()

	pvs := PhysicalVolumesFrom(losetup.Device())
	vgName := VolumeGroupName("test")
	pool := MustNewThinPool(vgName, "pool")
	poolSize := MustParseExtents("100%FREE")
	lvName := LogicalVolumeName("test")
	lvSize := MustParseSize("100M").Virtual()

	if err := lvm.VGCreate(ctx, vgName, pvs); err != nil {
		slog.Error(err.Error())
		return
	}
	defer func() {
		if err := lvm.VGRemove(ctx, vgName); err != nil {
			slog.Error(err.Error())
		}
	}()

	if err := lvm.LVCreate(ctx, vgName, pool.LogicalVolumeName, poolSize, Thin(true), ZeroVolume); err != nil {
		slog.Error(err.Error())
		return
	}
	defer func() {
		if err := lvm.LVRemove(ctx, vgName, pool.LogicalVolumeName); err != nil {
			slog.Error(err.Error())
		}
	}()

	if err := lvm.LVCreate(ctx, pool, lvName, lvSize); err != nil {
		slog.Error(err.Error())
		return
	}
	defer func() {
		if err := lvm.LVRemove(ctx, pool, lvName); err != nil {
			slog.Error(err.Error())
		}
	}()

	if err := lvm.LVResize(ctx, vgName, lvName, MustParsePrefixedSize("+10M")); err != nil {
		slog.Error(err.Error())
		return
	}
}
