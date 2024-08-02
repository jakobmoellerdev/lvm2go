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

type Devices []string

func (opt Devices) ApplyToLVsOptions(opts *LVsOptions) {
	opts.Devices = opt
}
func (opt Devices) ApplyToLVCreateOptions(opts *LVCreateOptions) {
	opts.Devices = opt
}
func (opt Devices) ApplyToLVRemoveOptions(opts *LVRemoveOptions) {
	opts.Devices = opt
}
func (opt Devices) ApplyToVGCreateOptions(opts *VGCreateOptions) {
	opts.Devices = opt
}
func (opt Devices) ApplyToVGRemoveOptions(opts *VGRemoveOptions) {
	opts.Devices = opt
}
func (opt Devices) ApplyToVGExtendOptions(opts *VGExtendOptions) {
	opts.Devices = opt
}
func (opt Devices) ApplyToVGReduceOptions(opts *VGReduceOptions) {
	opts.Devices = opt
}
func (opt Devices) ApplyToVGRenameOptions(opts *VGRenameOptions) {
	opts.Devices = opt
}
func (opt Devices) ApplyToVGChangeOptions(opts *VGChangeOptions) {
	opts.Devices = opt
}
func (opt Devices) ApplyToPVCreateOptions(opts *PVCreateOptions) {
	opts.Devices = opt
}
func (opt Devices) ApplyToPVRemoveOptions(opts *PVRemoveOptions) {
	opts.Devices = opt
}
func (opt Devices) ApplyToPVResizeOptions(opts *PVResizeOptions) {
	opts.Devices = opt
}
func (opt Devices) ApplyToPVChangeOptions(opts *PVChangeOptions) {
	opts.Devices = opt
}

func (opt Devices) ApplyToArgs(args Arguments) error {
	if len(opt) == 0 {
		return nil
	}
	args.AddOrReplaceAll([]string{"--devices", strings.Join(opt, ",")})
	return nil
}

type DevicesFile string

func (opt DevicesFile) ApplyToLVsOptions(opts *LVsOptions) {
	opts.DevicesFile = opt
}
func (opt DevicesFile) ApplyToLVCreateOptions(opts *LVCreateOptions) {
	opts.DevicesFile = opt
}
func (opt DevicesFile) ApplyToLVRemoveOptions(opts *LVRemoveOptions) {
	opts.DevicesFile = opt
}
func (opt DevicesFile) ApplyToVGsOptions(opts *VGsOptions) {
	opts.DevicesFile = opt
}
func (opt DevicesFile) ApplyToVGCreateOptions(opts *VGCreateOptions) {
	opts.DevicesFile = opt
}
func (opt DevicesFile) ApplyToVGRemoveOptions(opts *VGRemoveOptions) {
	opts.DevicesFile = opt
}
func (opt DevicesFile) ApplyToVGExtendOptions(opts *VGExtendOptions) {
	opts.DevicesFile = opt
}
func (opt DevicesFile) ApplyToVGReduceOptions(opts *VGReduceOptions) {
	opts.DevicesFile = opt
}
func (opt DevicesFile) ApplyToVGRenameOptions(opts *VGRenameOptions) {
	opts.DevicesFile = opt
}
func (opt DevicesFile) ApplyToVGChangeOptions(opts *VGChangeOptions) {
	opts.DevicesFile = opt
}
func (opt DevicesFile) ApplyToDevListOptions(opts *DevListOptions) {
	opts.DevicesFile = opt
}
func (opt DevicesFile) ApplyToDevCheckOptions(opts *DevCheckOptions) {
	opts.DevicesFile = opt
}
func (opt DevicesFile) ApplyToDevUpdateOptions(opts *DevUpdateOptions) {
	opts.DevicesFile = opt
}
func (opt DevicesFile) ApplyToDevModifyOptions(opts *DevModifyOptions) {
	opts.DevicesFile = opt
}
func (opt DevicesFile) ApplyToPVCreateOptions(opts *PVCreateOptions) {
	opts.DevicesFile = opt
}
func (opt DevicesFile) ApplyToPVRemoveOptions(opts *PVRemoveOptions) {
	opts.DevicesFile = opt
}
func (opt DevicesFile) ApplyToPVResizeOptions(opts *PVResizeOptions) {
	opts.DevicesFile = opt
}
func (opt DevicesFile) ApplyToPVChangeOptions(opts *PVChangeOptions) {
	opts.DevicesFile = opt
}

func (opt DevicesFile) ApplyToArgs(args Arguments) error {
	if opt == "" {
		return nil
	}
	args.AddOrReplaceAll([]string{"--devicesfile", string(opt)})
	return nil
}
