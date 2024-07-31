package lvm2go

import (
	"context"
	"fmt"
)

type (
	VGChangeOptions struct {
		VolumeGroupName

		MaximumLogicalVolumes
		MaximumPhysicalVolumes
		AllocationPolicy
		AutoActivation
		Tags
		DelTags

		CommonOptions
	}
	VGChangeOption interface {
		ApplyToVGChangeOptions(opts *VGChangeOptions)
	}
	VGChangeOptionsList []VGChangeOption
)

var (
	_ ArgumentGenerator = VGChangeOptionsList{}
	_ Argument          = (*VGChangeOptions)(nil)
)

func (c *client) VGChange(ctx context.Context, opts ...VGChangeOption) error {
	args, err := VGChangeOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunLVM(ctx, append([]string{"vgchange"}, args.GetRaw()...)...)
}

func (list VGChangeOptionsList) AsArgs() (Arguments, error) {
	args := NewArgs(ArgsTypeVGChange)
	options := VGChangeOptions{}
	for _, opt := range list {
		opt.ApplyToVGChangeOptions(&options)
	}
	if err := options.ApplyToArgs(args); err != nil {
		return nil, err
	}
	return args, nil
}

func (opts *VGChangeOptions) ApplyToArgs(args Arguments) error {
	if opts.VolumeGroupName == "" {
		return fmt.Errorf("VolumeGroupName is required for creation of a volume group")
	}

	for _, opt := range []Argument{
		opts.VolumeGroupName,
		opts.MaximumLogicalVolumes,
		opts.MaximumPhysicalVolumes,
		opts.AllocationPolicy,
		opts.AutoActivation,
		opts.Tags,
		opts.DelTags,
		opts.CommonOptions,
	} {
		if err := opt.ApplyToArgs(args); err != nil {
			return err

		}
	}

	return nil
}

func (opts *VGChangeOptions) ApplyToVGChangeOptions(new *VGChangeOptions) {
	*new = *opts
}
