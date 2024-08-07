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
	"testing"
	"time"

	. "github.com/jakobmoellerdev/lvm2go"
)

func TestVGReduceByForce(t *testing.T) {
	t.Parallel()
	SkipOrFailTestIfNotRoot(t)

	clnt := NewClient()
	ctx := context.Background()

	test := test{
		LoopDevices: []Size{
			MustParseSize("10M"),
			MustParseSize("10M"),
		},
		Volumes: []TestLogicalVolume{{
			Options: LVCreateOptionList{
				MustParseExtents("100%FREE"),
			},
		}},
	}

	infra := test.SetupDevicesAndVolumeGroup(t)

	vg, err := clnt.VG(ctx, infra.volumeGroup.Name)
	if err != nil {
		t.Fatal(err)
	}

	if expect := int64(len(infra.loopDevices)); vg.PvCount != expect {
		t.Fatalf("expected %d physical volumes, got %d", expect, vg.PvCount)
	}

	if err := infra.loopDevices[1].Close(); err != nil {
		t.Fatal(err)
	}

	if err := clnt.LVChange(ctx, infra.volumeGroup.Name, infra.lvs[0].LogicalVolumeName(), Deactivate); err != nil {
		if !IsDeviceNotFound(err) {
			t.Fatal(err)
		}
	}
	// Wait to allow the device to be removed
	time.Sleep(100 * time.Millisecond)
	if err := clnt.VGChange(ctx, infra.volumeGroup.Name, MaximumPhysicalVolumes(5)); err == nil {
		t.Fatal("expected error due to device missing")
	} else if !IsVGImmutableDueToMissingPVs(err) {
		t.Fatal(fmt.Errorf("unexpected error: %v", err))
	}

	if err = clnt.VGReduce(ctx, infra.volumeGroup.Name, RemoveMissing(true)); err == nil {
		t.Fatal("expected error due to device missing")
	}
	if !IsCouldNotFindDeviceWithUUID(err) {
		t.Fatal(fmt.Errorf("unexpected error: %v", err))
	}
	if !IsPartialLVNeedsRepairOrRemove(err) {
		t.Fatal(fmt.Errorf("unexpected error: %v", err))
	}
	if !IsThereAreStillPartialLVs(err) {
		t.Fatal(fmt.Errorf("unexpected error: %v", err))
	}
	if !IsVGMissingPVs(err) {
		t.Fatal(fmt.Errorf("unexpected error: %v", err))
	}

	if vg, pv, lastWritten, ok := VGMissingPVsDetails(err); !ok {
		t.Fatal("expected details")
	} else {
		if vg != string(infra.volumeGroup.Name) {
			t.Fatalf("expected volume group %s, got %s", infra.volumeGroup.Name, vg)
		}
		if len(pv) == 0 {
			t.Fatalf("expected physical volume, got %s", pv)
		}
		if lastWritten != infra.loopDevices[1].Device() {
			t.Fatalf("expected last written to %s, got %s", infra.loopDevices[1].Device(), lastWritten)
		}
	}

	if err = clnt.VGReduce(ctx, infra.volumeGroup.Name, RemoveMissing(true), Force(true)); err != nil {
		if !IsDeviceNotFound(err) {
			t.Fatal(err)
		}
	}

	if err := clnt.VGChange(ctx, infra.volumeGroup.Name, MaximumPhysicalVolumes(5)); err != nil {
		t.Fatal(err)
	}

	if vg, err = clnt.VG(ctx, infra.volumeGroup.Name); err != nil {
		t.Fatal(err)
	}

	if expect := int64(len(infra.loopDevices) - 1); vg.PvCount != expect {
		t.Fatalf("expected %d physical volume, got %d", expect, vg.PvCount)
	}

}

func TestVGReduceByMove(t *testing.T) {
	t.Parallel()
	SkipOrFailTestIfNotRoot(t)

	clnt := NewClient()
	ctx := context.Background()

	test := test{
		LoopDevices: []Size{
			MustParseSize("10M"),
			MustParseSize("10M"),
		},
		Volumes: []TestLogicalVolume{{
			Options: LVCreateOptionList{
				MustParseExtents("30%FREE"),
			},
		}},
	}

	infra := test.SetupDevicesAndVolumeGroup(t)

	vg, err := clnt.VG(ctx, infra.volumeGroup.Name)
	if err != nil {
		t.Fatal(err)
	}

	if expect := int64(len(infra.loopDevices)); vg.PvCount != expect {
		t.Fatalf("expected %d physical volumes, got %d", expect, vg.PvCount)
	}

	if err := clnt.PVMove(ctx,
		PhysicalVolumeName(infra.loopDevices[1].Device()),
		PhysicalVolumeName(infra.loopDevices[0].Device()),
	); err != nil && !IsNoDataToMove(err) {
		t.Fatal(err)
	}

	if err := infra.loopDevices[1].Close(); err != nil {
		t.Fatal(err)
	}

	if err := clnt.LVChange(ctx, infra.volumeGroup.Name, infra.lvs[0].LogicalVolumeName(), Deactivate); err != nil {
		if !IsDeviceNotFound(err) {
			t.Fatal(err)
		}
	}

	if err := clnt.VGChange(ctx, infra.volumeGroup.Name, MaximumPhysicalVolumes(5)); err == nil {
		t.Fatal("expected error due to device missing")
	} else if !IsVGImmutableDueToMissingPVs(err) {
		t.Fatal(fmt.Errorf("unexpected error: %v", err))
	}

	if err = clnt.VGReduce(ctx, infra.volumeGroup.Name, RemoveMissing(true)); err != nil {
		if !IsDeviceNotFound(err) {
			t.Fatal(err)
		}
	}

	if err := clnt.VGChange(ctx, infra.volumeGroup.Name, MaximumPhysicalVolumes(5)); err != nil {
		t.Fatal(err)
	}

	if vg, err = clnt.VG(ctx, infra.volumeGroup.Name); err != nil {
		t.Fatal(err)
	}

	if expect := int64(len(infra.loopDevices) - 1); vg.PvCount != expect {
		t.Fatalf("expected %d physical volume, got %d", expect, vg.PvCount)
	}

}
