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

	. "github.com/jakobmoellerdev/lvm2go"
)

func TestLVRename(t *testing.T) {
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

	for _, lv := range infra.lvs {
		oldName := lv.LogicalVolumeName()
		newName := LogicalVolumeName(fmt.Sprintf("%s-renamed", lv.LogicalVolumeName()))
		if err := clnt.LVRename(ctx, &LVRenameOptions{
			VolumeGroupName: infra.volumeGroup.Name,
			Old:             lv.LogicalVolumeName(),
			New:             newName,
		}); err != nil {
			t.Fatal(err)
		}

		if err := clnt.LVRename(ctx, infra.volumeGroup.Name, newName, oldName); err != nil {
			t.Fatal(err)
		}
	}

}
