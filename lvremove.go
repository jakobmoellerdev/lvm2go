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
)

type (
	LVRemoveOptions struct {
		LogicalVolumeName
		VolumeGroupName

		Force
		Tags
		Select

		CommonOptions
	}
	LVRemoveOption interface {
		ApplyToLVRemoveOptions(opts *LVRemoveOptions)
	}
	LVRemoveOptionsList []LVRemoveOption
)

var (
	_ ArgumentGenerator = LVRemoveOptionsList{}
	_ Argument          = (*LVRemoveOptions)(nil)
)

func (c *client) LVRemove(ctx context.Context, opts ...LVRemoveOption) error {
	args, err := LVRemoveOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunLVM(ctx, append([]string{"lvremove"}, args.GetRaw()...)...)
}

func (opts *LVRemoveOptions) ApplyToArgs(args Arguments) error {
	id, err := NewFQLogicalVolumeName(opts.VolumeGroupName, opts.LogicalVolumeName)
	if err != nil {
		return err
	}

	for _, arg := range []Argument{
		id,
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

func (list LVRemoveOptionsList) AsArgs() (Arguments, error) {
	args := NewArgs(ArgsTypeGeneric)
	options := LVRemoveOptions{}
	for _, opt := range list {
		opt.ApplyToLVRemoveOptions(&options)
	}
	if err := options.ApplyToArgs(args); err != nil {
		return nil, err
	}
	return args, nil
}

func (opts *LVRemoveOptions) ApplyToLVRemoveOptions(new *LVRemoveOptions) {
	*new = *opts
}
