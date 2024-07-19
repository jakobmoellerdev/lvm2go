package main

import (
	"context"
	"errors"
	"log/slog"
	"os"

	. "github.com/jakobmoellerdev/lvm2go"
)

func main() {
	if os.Geteuid() != 0 {
		panic("panicking because lvm2 requires root privileges for most operations.")
	}
	if err := run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func run() (err error) {
	ctx := context.Background()
	lvm := NewClient()
	vgName := VolumeGroupName("test")
	lvName := LogicalVolumeName("test")
	deviceSize := MustParseSize("1G")
	lvSize := MustParseSize("100M")

	var losetup LoopbackDevice
	if losetup, err = NewLoopbackDevice(deviceSize); err != nil {
		return
	}
	defer func() {
		err = errors.Join(err, losetup.Close())
	}()

	if err = lvm.VGCreate(ctx, vgName, PhysicalVolumesFrom(losetup.Device())); err != nil {
		return
	}
	defer func() {
		err = errors.Join(err, lvm.VGRemove(ctx, vgName))
	}()

	if err = lvm.LVCreate(ctx, vgName, lvName, lvSize); err != nil {
		return
	}
	defer func() {
		err = errors.Join(err, lvm.LVRemove(ctx, vgName, lvName))
	}()

	return
}
