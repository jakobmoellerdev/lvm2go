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
	"errors"
	"testing"

	. "github.com/jakobmoellerdev/lvm2go"
)

func TestPVAttributes(t *testing.T) {
	t.Parallel()
	type args struct {
		raw string
	}
	tests := []struct {
		name  string
		args  args
		want  PVAttributes
		error error
	}{
		{
			"pv used in vg",
			args{raw: "u--"},
			PVAttributes{
				DuplicateAllocatableUsed: Used,
				Exported:                 ExportedFalse,
				Missing:                  MissingFalse,
			},
			nil,
		},
		{
			"pv unused",
			args{raw: "a--"},
			PVAttributes{
				DuplicateAllocatableUsed: Allocatable,
				Exported:                 ExportedFalse,
				Missing:                  MissingFalse,
			},
			nil,
		},
		{
			"pv duplicate",
			args{raw: "d--"},
			PVAttributes{
				DuplicateAllocatableUsed: Duplicate,
				Exported:                 ExportedFalse,
				Missing:                  MissingFalse,
			},
			nil,
		},
		{
			"pv removed or broken or missing",
			args{raw: "u-m"},
			PVAttributes{
				DuplicateAllocatableUsed: Used,
				Exported:                 ExportedFalse,
				Missing:                  MissingTrue,
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePVAttributes(tt.args.raw)
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
