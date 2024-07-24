package lvm2go

import (
	"context"
	"testing"
)

func TestLVExtend(t *testing.T) {
	FailTestIfNotRoot(t)

	clnt := NewClient()
	ctx := context.Background()

	test := test{
		LoopDevices: []Size{
			MustParseSize("100M"),
		},
		Volumes: []TestLogicalVolume{{
			Options: LVCreateOptionList{
				MustParseExtents("10%FREE"),
			},
		}},
	}

	infra := test.SetupDevicesAndVolumeGroup(t)

	for _, lv := range infra.lvs {
		if err := clnt.LVExtend(
			ctx,
			infra.volumeGroup.Name,
			lv.LogicalVolumeName(),
			MustParsePrefixedExtents("100%FREE"),
		); err != nil {
			t.Fatal(err)
		}
	}

}