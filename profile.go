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
	"path/filepath"
	"strings"
)

var ErrInvalidProfileExtension = fmt.Errorf("profile extension must be empty or %q", LVMProfileExtension)

type Profile string

func (opt Profile) ApplyToArgs(args Arguments) error {
	if opt == "" {
		return nil
	}
	ext := filepath.Ext(string(opt))
	if ext != "" && ext != LVMProfileExtension {
		return ErrInvalidProfileExtension
	}

	path := filepath.Base(string(opt))
	split := strings.Split(path, ".")

	args.AddOrReplace("--profile", split[0])
	return nil
}

func (opt Profile) ApplyToVGsOptions(opts *VGsOptions) {
	opts.Profile = opt
}
func (opt Profile) ApplyToVGCreateOptions(opts *VGCreateOptions) {
	opts.Profile = opt
}
func (opt Profile) ApplyToVGRemoveOptions(opts *VGRemoveOptions) {
	opts.Profile = opt
}
func (opt Profile) ApplyToVGExtendOptions(opts *VGExtendOptions) {
	opts.Profile = opt
}
func (opt Profile) ApplyToVGReduceOptions(opts *VGReduceOptions) {
	opts.Profile = opt
}
func (opt Profile) ApplyToVGRenameOptions(opts *VGRenameOptions) {
	opts.Profile = opt
}
func (opt Profile) ApplyToVGChangeOptions(opts *VGChangeOptions) {
	opts.Profile = opt
}

func (opt Profile) ApplyToLVsOptions(opts *LVsOptions) {
	opts.Profile = opt
}
func (opt Profile) ApplyToLVCreateOptions(opts *LVCreateOptions) {
	opts.Profile = opt
}
func (opt Profile) ApplyToLVRemoveOptions(opts *LVRemoveOptions) {
	opts.Profile = opt
}
func (opt Profile) ApplyToLVResizeOptions(opts *LVResizeOptions) {
	opts.Profile = opt
}
func (opt Profile) ApplyToLVExtendOptions(opts *LVExtendOptions) {
	opts.Profile = opt
}
func (opt Profile) ApplyToLVReduceOptions(opts *LVReduceOptions) {
	opts.Profile = opt
}
func (opt Profile) ApplyToLVRenameOptions(opts *LVRenameOptions) {
	opts.Profile = opt
}
func (opt Profile) ApplyToLVChangeOptions(opts *LVChangeOptions) {
	opts.Profile = opt
}

func (opt Profile) ApplyToPVsOptions(opts *PVsOptions) {
	opts.Profile = opt
}
func (opt Profile) ApplyToPVCreateOptions(opts *PVCreateOptions) {
	opts.Profile = opt
}
func (opt Profile) ApplyToPVRemoveOptions(opts *PVRemoveOptions) {
	opts.Profile = opt
}
func (opt Profile) ApplyToPVResizeOptions(opts *PVResizeOptions) {
	opts.Profile = opt
}
func (opt Profile) ApplyToPVChangeOptions(opts *PVChangeOptions) {
	opts.Profile = opt
}
func (opt Profile) ApplyToPVMoveOptions(opts *PVMoveOptions) {
	opts.Profile = opt
}
func (opt Profile) ApplyToConfigOptions(opts *ConfigOptions) {
	opts.Profile = opt
}
