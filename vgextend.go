package lvm2go

import (
	"context"
	"fmt"
)

type (
	VGExtendOptions struct {
		VolumeGroupName
		PhysicalVolumeNames
		Force
		Zero
		CommonOptions
	}
	VGExtendOption interface {
		ApplyToVGExtendOptions(opts *VGExtendOptions)
	}
	VGExtendOptionsList []VGExtendOption
)

var (
	_ ArgumentGenerator = VGExtendOptionsList{}
	_ Argument          = (*VGExtendOptions)(nil)
	_ VGExtendOption    = (*VGExtendOptions)(nil)
)

func (c *client) VGExtend(ctx context.Context, opts ...VGExtendOption) error {
	args, err := VGExtendOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunLVM(ctx, append([]string{"vgextend"}, args.GetRaw()...)...)
}

func (opts *VGExtendOptions) ApplyToVGExtendOptions(new *VGExtendOptions) {
	*new = *opts
}

func (list VGExtendOptionsList) AsArgs() (Arguments, error) {
	args := NewArgs(ArgsTypeGeneric)
	options := VGExtendOptions{}
	for _, opt := range list {
		opt.ApplyToVGExtendOptions(&options)
	}
	if err := options.ApplyToArgs(args); err != nil {
		return nil, err
	}
	return args, nil
}

func (opts *VGExtendOptions) ApplyToArgs(args Arguments) error {
	if opts.VolumeGroupName == "" {
		return fmt.Errorf("VolumeGroupName is required for extension of a volume group")
	}

	if len(opts.PhysicalVolumeNames) == 0 {
		return fmt.Errorf("at least one PhysicalVolumeName is required for extension of a volume group")
	}

	for _, arg := range []Argument{
		opts.VolumeGroupName,
		opts.PhysicalVolumeNames,
		opts.Force,
		opts.Zero,
		opts.CommonOptions,
	} {
		if err := arg.ApplyToArgs(args); err != nil {
			return err
		}
	}

	return nil
}
