package lvm2go

import (
	"context"
	"log/slog"
	"testing"
)

func TestLVs(t *testing.T) {
	FailTestIfNotRoot(t)
	ctx := WithCustomEnvironment(context.Background(), map[string]string{})
	slog.SetDefault(slog.New(NewContextPropagatingSlogHandler(NewTestingHandler(t))))
	slog.SetLogLoggerLevel(slog.LevelDebug)
	for _, tc := range []test{
		{
			LoopDevices: []Size{
				MustParseSize("10M"),
			},
			Volumes: []TestLogicalVolume{
				{Options: LVCreateOptionList{
					MustParseSize("5M"),
				}},
			},
		},
		{
			LoopDevices: []Size{
				MustParseSize("5M"),
			},
			Volumes: []TestLogicalVolume{
				{Options: LVCreateOptionList{
					MustParseSize("5M"),
				}},
			},
		},
		{
			LoopDevices: []Size{
				MustParseSize("10M"),
			},
			Volumes: []TestLogicalVolume{
				{Options: LVCreateOptionList{
					MustParseSize("5M"),
				}},
				{Options: LVCreateOptionList{
					MustParseSize("5M"),
				}},
			},
		},
		{
			LoopDevices: []Size{
				{Val: TestExtentSize.Val * 2, Unit: TestExtentSize.Unit},
			},
			Volumes: []TestLogicalVolume{
				{Options: LVCreateOptionList{
					MustParseExtents("1"),
				}},
				{Options: LVCreateOptionList{
					MustParseExtents("1"),
				}},
			},
		},
	} {
		t.Run(tc.String(), func(t *testing.T) {
			FailTestIfNotRoot(t)
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			clnt := GetTestClient(ctx)
			infra := tc.SetupDevicesAndVolumeGroup(t)

			lvs, err := clnt.LVs(ctx, infra.volumeGroup.Name)
			if err != nil {
				t.Fatal(err)
			}

			if len(lvs) != len(tc.Volumes) {
				t.Fatalf("Expected %d logical volumes, got %d", len(tc.Volumes), len(lvs))
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
