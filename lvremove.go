package lvm2go

import (
	"context"
	"fmt"
)

type (
	LVRemoveOptions struct {
		LogicalVolumeName
		VolumeGroupName

		Force
		Tags
		Select

		CommonOptions
	}
	LVRemoveOption interface {
		ApplyToLVRemoveOptions(opts *LVRemoveOptions)
	}
	LVRemoveOptionsList []LVRemoveOption
)

var (
	_ ArgumentGenerator = LVRemoveOptionsList{}
	_ Argument          = (*LVRemoveOptions)(nil)
)

func (c *client) LVRemove(ctx context.Context, opts ...LVRemoveOption) error {
	args, err := LVRemoveOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunLVM(ctx, append([]string{"lvremove"}, args.GetRaw()...)...)
}

func (opts *LVRemoveOptions) ApplyToArgs(args Arguments) error {
	if opts.LogicalVolumeName == "" {
		return fmt.Errorf("LogicalVolumeName is required for removal of a logical volume")
	}

	if opts.VolumeGroupName == "" {
		return fmt.Errorf("VolumeGroupName is required for removal of a logical volume")
	}

	var identifier []Argument
	fq, err := NewFQLogicalVolumeName(opts.VolumeGroupName, opts.LogicalVolumeName)
	if err != nil {
		return err
	}
	identifier = []Argument{fq}

	for _, arg := range append(identifier,
		opts.Tags,
		opts.Force,
		opts.CommonOptions,
	) {
		if err := arg.ApplyToArgs(args); err != nil {
			return err
		}
	}

	return nil
}

func (list LVRemoveOptionsList) AsArgs() (Arguments, error) {
	args := NewArgs(ArgsTypeGeneric)
	options := LVRemoveOptions{}
	for _, opt := range list {
		opt.ApplyToLVRemoveOptions(&options)
	}
	if err := options.ApplyToArgs(args); err != nil {
		return nil, err
	}
	return args, nil
}

func (opts *LVRemoveOptions) ApplyToLVRemoveOptions(new *LVRemoveOptions) {
	*new = *opts
}
