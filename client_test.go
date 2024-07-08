package lvm2go

import (
	"context"
	"log/slog"
	"testing"
)

func TestLVs(t *testing.T) {
	FailTestIfNotRoot(t)
	slog.SetDefault(slog.New(NewContextPropagatingSlogHandler(NewTestingHandler(t))))
	slog.SetLogLoggerLevel(slog.LevelDebug)
	for _, tc := range []test{
		{
			loopDevices: []Size{
				MustParseSize("1G"),
			},
			lvs: []TestLogicalVolume{
				{Options: LVCreateOptionList{
					MustParseSize("100M"),
				}},
			},
		},
		{
			loopDevices: []Size{
				MustParseSize("100M"),
			},
			lvs: []TestLogicalVolume{
				{Options: LVCreateOptionList{
					MustParseSize("100M"),
				}},
			},
		},
		{
			loopDevices: []Size{
				MustParseSize("1G"),
			},
			lvs: []TestLogicalVolume{
				{Options: LVCreateOptionList{
					MustParseSize("1G"),
				}},
			},
		},
		{
			loopDevices: []Size{
				MustParseSize("4G"),
			},
			lvs: []TestLogicalVolume{
				{Options: LVCreateOptionList{
					MustParseSize("2G"),
				}},
				{Options: LVCreateOptionList{
					MustParseSize("2G"),
				}},
			},
		},
	} {
		t.Run(tc.String(), func(t *testing.T) {
			FailTestIfNotRoot(t)
			ctx := context.Background()
			clnt := GetTestClient(ctx)
			infra := tc.SetupDevicesAndVolumeGroup(t)

			lvs, err := clnt.LVs(ctx, infra.volumeGroup.Name)
			if err != nil {
				t.Fatal(err)
			}

			if len(lvs) != len(tc.lvs) {
				t.Fatalf("Expected %d logical volumes, got %d", len(tc.lvs), len(lvs))
			}

			for _, expected := range infra.lvs {
				found := false
				for _, lv := range lvs {
					if lv.Name != expected.LogicalVolumeName() {
						continue
					}
					if eq, err := lv.Size.IsEqualTo(expected.Size()); err != nil || !eq {
						if err != nil {
							t.Fatalf("Size inconsistency: %s", err)
						}
					}
					found = true
					break
				}
				if !found {
					t.Fatalf("Expected logical volume %s not found in LVs report", expected)
				}
			}
		})
	}
}
