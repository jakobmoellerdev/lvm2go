package lvm2go

import (
	"context"
	"fmt"
)

type (
	LVRemoveOptions struct {
		LogicalVolumeName
		VolumeGroupName
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

	return RunLVM(ctx, append([]string{"lvremove"}, args.GetRaw()...)...)
}

func (opts *LVRemoveOptions) ApplyToArgs(args Arguments) error {
	if opts.VolumeGroupName == "" {
		return fmt.Errorf("VolumeGroupName is required for removal of a volume group")
	}

	if opts.LogicalVolumeName == "" {
		return fmt.Errorf("LogicalVolumeName is required for removal of a logical volume")
	}

	logicalVolumeName, err := NewFQLogicalVolumeName(opts.VolumeGroupName, opts.LogicalVolumeName)
	if err != nil {
		return err
	}

	for _, arg := range []Argument{
		logicalVolumeName,
		opts.Tags,
		Force(true), // Force is required for removal without confirmation
		opts.CommonOptions,
	} {
		if err := arg.ApplyToArgs(args); err != nil {
			return err
		}
	}

	return nil
}

func (list LVRemoveOptionsList) AsArgs() (Arguments, error) {
	args := NewArgs(ArgsTypeLVRemove)
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
