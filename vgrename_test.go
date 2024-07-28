package lvm2go_test

import (
	"context"
	"testing"

	. "github.com/jakobmoellerdev/lvm2go"
)

func TestVGRename(t *testing.T) {
	t.Parallel()
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
