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
	LVRenameOptions struct {
		VolumeGroupName
		Old LogicalVolumeName
		New LogicalVolumeName
		CommonOptions
	}
	LVRenameOption interface {
		ApplyToLVRenameOptions(opts *LVRenameOptions)
	}
	LVRenameOptionsList []LVRenameOption
)

func (opts *LVRenameOptions) SetOldOrNew(name LogicalVolumeName) {
	if opts.Old == "" {
		opts.Old = name
	} else if opts.New == "" {
		opts.New = name
	} else {
		opts.Old = opts.New
		opts.New = name
	}
}

var (
	_ ArgumentGenerator = LVRenameOptionsList{}
	_ Argument          = (*LVRenameOptions)(nil)
)

func (c *client) LVRename(ctx context.Context, opts ...LVRenameOption) error {
	args, err := LVRenameOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunLVM(ctx, append([]string{"lvrename"}, args.GetRaw()...)...)
}

func (list LVRenameOptionsList) AsArgs() (Arguments, error) {
	args := NewArgs(ArgsTypeLVRename)
	options := LVRenameOptions{}
	for _, opt := range list {
		opt.ApplyToLVRenameOptions(&options)
	}
	if err := options.ApplyToArgs(args); err != nil {
		return nil, err
	}
	return args, nil
}

func (opts *LVRenameOptions) ApplyToLVRenameOptions(other *LVRenameOptions) {
	*other = *opts
}

func (opts *LVRenameOptions) ApplyToArgs(args Arguments) error {
	if opts.VolumeGroupName == "" {
		return ErrVolumeGroupNameRequired
	}
	if opts.Old == "" {
		return fmt.Errorf("old is empty: %w", ErrLogicalVolumeNameRequired)
	}
	if opts.New == "" {
		return fmt.Errorf("new is empty: %w", ErrLogicalVolumeNameRequired)
	}

	for _, arg := range []Argument{
		opts.VolumeGroupName,
		opts.Old,
		opts.New,
		opts.CommonOptions,
	} {
		if err := arg.ApplyToArgs(args); err != nil {
			return err
		}
	}

	return nil
}
