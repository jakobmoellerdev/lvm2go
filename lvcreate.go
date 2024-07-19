package lvm2go

import (
	"context"
	"fmt"
)

type (
	LVCreateOptions struct {
		VolumeGroupName
		LogicalVolumeName

		Tags
		Size
		Extents
		VirtualSize

		AllocationPolicy
		ActivationState
		Zero
		ChunkSize
		Type
		Thin
		*ThinPool

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
		return fmt.Errorf("LogicalVolumeName or ThinPoolName is required for creation of a logical volume")
	}

	if opts.Extents.Val > 0 && opts.Size.Val > 0 && opts.VirtualSize.Val > 0 {
		return fmt.Errorf("size, virtual size and extents are mutually exclusive")
	} else if opts.Extents.Val <= 0 && opts.Size.Val <= 0 && opts.VirtualSize.Val <= 0 {
		return fmt.Errorf("size, virtual size or extents must be specified")
	}

	if opts.Type == TypeThin && opts.ThinPool == nil {
		return fmt.Errorf("ThinPool is required for Thin Logical Volume")
	}

	if opts.ThinPool != nil && opts.VolumeGroupName != "" {
		return fmt.Errorf("ThinPool and VolumeGroupName are mutually exclusive. VolumeGroupName is a part of ThinPool name")
	}

	var identifier []Argument

	if opts.ThinPool != nil {
		identifier = []Argument{opts.ThinPool, opts.LogicalVolumeName}
	} else {
		identifier = []Argument{opts.VolumeGroupName, opts.LogicalVolumeName}
	}

	var sizeArgument Argument
	if opts.Extents.Val > 0 {
		sizeArgument = opts.Extents
	} else if opts.Size.Val > 0 {
		sizeArgument = opts.Size
	} else {
		sizeArgument = opts.VirtualSize
	}

	for _, arg := range append(identifier,
		sizeArgument,
		opts.AllocationPolicy,
		opts.Thin,
		opts.Type,
		opts.ActivationState,
		opts.Zero,
		opts.Tags,
		opts.CommonOptions,
	) {
		if err := arg.ApplyToArgs(args); err != nil {
			return err
		}
	}

	return nil
}

func (opts *LVCreateOptions) ApplyToLVCreateOptions(new *LVCreateOptions) {
	*new = *opts
}
