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
	PVCreateOptions struct {
		PhysicalVolumeName
		Force
		Zero
		DataAlignment
		DataAlignmentOffset
		MetadataSize
		CommonOptions
	}
	PVCreateOption interface {
		ApplyToPVCreateOptions(opts *PVCreateOptions)
	}
	PVCreateOptionsList []PVCreateOption
)

var (
	_ ArgumentGenerator = PVCreateOptionsList{}
	_ Argument          = (*PVCreateOptions)(nil)
)

func (c *client) PVCreate(ctx context.Context, opts ...PVCreateOption) error {
	args, err := PVCreateOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunLVM(ctx, append([]string{"pvcreate"}, args.GetRaw()...)...)
}

func (opts *PVCreateOptions) ApplyToPVCreateOptions(new *PVCreateOptions) {
	*new = *opts
}

func (list PVCreateOptionsList) AsArgs() (Arguments, error) {
	args := NewArgs(ArgsTypeGeneric)
	options := PVCreateOptions{}
	for _, opt := range list {
		opt.ApplyToPVCreateOptions(&options)
	}
	if err := options.ApplyToArgs(args); err != nil {
		return nil, err
	}
	return args, nil
}

func (opts *PVCreateOptions) ApplyToArgs(args Arguments) error {
	if opts.PhysicalVolumeName == "" {
		return ErrPhysicalVolumeNameRequired
	}

	for _, arg := range []Argument{
		opts.PhysicalVolumeName,
		opts.Force,
		opts.Zero,
		opts.DataAlignment,
		opts.DataAlignmentOffset,
		opts.MetadataSize,
		opts.CommonOptions,
	} {
		if err := arg.ApplyToArgs(args); err != nil {
			return err
		}
	}

	return nil
}
