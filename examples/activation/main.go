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

package main

import (
	"context"
	"log/slog"
	"os"

	. "github.com/jakobmoellerdev/lvm2go"
)

func main() {
	if os.Geteuid() != 0 {
		panic("panicking because it requires root privileges to setup its environment.")
	}

	ctx := context.Background()
	lvm := NewClient()

	deviceSize := MustParseSize("200M")

	losetup, err := NewLoopbackDevice(deviceSize)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer func() {
		if err := losetup.Close(); err != nil {
			slog.Error(err.Error())
		}
	}()

	pvs := PhysicalVolumesFrom(losetup.Device())
	vgName := VolumeGroupName("test")
	lvName := LogicalVolumeName("test")
	lvSize := MustParseSize("100M")

	if err := lvm.VGCreate(ctx, vgName, pvs); err != nil {
		slog.Error(err.Error())
		return
	}
	defer func() {
		if err := lvm.VGRemove(ctx, vgName); err != nil {
			slog.Error(err.Error())
		}
	}()

	if err := lvm.LVCreate(ctx, vgName, lvName, lvSize, Deactivate); err != nil {
		slog.Error(err.Error())
		return
	}
	defer func() {
		if err := lvm.LVRemove(ctx, vgName, lvName); err != nil {
			slog.Error(err.Error())
		}
	}()

	state, err := getActivationState(ctx, lvm, vgName, lvName)
	if err != nil {
		slog.Error(err.Error())
		return
	} else if state == StateUnknown {
		slog.Error("state was unknown")
		return
	} else if state == StateActive {
		slog.Error("expected logical volume to be inactive", slog.String("state", string(state)))
		return
	}

	if err := lvm.LVChange(ctx, vgName, lvName, Activate); err != nil {
		slog.Error(err.Error())
		return
	}

	state, err = getActivationState(ctx, lvm, vgName, lvName)
	if err != nil {
		slog.Error(err.Error())
		return
	} else if state == StateUnknown {
		slog.Error("state was unknown")
		return
	} else if state != StateActive {
		slog.Error("expected logical volume to be active", slog.String("state", string(state)))
		return
	}
}

func getActivationState(ctx context.Context, lvm Client, vgName VolumeGroupName, lvName LogicalVolumeName) (State, error) {
	lvs, err := lvm.LVs(ctx, vgName, NewMatchesAllSelector(map[string]string{"lv_name": string(lvName)}))
	if err != nil {
		return StateUnknown, err
	}
	for _, lv := range lvs {
		if lv.Name == lvName {
			return lv.Attr.State, nil
		}
	}
	return StateUnknown, nil
}
