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
	VGRenameOptions struct {
		Old VolumeGroupName
		New VolumeGroupName
		Force
		CommonOptions
	}
	VGRenameOption interface {
		ApplyToVGRenameOptions(opts *VGRenameOptions)
	}
	VGRenameOptionsList []VGRenameOption
)

func (opts *VGRenameOptions) SetOldOrNew(name VolumeGroupName) {
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
	_ ArgumentGenerator = VGRenameOptionsList{}
	_ Argument          = (*VGRenameOptions)(nil)
)

func (c *client) VGRename(ctx context.Context, opts ...VGRenameOption) error {
	args, err := VGRenameOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunLVM(ctx, append([]string{"vgrename"}, args.GetRaw()...)...)
}

func (list VGRenameOptionsList) AsArgs() (Arguments, error) {
	args := NewArgs(ArgsTypeGeneric)
	options := VGRenameOptions{}
	for _, opt := range list {
		opt.ApplyToVGRenameOptions(&options)
	}
	if err := options.ApplyToArgs(args); err != nil {
		return nil, err
	}
	return args, nil
}

func (opts *VGRenameOptions) ApplyToVGRenameOptions(new *VGRenameOptions) {
	*new = *opts
}

func (opts *VGRenameOptions) ApplyToArgs(args Arguments) error {
	if opts.Old == "" {
		return fmt.Errorf("old is empty: %w", ErrVolumeGroupNameRequired)
	}
	if opts.New == "" {
		return fmt.Errorf("new is empty: %w", ErrVolumeGroupNameRequired)
	}

	for _, arg := range []Argument{
		opts.Old,
		opts.New,
		opts.Force,
		opts.CommonOptions,
	} {
		if err := arg.ApplyToArgs(args); err != nil {
			return err
		}
	}

	return nil
}
