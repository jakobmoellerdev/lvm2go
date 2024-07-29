package lvm2go_test

import (
	"context"
	"slices"
	"testing"

	. "github.com/jakobmoellerdev/lvm2go"
)

func TestVGChange(t *testing.T) {
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

	infra := test.SetupDevicesAndVolumeGroup(t)

	testTags := Tags{"test"}

	getVGTags := func() Tags {
		vg, err := clnt.VG(ctx, infra.volumeGroup.Name)
		if err != nil {
			t.Fatal(err)
		}
		return vg.Tags
	}

	if err := clnt.VGChange(ctx, infra.volumeGroup.Name, testTags); err != nil {
		t.Fatal(err)
	}

	if tags := getVGTags(); len(tags) == 0 {
		t.Fatalf("expected tags, got %v", tags)
	} else {
		for _, testTag := range testTags {
			if !slices.Contains(tags, testTag) {
				t.Fatalf("expected tag %s, got %v", testTag, tags)
			}
		}
	}

	if err := clnt.VGChange(ctx, infra.volumeGroup.Name, DelTags(testTags)); err != nil {
		t.Fatal(err)
	}

	if tags := getVGTags(); len(tags) != 0 {
		t.Fatalf("expected 0 tags, got %d", len(tags))
	}

}
