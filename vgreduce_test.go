package lvm2go_test

import (
	"context"
	"fmt"
	"testing"

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
		t.Fatal(err)
	}

	if err := clnt.VGChange(ctx, infra.volumeGroup.Name, MaximumPhysicalVolumes(5)); err == nil {
		t.Fatal("expected error due to device missing")
	} else if !IsLVMErrVGImmutableDueToMissingPVs(err) {
		t.Fatal(fmt.Errorf("unexpected error: %v", err))
	}

	if err = clnt.VGReduce(ctx, infra.volumeGroup.Name, RemoveMissing(true)); err == nil {
		t.Fatal("expected error due to device missing")
	}
	if !IsLVMCouldNotFindDeviceWithUUID(err) {
		t.Fatal(fmt.Errorf("unexpected error: %v", err))
	}
	if !IsLVMPartialLVNeedsRepairOrRemove(err) {
		t.Fatal(fmt.Errorf("unexpected error: %v", err))
	}
	if !IsLVMErrThereAreStillPartialLVs(err) {
		t.Fatal(fmt.Errorf("unexpected error: %v", err))
	}
	if !IsLVMErrVGMissingPVs(err) {
		t.Fatal(fmt.Errorf("unexpected error: %v", err))
	}

	if vg, pv, lastWritten, ok := LVMErrVGMissingPVsDetails(err); !ok {
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
		t.Fatal(err)
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
	); err != nil && !IsLVMErrNoDataToMove(err) {
		t.Fatal(err)
	}

	if err := infra.loopDevices[1].Close(); err != nil {
		t.Fatal(err)
	}

	if err := clnt.LVChange(ctx, infra.volumeGroup.Name, infra.lvs[0].LogicalVolumeName(), Deactivate); err != nil {
		t.Fatal(err)
	}

	if err := clnt.VGChange(ctx, infra.volumeGroup.Name, MaximumPhysicalVolumes(5)); err == nil {
		t.Fatal("expected error due to device missing")
	} else if !IsLVMErrVGImmutableDueToMissingPVs(err) {
		t.Fatal(fmt.Errorf("unexpected error: %v", err))
	}

	if err = clnt.VGReduce(ctx, infra.volumeGroup.Name, RemoveMissing(true)); err != nil {
		t.Fatal(err)
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
