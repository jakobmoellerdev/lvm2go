package lvm2go

import (
	"strings"
)

type Devices []string

func (opt Devices) ApplyToLVsOptions(opts *LVsOptions) {
	opts.Devices = opt
}
func (opt Devices) ApplyToVGsOptions(opts *VGsOptions) {
	opts.Devices = opt
}
func (opt Devices) ApplyToVGCreateOptions(opts *VGCreateOptions) {
	opts.Devices = opt
}
func (opt Devices) ApplyToVGRemoveOptions(opts *VGRemoveOptions) {
	opts.Devices = opt
}
func (opt Devices) ApplyToLVCreateOptions(opts *LVCreateOptions) {
	opts.Devices = opt
}
func (opt Devices) ApplyToLVRemoveOptions(opts *LVRemoveOptions) {
	opts.Devices = opt
}

func (opt Devices) ApplyToArgs(args Arguments) error {
	if len(opt) == 0 {
		return nil
	}
	args.AppendAll([]string{"--devices", strings.Join(opt, ",")})
	return nil
}

type DevicesFile string

func (opt DevicesFile) ApplyToLVsOptions(opts *LVsOptions) {
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
func (opt DevicesFile) ApplyToLVCreateOptions(opts *LVCreateOptions) {
	opts.DevicesFile = opt
}
func (opt DevicesFile) ApplyToLVRemoveOptions(opts *LVRemoveOptions) {
	opts.DevicesFile = opt
}

func (opt DevicesFile) ApplyToArgs(args Arguments) error {
	if opt == "" {
		return nil
	}
	args.AppendAll([]string{"--devicesfile", string(opt)})
	return nil
}
