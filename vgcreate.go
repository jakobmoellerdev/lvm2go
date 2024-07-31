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
