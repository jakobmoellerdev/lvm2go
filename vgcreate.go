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

		Force
		Zero

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

	return RunLVM(ctx, append([]string{"vgcreate"}, args.GetRaw()...)...)
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

	if err := opts.VolumeGroupName.ApplyToArgs(args); err != nil {
		return err
	}

	if err := opts.PhysicalVolumeNames.ApplyToArgs(args); err != nil {
		return err
	}

	if err := opts.Tags.ApplyToArgs(args); err != nil {
		return err
	}

	if err := opts.Force.ApplyToArgs(args); err != nil {
		return err
	}

	if err := opts.Zero.ApplyToArgs(args); err != nil {
		return err
	}

	if err := opts.CommonOptions.ApplyToArgs(args); err != nil {
		return err
	}

	return nil
}

func (opts *VGCreateOptions) ApplyToVGCreateOptions(new *VGCreateOptions) {
	*new = *opts
}
