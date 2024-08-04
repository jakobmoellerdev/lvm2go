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
	"context"
	"errors"
	"fmt"
)

var ErrNoDevicesSpecifiedForModification = errors.New("no devices specified for modification")

type ModifyDeviceType string

const (
	DelDev       ModifyDeviceType = "deldev"
	DelDevByPVID ModifyDeviceType = "delpvid"
	AddDev       ModifyDeviceType = "adddev"
	AddDevByPVID ModifyDeviceType = "addpvid"
)

func AddDevice(device string) ModifyDevice {
	return ModifyDevice{
		Device:           device,
		ModifyDeviceType: AddDev,
	}
}

func AddDeviceByPVID(pvid string) ModifyDevice {
	return ModifyDevice{
		Device:           pvid,
		ModifyDeviceType: AddDevByPVID,
	}
}

func DelDevice(device string) ModifyDevice {
	return ModifyDevice{
		Device:           device,
		ModifyDeviceType: DelDev,
	}
}

func DelDeviceByPVID(pvid string) ModifyDevice {
	return ModifyDevice{
		Device:           pvid,
		ModifyDeviceType: DelDevByPVID,
	}
}

type ModifyDevice struct {
	Device string
	ModifyDeviceType
}

func (opt ModifyDevice) ApplyToDevModifyOptions(opts *DevModifyOptions) {
	opts.ModifyDevice = opt
}

func (opt ModifyDevice) ApplyToArgs(args Arguments) error {
	if len(opt.Device) == 0 {
		return nil
	}
	args.AddOrReplaceAll([]string{fmt.Sprintf("--%s", string(opt.ModifyDeviceType)), opt.Device})
	return nil
}

type (
	DevModifyOptions struct {
		DevicesFile

		ModifyDevice

		DeviceIDType
	}
	DevModifyOption interface {
		ApplyToDevModifyOptions(opts *DevModifyOptions)
	}
	DevModifyOptionsList []DevModifyOption
)

var (
	_ ArgumentGenerator = DevModifyOptionsList{}
	_ Argument          = (*DevModifyOptions)(nil)
)

func (c *client) DevModify(ctx context.Context, opts ...DevModifyOption) error {
	args, err := DevModifyOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunRaw(
		ctx,
		NoOpRawOutputProcessor(false),
		append([]string{"lvmdevices"}, args.GetRaw()...)...,
	)
}

func (list DevModifyOptionsList) AsArgs() (Arguments, error) {
	args := NewArgs(ArgsTypeGeneric)
	options := DevModifyOptions{}
	for _, opt := range list {
		opt.ApplyToDevModifyOptions(&options)
	}
	if err := options.ApplyToArgs(args); err != nil {
		return nil, err
	}
	return args, nil
}

func (opts *DevModifyOptions) ApplyToDevModifyOptions(new *DevModifyOptions) {
	*new = *opts
}

func (opts *DevModifyOptions) ApplyToArgs(args Arguments) error {
	if err := opts.DevicesFile.ApplyToArgs(args); err != nil {
		return err
	}

	if len(opts.ModifyDevice.Device) == 0 {
		return ErrNoDevicesSpecifiedForModification
	}

	if err := opts.ModifyDevice.ApplyToArgs(args); err != nil {
		return err
	}
	return nil
}
