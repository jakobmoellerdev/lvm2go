/*
 Copyright 2024 The lvm2go Authors.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package lvm2go_test

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	. "github.com/jakobmoellerdev/lvm2go"
)

func TestLVs(t *testing.T) {
	t.Parallel()
	slog.SetDefault(slog.New(NewContextPropagatingSlogHandler(NewTestingHandler(t))))
	slog.SetLogLoggerLevel(slog.LevelDebug)

	SkipOrFailTestIfNotRoot(t)
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
			t.Parallel()
			ctx := WithCustomEnvironment(context.Background(), map[string]string{})
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			SkipOrFailTestIfNotRoot(t)
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

			var pvs []*PhysicalVolume
			success := false
			for attempt := 0; attempt < 3; attempt++ {
				if pvs, err = clnt.PVs(ctx, infra.volumeGroup.Name); err != nil {
					t.Logf("failed to get physical volumes: %s", err)
				}
				if len(pvs) != len(infra.loopDevices) {
					t.Logf("expected %d physical volumes, got %d, pvs may not be updated yet", len(infra.loopDevices), len(pvs))
				}
				success = true
			}
			if !success {
				t.Fatalf("failed to get physical volumes: %s", err)
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
