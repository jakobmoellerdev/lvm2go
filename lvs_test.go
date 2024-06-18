package lvm2go

import (
	"context"
	"log/slog"
	"os"
	"testing"
)

func TestLVs(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	SkipTestIfNotRoot(t)

	loop := MakeTestLoopbackDevice(t, "1G")
	volumeGroup := MakeTestVolumeGroup(t, loop.Device)
	logicalVolume := volumeGroup.MakeTestLogicalVolume("100M")

	lvs, err := LVs(context.Background(), LVsOptions{
		VolumeGroupName: volumeGroup.Name,
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(lvs) != 1 {
		t.Fatalf("Expected 1 logical volume, got %d", len(lvs))
	}

	if lvs[0].Name != string(logicalVolume) {
		t.Fatalf("Expected logical volume name to be %s, got %s", logicalVolume, lvs[0].Name)
	}

	if lvs[0].Size != 100*1024*1024 {
		t.Fatalf("Expected logical volume size to be 100M, got %d", lvs[0].Size)
	}

	if lvs[0].VolumeGroupName != volumeGroup.Name {
		t.Fatalf("Expected volume group name to be %s, got %s", volumeGroup.Name, lvs[0].VolumeGroupName)
	}
}

func SkipTestIfNotRoot(t *testing.T) {
	if os.Geteuid() != 0 {
		t.Skip("Skipping test because it requires root privileges to setup its environment.")
	}
}
