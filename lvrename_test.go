package lvm2go_test

import (
	"context"
	"fmt"
	"testing"

	. "github.com/jakobmoellerdev/lvm2go"
)

func TestLVRename(t *testing.T) {
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

	for _, lv := range infra.lvs {
		oldName := lv.LogicalVolumeName()
		newName := LogicalVolumeName(fmt.Sprintf("%s-renamed", lv.LogicalVolumeName()))
		if err := clnt.LVRename(ctx, &LVRenameOptions{
			VolumeGroupName: infra.volumeGroup.Name,
			Old:             lv.LogicalVolumeName(),
			New:             newName,
		}); err != nil {
			t.Fatal(err)
		}

		if err := clnt.LVRename(ctx, infra.volumeGroup.Name, newName, oldName); err != nil {
			t.Fatal(err)
		}
	}

}
