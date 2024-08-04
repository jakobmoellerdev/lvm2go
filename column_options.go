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
	"strings"
)

var (
	DefaultLVsColumnOptions = ColumnOptions{
		"lv_all",
	}
	DefaultVGsColumnOptions = ColumnOptions{
		"vg_all",
	}
	DefaultPVsColumnOptions = ColumnOptions{
		"pv_all",
	}
)

type ColumnOptions []string

func (opt ColumnOptions) ApplyToLVsOptions(opts *LVsOptions) {
	opts.ColumnOptions = opt
}

func (opt ColumnOptions) ApplyToVGsOptions(opts *VGsOptions) {
	opts.ColumnOptions = opt
}

func (opt ColumnOptions) ApplyToArgs(args Arguments) error {
	var optionsString string
	if len(opt) > 0 {
		optionsString = strings.Join(opt, ",")
	} else {
		switch args.GetType() {
		case ArgsTypeVGs:
			optionsString = strings.Join(DefaultVGsColumnOptions, ",")
		case ArgsTypeLVs:
			optionsString = strings.Join(DefaultLVsColumnOptions, ",")
		case ArgsTypePVs:
			optionsString = strings.Join(DefaultPVsColumnOptions, ",")
		}
	}
	args.AddOrReplaceAll([]string{"--options", optionsString})
	return nil
}
