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
	LVResizeOptions struct {
		LogicalVolumeName
		VolumeGroupName

		PrefixedSize

		CommonOptions
	}
	LVResizeOption interface {
		ApplyToLVResizeOptions(opts *LVResizeOptions)
	}
	LVResizeOptionsList []LVResizeOption
)

var (
	_ ArgumentGenerator = LVResizeOptionsList{}
	_ Argument          = (*LVResizeOptions)(nil)
)

func (c *client) LVResize(ctx context.Context, opts ...LVResizeOption) error {
	args, err := LVResizeOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunLVM(ctx, append([]string{"lvresize"}, args.GetRaw()...)...)
}

func (list LVResizeOptionsList) AsArgs() (Arguments, error) {
	args := NewArgs(ArgsTypeGeneric)
	options := LVResizeOptions{}
	for _, opt := range list {
		opt.ApplyToLVResizeOptions(&options)
	}
	if err := options.ApplyToArgs(args); err != nil {
		return nil, err
	}
	return args, nil
}

func (opts *LVResizeOptions) ApplyToArgs(args Arguments) error {
	id, err := NewFQLogicalVolumeName(opts.VolumeGroupName, opts.LogicalVolumeName)
	if err != nil {
		return err
	}

	for _, opt := range []Argument{
		id,
		opts.PrefixedSize,
		opts.CommonOptions,
	} {
		if err := opt.ApplyToArgs(args); err != nil {
			return err
		}
	}

	return nil
}
