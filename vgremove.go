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
	"fmt"
)

type (
	VGRemoveOptions struct {
		VolumeGroupName
		Tags
		Select

		Force

		CommonOptions
	}
	VGRemoveOption interface {
		ApplyToVGRemoveOptions(opts *VGRemoveOptions)
	}
	VGRemoveOptionsList []VGRemoveOption
)

var (
	_ ArgumentGenerator = VGRemoveOptionsList{}
	_ Argument          = (*VGRemoveOptions)(nil)
)

func (c *client) VGRemove(ctx context.Context, opts ...VGRemoveOption) error {
	args, err := VGRemoveOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunLVM(ctx, append([]string{"vgremove"}, args.GetRaw()...)...)
}

func (opts *VGRemoveOptions) ApplyToArgs(args Arguments) error {
	if opts.VolumeGroupName == "" {
		return fmt.Errorf("VolumeGroupName is required for removal of a volume group")
	}

	for _, arg := range []Argument{
		opts.VolumeGroupName,
		opts.Tags,
		opts.Force,
		opts.CommonOptions,
	} {
		if err := arg.ApplyToArgs(args); err != nil {
			return err
		}
	}

	return nil
}

func (opts *VGRemoveOptions) ApplyToVGRemoveOptions(new *VGRemoveOptions) {
	*new = *opts
}

var (
	_ ArgumentGenerator = VGRemoveOptionsList{}
	_ Argument          = (*VGRemoveOptions)(nil)
)

func (list VGRemoveOptionsList) AsArgs() (Arguments, error) {
	args := NewArgs(ArgsTypeGeneric)
	options := VGRemoveOptions{}
	for _, opt := range list {
		opt.ApplyToVGRemoveOptions(&options)
	}
	if err := options.ApplyToArgs(args); err != nil {
		return nil, err
	}
	return args, nil
}
