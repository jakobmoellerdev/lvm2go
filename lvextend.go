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

type (
	LVExtendOptions struct {
		VolumeGroupName
		LogicalVolumeName

		PoolMetadataPrefixedSize
		PrefixedSize
		PrefixedExtents

		CommonOptions
	}
	LVExtendOption interface {
		ApplyToLVExtendOptions(opts *LVExtendOptions)
	}
	LVExtendOptionsList []LVExtendOption
)

var (
	_ ArgumentGenerator = LVExtendOptionsList{}
	_ Argument          = (*LVExtendOptions)(nil)
	_ LVExtendOption    = (*LVExtendOptions)(nil)
)

func (c *client) LVExtend(ctx context.Context, opts ...LVExtendOption) error {
	args, err := LVExtendOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunLVM(ctx, append([]string{"lvextend"}, args.GetRaw()...)...)
}

func (opts *LVExtendOptions) ApplyToLVExtendOptions(new *LVExtendOptions) {
	*new = *opts
}

func (list LVExtendOptionsList) AsArgs() (Arguments, error) {
	args := NewArgs(ArgsTypeGeneric)
	options := LVExtendOptions{}
	for _, opt := range list {
		opt.ApplyToLVExtendOptions(&options)
	}
	if err := options.ApplyToArgs(args); err != nil {
		return nil, err
	}
	return args, nil

}

func (opts *LVExtendOptions) ApplyToArgs(args Arguments) error {
	id, err := NewFQLogicalVolumeName(opts.VolumeGroupName, opts.LogicalVolumeName)
	if err != nil {
		return err
	}

	if opts.Extents.Val > 0 && opts.PrefixedSize.Val > 0 {
		return fmt.Errorf("size and extents are mutually exclusive")
	} else if opts.Extents.Val <= 0 && opts.PrefixedSize.Val <= 0 {
		return fmt.Errorf("size or extents must be specified")
	}

	if opts.PrefixedSize.SizePrefix == SizePrefixMinus {
		return fmt.Errorf("size prefix must be positive")
	} else if opts.PrefixedExtents.SizePrefix == SizePrefixMinus {
		return fmt.Errorf("extents prefix must be positive")
	} else if opts.PoolMetadataPrefixedSize.SizePrefix == SizePrefixMinus {
		return fmt.Errorf("pool metadata size prefix must be positive")
	}

	if opts.PoolMetadataPrefixedSize.Val == 0 && opts.PrefixedSize.Val == 0 && opts.Extents.Val == 0 {
		return errors.New("PoolMetadataPrefixedSize, Size or Extents is required")
	}

	for _, arg := range []Argument{
		id,
		opts.PrefixedSize,
		opts.PrefixedExtents,
		opts.PoolMetadataPrefixedSize,
		opts.CommonOptions,
	} {
		if err := arg.ApplyToArgs(args); err != nil {
			return err
		}
	}

	return nil
}
