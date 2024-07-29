package lvm2go_test

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	. "github.com/jakobmoellerdev/lvm2go"
)

func TestLVs(t *testing.T) {
	SkipOrFailTestIfNotRoot(t)
	ctx := WithCustomEnvironment(context.Background(), map[string]string{})
	slog.SetDefault(slog.New(NewContextPropagatingSlogHandler(NewTestingHandler(t))))
	slog.SetLogLoggerLevel(slog.LevelDebug)
	for i, tc := range []test{
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
		t.Run(fmt.Sprintf("[%v]%s", i, tc.String()), func(t *testing.T) {
			SkipOrFailTestIfNotRoot(t)
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			clnt := GetTestClient(ctx)
			infra := tc.SetupDevicesAndVolumeGroup(t)

			lvs, err := clnt.LVs(ctx, infra.volumeGroup.Name)
			if err != nil {
				t.Fatal(err)
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

			vg, err := clnt.VG(ctx, infra.volumeGroup.Name)
			if err != nil {
				t.Fatal(err)
			}

			if vg.Name != infra.volumeGroup.Name {
				t.Fatalf("Expected volume group %s, got %s", infra.volumeGroup.Name, vg.Name)
			}

			pvs, err := clnt.PVs(ctx)
			if err != nil {
				t.Fatal(err)
			}
			if len(pvs) != len(infra.loopDevices) {
				t.Fatalf("Expected %d physical volumes, got %d", len(infra.loopDevices), len(pvs))
			}

			for _, pv := range pvs {
				found := false
				for _, ld := range infra.loopDevices {
					if string(pv.Name) == ld.Device() {
						found = true
						break
					}
				}
				if !found {
					t.Fatalf("physical volume %s in PVs report is not part of the original loop devices", pv.Name)
				}
			}

		})
	}
}
