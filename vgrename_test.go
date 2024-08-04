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

func TestVGRename(t *testing.T) {
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

	infra := test.SetupDevicesAndVolumeGroup(t)
	n1, n2 := infra.volumeGroup.Name, infra.volumeGroup.Name+"-new"

	if err := clnt.VGRename(ctx, &VGRenameOptions{
		Old: n1,
		New: n2,
	}); err != nil {
		t.Fatal(err)
	}

	if err := clnt.VGRename(ctx, n2, n1); err != nil {
		t.Fatal(err)
	}
}
