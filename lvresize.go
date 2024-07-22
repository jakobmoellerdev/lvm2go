package lvm2go

import (
	"context"
)

type (
	LVResizeOptions struct {
		LogicalVolumeName
		VolumeGroupName

		PrefixedSize

		CommonOptions
	}
	LVResizeOption interface {
		ApplyToLVResizeOptions(opts *LVResizeOptions)
	}
	LVResizeOptionsList []LVResizeOption
)

var (
	_ ArgumentGenerator = LVResizeOptionsList{}
	_ Argument          = (*LVResizeOptions)(nil)
)

func (c *client) LVResize(ctx context.Context, opts ...LVResizeOption) error {
	args, err := LVResizeOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunLVM(ctx, append([]string{"lvresize"}, args.GetRaw()...)...)
}

func (list LVResizeOptionsList) AsArgs() (Arguments, error) {
	args := NewArgs(ArgsTypeGeneric)
	options := LVResizeOptions{}
	for _, opt := range list {
		opt.ApplyToLVResizeOptions(&options)
	}
	if err := options.ApplyToArgs(args); err != nil {
		return nil, err
	}
	return args, nil
}

func (opts *LVResizeOptions) ApplyToArgs(args Arguments) error {
	id, err := NewFQLogicalVolumeName(opts.VolumeGroupName, opts.LogicalVolumeName)
	if err != nil {
		return err
	}

	for _, opt := range []Argument{
		id,
		opts.PrefixedSize,
		opts.CommonOptions,
	} {
		if err := opt.ApplyToArgs(args); err != nil {
			return err
		}
	}

	return nil
}
