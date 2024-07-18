package lvm2go

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var ErrEmptyDevicesFile = errors.New("devices file cannot be empty to be validated")
var ErrDevicesFileCannotBeDir = errors.New("devices file cannot be a directory")

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
	if err := opt.Validate(); err != nil {
		return err
	}
	args.AddOrReplaceAll([]string{"--devicesfile", string(opt)})
	return nil
}

func (opt DevicesFile) Validate() error {
	if opt == "" {
		return fmt.Errorf("%q is empty: %w", opt, ErrEmptyDevicesFile)
	}

	fi, err := os.Stat(string(opt))
	if err != nil {
		return fmt.Errorf("%q is not a valid devices file: %w", opt, err)
	}
	if fi.IsDir() {
		return fmt.Errorf("%q is a directory: %w", opt, ErrDevicesFileCannotBeDir)
	}

	return nil
}
