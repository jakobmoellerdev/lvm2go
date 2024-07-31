package lvm2go

import (
	"context"
	"fmt"
)

type (
	VGReduceOptions struct {
		VolumeGroupName
		PhysicalVolumeNames
		RemoveMissing
		Force
		CommonOptions
	}
	VGReduceOption interface {
		ApplyToVGReduceOptions(opts *VGReduceOptions)
	}
	VGReduceOptionsList []VGReduceOption
)

var (
	_ ArgumentGenerator = VGReduceOptionsList{}
	_ Argument          = (*VGReduceOptions)(nil)
)

func (c *client) VGReduce(ctx context.Context, opts ...VGReduceOption) error {
	args, err := VGReduceOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunLVM(ctx, append([]string{"vgreduce"}, args.GetRaw()...)...)
}

func (list VGReduceOptionsList) AsArgs() (Arguments, error) {
	args := NewArgs(ArgsTypeGeneric)
	options := VGReduceOptions{}
	for _, opt := range list {
		opt.ApplyToVGReduceOptions(&options)
	}
	if err := options.ApplyToArgs(args); err != nil {
		return nil, err
	}
	return args, nil
}

func (opts *VGReduceOptions) ApplyToArgs(args Arguments) error {
	if opts.VolumeGroupName == "" {
		return fmt.Errorf("VolumeGroupName is required for extension of a volume group")
	}

	if len(opts.PhysicalVolumeNames) == 0 && !opts.RemoveMissing {
		return fmt.Errorf("at least one PhysicalVolumeName is required for reduction of a volume group")
	}

	for _, arg := range []Argument{
		opts.RemoveMissing,
		opts.VolumeGroupName,
		opts.PhysicalVolumeNames,
		opts.Force,
		opts.CommonOptions,
	} {
		if err := arg.ApplyToArgs(args); err != nil {
			return err
		}
	}

	return nil
}

func (opts *VGReduceOptions) ApplyToVGReduceOptions(new *VGReduceOptions) {
	*new = *opts
}
