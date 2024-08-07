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
	"testing"

	. "github.com/jakobmoellerdev/lvm2go"
)

func TestMaximumLogicalVolumes(t *testing.T) {
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
				MustParseSize("4M"),
			},
		}},
		AdditionalVolumeGroupOptions: []VGCreateOption{
			MaximumLogicalVolumes(1),
			MaximumPhysicalVolumes(1),
		},
	}

	infra := test.SetupDevicesAndVolumeGroup(t)

	t.Run("maximum logical volumes", func(t *testing.T) {
		if err := clnt.LVCreate(
			ctx,
			infra.volumeGroup.Name,
			LogicalVolumeName("exceedinglimit"),
			MustParseSize("4M"),
		); err == nil {
			t.Fatal("expected error")
		} else if !IsMaximumLogicalVolumesReached(err) {
			t.Fatalf("expected maximum number of logical volumes reached error, but got %s", err)
		}

		if err := clnt.VGChange(ctx, infra.volumeGroup.Name, MaximumLogicalVolumes(2)); err != nil {
			t.Fatal(err)
		}

		lvName := LogicalVolumeName("withinlimit")
		if err := clnt.LVCreate(ctx, infra.volumeGroup.Name, lvName, MustParseSize("4M")); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			if err := clnt.LVRemove(ctx, infra.volumeGroup.Name, lvName); err != nil {
				if !IsSkippableErrorForCleanup(err) {
					t.Fatal(err)
				}
			}
		})
	})

	t.Run("maximum physical volumes", func(t *testing.T) {
		additionalLoop := MakeTestLoopbackDevice(t, MustParseSize("10M"))
		if err := clnt.VGExtend(
			ctx,
			infra.volumeGroup.Name,
			PhysicalVolumeName(additionalLoop.Device()),
		); err == nil {
			t.Fatal("expected error")
		} else if !IsMaximumPhysicalVolumesReached(err) {
			t.Fatalf("expected maximum number of physical volumes reached error, but got %s", err)
		}

		if err := clnt.VGChange(ctx, infra.volumeGroup.Name, MaximumPhysicalVolumes(2)); err != nil {
			t.Fatal(err)
		}

		if err := clnt.VGExtend(
			ctx,
			infra.volumeGroup.Name,
			PhysicalVolumeName(additionalLoop.Device()),
		); err != nil {
			t.Fatal(err)
		}
	})

}
