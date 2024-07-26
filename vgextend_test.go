package lvm2go_test

import (
	"context"
	"testing"

	. "github.com/jakobmoellerdev/lvm2go"
)

func TestVGExtend(t *testing.T) {
	FailTestIfNotRoot(t)

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

	vgs, err := clnt.VGs(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(vgs) != 1 {
		t.Fatalf("expected 1 volume group, got %d", len(vgs))
	}

	if vgs[0].Name != infra.volumeGroup.Name {
		t.Fatalf("expected volume group %s, got %s", infra.volumeGroup.Name, vgs[0].Name)
	}

	vg := vgs[0]

	if vg.PvCount != 3 {
		t.Fatalf("expected 3 physical volumes, got %d", vg.PvCount)
	}

}
