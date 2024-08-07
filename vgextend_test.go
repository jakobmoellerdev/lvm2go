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

func TestVGExtend(t *testing.T) {
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

	vg, err := clnt.VG(ctx, infra.volumeGroup.Name)
	if err != nil {
		t.Fatal(err)
	}

	if vg.PvCount != 3 {
		t.Fatalf("expected 3 physical volumes, got %d", vg.PvCount)
	}

	if err := clnt.VGReduce(ctx, infra.volumeGroup.Name, addedDevices.PhysicalVolumeNames()); err != nil {
		if IsSkippableErrorForCleanup(err) {
			t.Logf("vgreduce on loop devices failed due to skippable error, ignoring: %s", err)
		} else {
			t.Fatal(err)
		}
	}

	vg, err = clnt.VG(ctx, infra.volumeGroup.Name)
	if err != nil {
		t.Fatal(err)
	}

	if vg.PvCount != 1 {
		t.Fatalf("expected 3 physical volumes, got %d", vg.PvCount)
	}

}
