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

package lvm2go

import (
	"fmt"
)

type Thin bool

func (opt Thin) ApplyToArgs(args Arguments) error {
	if opt {
		args.AddOrReplaceAll([]string{"--thin"})
	}
	return nil
}

func (opt Thin) ApplyToLVCreateOptions(opts *LVCreateOptions) {
	opts.Thin = opt
}

type ThinPool FQLogicalVolumeName

func (opt *ThinPool) ApplyToLVCreateOptions(opts *LVCreateOptions) {
	opts.ThinPool = opt
}

func (opt *ThinPool) ApplyToLVRemoveOptions(opts *LVRemoveOptions) {
	opts.LogicalVolumeName, opts.VolumeGroupName = opt.LogicalVolumeName, opt.VolumeGroupName
}

func (opt *ThinPool) ApplyToArgs(args Arguments) error {
	if opt == nil {
		return nil
	}

	if err := (*FQLogicalVolumeName)(opt).Validate(); err != nil {
		return err
	}

	args.AddOrReplace(fmt.Sprintf(
		"--thinpool=%s",
		fmt.Sprintf("%s/%s", opt.VolumeGroupName, opt.LogicalVolumeName),
	))

	return nil
}

func MustNewThinPool(vg VolumeGroupName, lv LogicalVolumeName) *ThinPool {
	fq, err := NewThinPool(vg, lv)
	if err != nil {
		panic(err)
	}
	return fq
}

func NewThinPool(vg VolumeGroupName, lv LogicalVolumeName) (*ThinPool, error) {
	fq, err := NewFQLogicalVolumeName(vg, lv)
	if err != nil {
		return nil, err
	}
	return (*ThinPool)(fq), fq.Validate()
}
