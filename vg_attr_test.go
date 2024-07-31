package lvm2go_test

import (
	"errors"
	"testing"

	. "github.com/jakobmoellerdev/lvm2go"
)

func TestVGAttributes(t *testing.T) {
	t.Parallel()
	type args struct {
		raw string
	}
	tests := []struct {
		name  string
		args  args
		want  VGAttributes
		error error
	}{
		{
			"standard active vg",
			args{raw: "wz--n-"},
			VGAttributes{
				VGPermissions:          VGPermissionsWriteable,
				Resizeable:             ResizeableTrue,
				Exported:               ExportedFalse,
				PartialAttr:            PartialAttrFalse,
				VGAllocationPolicyAttr: VGAllocationPolicyAttrNormal,
				ClusteredOrShared:      ClusteredOrSharedFalse,
			},
			nil,
		},
		{
			"standard active vg (non-resizable)",
			args{raw: "w---n-"},
			VGAttributes{
				VGPermissions:          VGPermissionsWriteable,
				Resizeable:             ResizeableFalse,
				Exported:               ExportedFalse,
				PartialAttr:            PartialAttrFalse,
				VGAllocationPolicyAttr: VGAllocationPolicyAttrNormal,
				ClusteredOrShared:      ClusteredOrSharedFalse,
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseVGAttributes(tt.args.raw)
			if !errors.Is(err, tt.error) {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.want != got {
				t.Errorf("unexpected result: %v - expected %v", got, tt.want)
			}
			if tt.args.raw != got.String() {
				t.Errorf("unexpected string: %v - expected %v", got.String(), tt.args.raw)
			}
		})
	}
}
