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
	VGCreateOptions struct {
		VolumeGroupName
		Tags

		PhysicalVolumeNames

		MaximumLogicalVolumes
		MaximumPhysicalVolumes

		AutoActivation
		Force
		Zero
		PhysicalExtentSize
		DataAlignment
		DataAlignmentOffset
		MetadataSize
		AllocationPolicy

		CommonOptions
	}
	VGCreateOption interface {
		ApplyToVGCreateOptions(opts *VGCreateOptions)
	}
	VGCreateOptionList []VGCreateOption
)

var (
	_ ArgumentGenerator = VGCreateOptionList{}
	_ Argument          = (*VGCreateOptions)(nil)
)

func (c *client) VGCreate(ctx context.Context, opts ...VGCreateOption) error {
	args, err := VGCreateOptionList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunLVM(ctx, append([]string{"vgcreate"}, args.GetRaw()...)...)
}

func (list VGCreateOptionList) AsArgs() (Arguments, error) {
	args := NewArgs(ArgsTypeVGCreate)
	options := VGCreateOptions{}
	for _, opt := range list {
		opt.ApplyToVGCreateOptions(&options)
	}
	if err := options.ApplyToArgs(args); err != nil {
		return nil, err
	}
	return args, nil
}

func (opts *VGCreateOptions) ApplyToArgs(args Arguments) error {
	if opts.VolumeGroupName == "" {
		return fmt.Errorf("VolumeGroupName is required for creation of a volume group")
	}

	if len(opts.PhysicalVolumeNames) == 0 {
		return fmt.Errorf("PhysicalVolumeNames is required for creation of a volume group")
	}

	for _, opt := range []Argument{
		opts.VolumeGroupName,
		opts.PhysicalVolumeNames,
		opts.MaximumLogicalVolumes,
		opts.MaximumPhysicalVolumes,
		opts.Tags,
		opts.Force,
		opts.Zero,
		opts.PhysicalExtentSize,
		opts.DataAlignment,
		opts.DataAlignmentOffset,
		opts.MetadataSize,
		opts.AllocationPolicy,
		opts.AutoActivation,
		opts.CommonOptions,
	} {
		if err := opt.ApplyToArgs(args); err != nil {
			return err

		}
	}

	return nil
}

func (opts *VGCreateOptions) ApplyToVGCreateOptions(new *VGCreateOptions) {
	*new = *opts
}
