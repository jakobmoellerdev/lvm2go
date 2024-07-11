package lvm2go

import (
	"context"
	"errors"
)

type (
	LVChangeOptions struct {
		VolumeGroupName
		LogicalVolumeName
		FQLogicalVolumeName

		Permission

		Tags
		DelTags

		*Zero
		Yes
		Activate
		ActivationMode
		AllocationPolicy
		*ErrorWhenFull
		Partial
		SyncAction
		Rebuild
		Resync
		Discards
		*Deduplication
		*Compression
		AutoActivation

		CommonOptions
	}
	LVChangeOption interface {
		ApplyToLVChangeOptions(opts *LVChangeOptions)
	}
	LVChangeOptionsList []LVChangeOption
)

var (
	_ ArgumentGenerator = LVChangeOptionsList{}
	_ Argument          = (*LVChangeOptions)(nil)
)

func (c *client) LVChange(ctx context.Context, opts ...LVChangeOption) error {
	args, err := LVChangeOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunLVM(ctx, append([]string{"lvchange"}, args.GetRaw()...)...)
}

func (list LVChangeOptionsList) AsArgs() (Arguments, error) {
	args := NewArgs(ArgsTypeGeneric)
	options := LVChangeOptions{}
	for _, opt := range list {
		opt.ApplyToLVChangeOptions(&options)
	}
	if err := options.ApplyToArgs(args); err != nil {
		return nil, err
	}
	return args, nil
}

func (opts *LVChangeOptions) ApplyToArgs(args Arguments) error {
	if opts.FQLogicalVolumeName == "" {
		if opts.VolumeGroupName == "" {
			return errors.New("VolumeGroupName is required")
		}
		if opts.LogicalVolumeName == "" {
			return errors.New("LogicalVolumeName is required")
		}
	} else {
		var err error
		opts.FQLogicalVolumeName, err = NewFQLogicalVolumeName(opts.VolumeGroupName, opts.LogicalVolumeName)
		if err != nil {
			return err
		}
	}

	for _, arg := range []Argument{
		opts.FQLogicalVolumeName,
		opts.Permission,
		opts.Tags,
		opts.DelTags,
		opts.Zero,
		opts.Yes,
		opts.Activate,
		opts.ActivationMode,
		opts.AllocationPolicy,
		opts.ErrorWhenFull,
		opts.Partial,
		opts.SyncAction,
		opts.Rebuild,
		opts.Resync,
		opts.Discards,
		opts.Deduplication,
		opts.Compression,
		opts.AutoActivation,
		opts.CommonOptions,
	} {
		if err := arg.ApplyToArgs(args); err != nil {
			return err
		}
	}

	return nil
}
