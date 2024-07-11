package lvm2go

import (
	"context"
	"fmt"
)

type (
	LVCreateOptions struct {
		LogicalVolumeName
		VolumeGroupName
		Tags

		Size
		Extents
		VirtualSize

		AllocationPolicy
		Activate
		*Zero
		ChunkSize
		Type
		Thin

		Stripes
		Mirrors
		StripeSize

		CommonOptions
	}
	LVCreateOption interface {
		ApplyToLVCreateOptions(opts *LVCreateOptions)
	}
	LVCreateOptionList []LVCreateOption
)

func (list LVCreateOptionList) ApplyToLVCreateOptions(opts *LVCreateOptions) {
	for _, opt := range list {
		opt.ApplyToLVCreateOptions(opts)
	}
}

var (
	_ ArgumentGenerator = LVCreateOptionList{}
	_ Argument          = (*LVCreateOptions)(nil)
)

func (c *client) LVCreate(ctx context.Context, opts ...LVCreateOption) error {
	args, err := LVCreateOptionList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunLVM(ctx, append([]string{"lvcreate"}, args.GetRaw()...)...)
}

func (list LVCreateOptionList) AsArgs() (Arguments, error) {
	args := NewArgs(ArgsTypeLVCreate)
	options := LVCreateOptions{}
	for _, opt := range list {
		opt.ApplyToLVCreateOptions(&options)
	}
	if err := options.ApplyToArgs(args); err != nil {
		return nil, err
	}
	return args, nil
}

func (opts *LVCreateOptions) ApplyToArgs(args Arguments) error {
	if opts.LogicalVolumeName == "" {
		return fmt.Errorf("LogicalVolumeName is required for creation of a logical volume")
	}

	if opts.VolumeGroupName == "" {
		return fmt.Errorf("VolumeGroupName is required for creation of a logical volume")
	}

	if opts.Extents.Val > 0 && opts.Size.Val > 0 {
		return fmt.Errorf("size and extents are mutually exclusive")
	} else if opts.Extents.Val <= 0 && opts.Size.Val <= 0 {
		return fmt.Errorf("size or extents must be specified")
	}

	var sizeArgument Argument
	if opts.Extents.Val > 0 {
		sizeArgument = opts.Extents
	} else {
		sizeArgument = opts.Size
	}

	for _, arg := range []Argument{
		opts.VolumeGroupName,
		opts.LogicalVolumeName,
		sizeArgument,
		opts.AllocationPolicy,
		opts.Activate,
		opts.Zero,
		opts.Tags,
		opts.CommonOptions,
	} {
		if err := arg.ApplyToArgs(args); err != nil {
			return err
		}
	}

	return nil
}

func (opts *LVCreateOptions) ApplyToLVCreateOptions(new *LVCreateOptions) {
	*new = *opts
}
