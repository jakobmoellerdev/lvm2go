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

type MaximumLogicalVolumes int

func (opt MaximumLogicalVolumes) ApplyToVGChangeOptions(opts *VGChangeOptions) {
	opts.MaximumLogicalVolumes = opt
}

func (opt MaximumLogicalVolumes) ApplyToVGCreateOptions(opts *VGCreateOptions) {
	opts.MaximumLogicalVolumes = opt
}

func (opt MaximumLogicalVolumes) ApplyToArgs(args Arguments) error {
	if opt == 0 {
		return nil
	}
	switch args.GetType() {
	case ArgsTypeVGChange:
		args.AddOrReplace(fmt.Sprintf("--logicalvolume=%d", opt))
	default:
		args.AddOrReplace(fmt.Sprintf("--maxlogicalvolumes=%d", opt))
	}
	return nil
}
