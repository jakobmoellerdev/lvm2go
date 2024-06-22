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
		AllocationPolicy
		Activate
		Zero

		CommonOptions
	}
	LVCreateOption interface {
		ApplyToLVCreateOptions(opts *LVCreateOptions)
	}
	LVCreateOptionList []LVCreateOption
)

var (
	_ ArgumentGenerator = LVCreateOptionList{}
	_ Argument          = (*LVCreateOptions)(nil)
)

func (c *client) LVCreate(ctx context.Context, opts ...LVCreateOption) error {
	args, err := LVCreateOptionList(opts).AsArgs()
	if err != nil {
		return err
	}

	return RunLVM(ctx, append([]string{"lvcreate"}, args.GetRaw()...)...)
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

	if opts.Size.Val <= 0 {
		return fmt.Errorf("size is required for creation of a logical volume")
	}

	if err := opts.VolumeGroupName.ApplyToArgs(args); err != nil {
		return err
	}

	if err := opts.LogicalVolumeName.ApplyToArgs(args); err != nil {
		return err
	}

	if err := opts.Size.ApplyToArgs(args); err != nil {
		return err
	}

	if err := opts.AllocationPolicy.ApplyToArgs(args); err != nil {
		return err
	}

	if err := opts.Activate.ApplyToArgs(args); err != nil {
		return err
	}

	if err := opts.Zero.ApplyToArgs(args); err != nil {
		return err
	}

	if err := opts.Tags.ApplyToArgs(args); err != nil {
		return err
	}

	if err := opts.CommonOptions.ApplyToArgs(args); err != nil {
		return err
	}

	return nil
}

func (opts *LVCreateOptions) ApplyToLVCreateOptions(new *LVCreateOptions) {
	*new = *opts
}
