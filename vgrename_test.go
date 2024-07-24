package lvm2go

import (
	"context"
	"testing"
)

func TestVGRename(t *testing.T) {
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

	infra := test.SetupDevicesAndVolumeGroup(t)
	n1, n2 := infra.volumeGroup.Name, infra.volumeGroup.Name+"-new"

	if err := clnt.VGRename(ctx, &VGRenameOptions{
		Old: n1,
		New: n2,
	}); err != nil {
		t.Fatal(err)
	}

	if err := clnt.VGRename(ctx, n2, n1); err != nil {
		t.Fatal(err)
	}
}
