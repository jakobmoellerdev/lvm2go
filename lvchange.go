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
	LVChangeOptions struct {
		VolumeGroupName
		LogicalVolumeName

		Permission

		Tags
		DelTags

		Zero
		RequestConfirm
		ActivationState
		ActivationMode
		AllocationPolicy
		*ErrorWhenFull
		Partial
		SyncAction
		Rebuild
		Resync
		Discards
		*Deduplication
		*Compression
		AutoActivation

		CommonOptions
	}
	LVChangeOption interface {
		ApplyToLVChangeOptions(opts *LVChangeOptions)
	}
	LVChangeOptionsList []LVChangeOption
)

var (
	_ ArgumentGenerator = LVChangeOptionsList{}
	_ Argument          = (*LVChangeOptions)(nil)
)

func (c *client) LVChange(ctx context.Context, opts ...LVChangeOption) error {
	args, err := LVChangeOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunLVM(ctx, append([]string{"lvchange"}, args.GetRaw()...)...)
}

func (opts *LVChangeOptions) ApplyToLVChangeOptions(new *LVChangeOptions) {
	*new = *opts
}

func (list LVChangeOptionsList) AsArgs() (Arguments, error) {
	args := NewArgs(ArgsTypeLVChange)
	options := LVChangeOptions{}
	for _, opt := range list {
		opt.ApplyToLVChangeOptions(&options)
	}
	if err := options.ApplyToArgs(args); err != nil {
		return nil, err
	}
	return args, nil
}

func (opts *LVChangeOptions) ApplyToArgs(args Arguments) error {
	id, err := NewFQLogicalVolumeName(opts.VolumeGroupName, opts.LogicalVolumeName)
	if err != nil {
		return err
	}

	for _, arg := range []Argument{
		id,
		opts.Permission,
		opts.Tags,
		opts.DelTags,
		opts.Zero,
		opts.RequestConfirm,
		opts.ActivationState,
		opts.ActivationMode,
		opts.AllocationPolicy,
		opts.ErrorWhenFull,
		opts.Partial,
		opts.SyncAction,
		opts.Rebuild,
		opts.Resync,
		opts.Discards,
		opts.Deduplication,
		opts.Compression,
		opts.AutoActivation,
		opts.CommonOptions,
	} {
		if err := arg.ApplyToArgs(args); err != nil {
			return err
		}
	}

	return nil
}
