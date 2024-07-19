package lvm2go

import (
	"context"
	"errors"
	"fmt"
)

type (
	LVExtendOptions struct {
		VolumeGroupName
		LogicalVolumeName

		PoolMetadataPrefixedSize
		Size
		Extents

		CommonOptions
	}
	LVExtendOption interface {
		ApplyToLVExtendOptions(opts *LVExtendOptions)
	}
	LVExtendOptionsList []LVExtendOption
)

var (
	_ ArgumentGenerator = LVExtendOptionsList{}
	_ Argument          = (*LVExtendOptions)(nil)
	_ LVExtendOption    = (*LVExtendOptions)(nil)
)

func (c *client) LVExtend(ctx context.Context, opts ...LVExtendOption) error {
	args, err := LVExtendOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunLVM(ctx, append([]string{"lvextend"}, args.GetRaw()...)...)
}

func (opts *LVExtendOptions) ApplyToLVExtendOptions(new *LVExtendOptions) {
	*new = *opts
}

func (list LVExtendOptionsList) AsArgs() (Arguments, error) {
	args := NewArgs(ArgsTypeGeneric)
	options := LVExtendOptions{}
	for _, opt := range list {
		opt.ApplyToLVExtendOptions(&options)
	}
	if err := options.ApplyToArgs(args); err != nil {
		return nil, err
	}
	return args, nil

}

func (opts *LVExtendOptions) ApplyToArgs(args Arguments) error {
	if opts.VolumeGroupName == "" {
		return errors.New("VolumeGroupName is required")
	}
	if opts.LogicalVolumeName == "" {
		return errors.New("LogicalVolumeName is required")
	}

	if opts.Extents.Val > 0 && opts.Size.Val > 0 {
		return fmt.Errorf("size and extents are mutually exclusive")
	} else if opts.Extents.Val <= 0 && opts.Size.Val <= 0 {
		return fmt.Errorf("size or extents must be specified")
	}

	if opts.PoolMetadataPrefixedSize.Val == 0 && opts.Size.Val == 0 && opts.Extents.Val == 0 {
		return errors.New("PoolMetadataPrefixedSize, Size or Extents is required")
	}

	fqLogicalVolumeName, err := NewFQLogicalVolumeName(opts.VolumeGroupName, opts.LogicalVolumeName)
	if err != nil {
		return err
	}

	for _, arg := range []Argument{
		fqLogicalVolumeName,
		opts.Size,
		opts.PoolMetadataPrefixedSize,
		opts.Extents,
		opts.CommonOptions,
	} {
		if err := arg.ApplyToArgs(args); err != nil {
			return err
		}
	}

	return nil
}
