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

func TestLVAttributes(t *testing.T) {
	t.Parallel()
	type args struct {
		raw string
	}
	tests := []struct {
		name  string
		args  args
		want  LVAttributes
		error error
	}{
		{
			"RAID LVMConfig without Initial Sync",
			args{raw: "Rwi-a-r---"},
			LVAttributes{
				VolumeType:             VolumeTypeRAIDNoInitialSync,
				LVPermissions:          LVPermissionsWriteable,
				LVAllocationPolicyAttr: LVAllocationPolicyAttrInherited,
				Minor:                  MinorFalse,
				State:                  StateActive,
				Open:                   OpenFalse,
				OpenTarget:             OpenTargetRaid,
				ZeroAttr:               ZeroAttrFalse,
				VolumeHealth:           VolumeHealthOK,
				SkipActivation:         SkipActivationFalse,
			},
			nil,
		},
		{
			"ThinPool with Zeroing",
			args{raw: "twi-a-tz--"},
			LVAttributes{
				VolumeType:             VolumeTypeThinPool,
				LVPermissions:          LVPermissionsWriteable,
				LVAllocationPolicyAttr: LVAllocationPolicyAttrInherited,
				Minor:                  MinorFalse,
				State:                  StateActive,
				Open:                   OpenFalse,
				OpenTarget:             OpenTargetThin,
				ZeroAttr:               ZeroAttrTrue,
				VolumeHealth:           VolumeHealthOK,
				SkipActivation:         SkipActivationFalse,
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseLVAttributes(tt.args.raw)
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

func TestVerifyHealth(t *testing.T) {
	tests := []struct {
		name    string
		rawAttr string
		wantErr error
	}{
		{
			name:    "Partial Activation",
			rawAttr: "--------p-",
			wantErr: ErrPartialActivation,
		},
		{
			name:    "Unknown Volume Health",
			rawAttr: "--------X-",
			wantErr: ErrUnknownVolumeHealth,
		},
		{
			name:    "Write Cache Error",
			rawAttr: "--------E-",
			wantErr: ErrWriteCacheError,
		},
		{
			name:    "Thin Pool Failed",
			rawAttr: "t-------F-",
			wantErr: ErrThinPoolFailed,
		},
		{
			name:    "Thin Pool Out of Data Space",
			rawAttr: "t-------D-",
			wantErr: ErrThinPoolOutOfDataSpace,
		},
		{
			name:    "Thin Volume Failed",
			rawAttr: "V-------F-",
			wantErr: ErrThinVolumeFailed,
		},
		{
			name:    "RAID Refresh Needed",
			rawAttr: "r-------r-",
			wantErr: ErrRAIDRefreshNeeded,
		},
		{
			name:    "RAID Mismatches Exist",
			rawAttr: "r-------m-",
			wantErr: ErrRAIDMismatchesExist,
		},
		{
			name:    "RAID Reshaping",
			rawAttr: "r-------s-",
			wantErr: ErrRAIDReshaping,
		},
		{
			name:    "RAID Reshape Removed",
			rawAttr: "r-------R-",
			wantErr: ErrRAIDReshapeRemoved,
		},
		{
			name:    "RAID Write Mostly",
			rawAttr: "r-------w-",
			wantErr: ErrRAIDWriteMostly,
		},
		{
			name:    "Logical Volume Suspended",
			rawAttr: "----s-----",
			wantErr: ErrLogicalVolumeSuspended,
		},
		{
			name:    "Invalid Snapshot",
			rawAttr: "----I-----",
			wantErr: ErrInvalidSnapshot,
		},
		{
			name:    "Snapshot Merge Failed",
			rawAttr: "----m-----",
			wantErr: ErrSnapshotMergeFailed,
		},
		{
			name:    "Mapped Device Present With Inactive Tables",
			rawAttr: "----i-----",
			wantErr: ErrMappedDevicePresentWithInactiveTables,
		},
		{
			name:    "Mapped Device Present Without Tables",
			rawAttr: "----d-----",
			wantErr: ErrMappedDevicePresentWithoutTables,
		},
		{
			name:    "Thin Pool Check Needed",
			rawAttr: "----c-----",
			wantErr: ErrThinPoolCheckNeeded,
		},
		{
			name:    "Unknown Volume State",
			rawAttr: "----X-----",
			wantErr: ErrUnknownVolumeState,
		},
		{
			name:    "Historical Volume State",
			rawAttr: "----h-----",
			wantErr: ErrHistoricalVolumeState,
		},
		{
			name:    "Logical Volume Underlying Device State Unknown",
			rawAttr: "-----X----",
			wantErr: ErrLogicalVolumeUnderlyingDeviceStateUnknown,
		},
		{
			name:    "Healthy Volume",
			rawAttr: "-wi-a-----",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr, err := ParseLVAttributes(tt.rawAttr)
			if err != nil {
				t.Fatalf("ParsedLvAttr() error = %v", err)
			}
			if err := attr.VerifyHealth(); !errors.Is(err, tt.wantErr) {
				t.Errorf("VerifyHealth() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
