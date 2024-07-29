package lvm2go_test

import (
	"context"
	"testing"

	. "github.com/jakobmoellerdev/lvm2go"
)

func TestVGExtend(t *testing.T) {
	t.Parallel()
	SkipOrFailTestIfNotRoot(t)

	clnt := NewClient()
	ctx := context.Background()

	test := test{
		LoopDevices: []Size{
			MustParseSize("10M"),
		},
		Volumes: []TestLogicalVolume{{
			Options: LVCreateOptionList{
				MustParseExtents("100%FREE"),
			},
		}},
	}

	addedDevices := LoopbackDevices{
		MakeTestLoopbackDevice(t, MustParseSize("10M")),
		MakeTestLoopbackDevice(t, MustParseSize("10M")),
	}

	infra := test.SetupDevicesAndVolumeGroup(t)

	if err := clnt.VGExtend(ctx, infra.volumeGroup.Name, addedDevices.PhysicalVolumeNames()); err != nil {
		t.Fatal(err)
	}

	vg, err := clnt.VG(ctx, infra.volumeGroup.Name)
	if err != nil {
		t.Fatal(err)
	}

	if vg.PvCount != 3 {
		t.Fatalf("expected 3 physical volumes, got %d", vg.PvCount)
	}

	if err := clnt.VGReduce(ctx, infra.volumeGroup.Name, addedDevices.PhysicalVolumeNames()); err != nil {
		t.Fatal(err)
	}

	vg, err = clnt.VG(ctx, infra.volumeGroup.Name)
	if err != nil {
		t.Fatal(err)
	}

	if vg.PvCount != 1 {
		t.Fatalf("expected 3 physical volumes, got %d", vg.PvCount)
	}

}
