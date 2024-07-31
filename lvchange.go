package lvm2go

import (
	"context"
)

type (
	LVChangeOptions struct {
		VolumeGroupName
		LogicalVolumeName

		Permission

		Tags
		DelTags

		Zero
		RequestConfirm
		ActivationState
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

func (opts *LVChangeOptions) ApplyToLVChangeOptions(new *LVChangeOptions) {
	*new = *opts
}

func (list LVChangeOptionsList) AsArgs() (Arguments, error) {
	args := NewArgs(ArgsTypeLVChange)
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
	id, err := NewFQLogicalVolumeName(opts.VolumeGroupName, opts.LogicalVolumeName)
	if err != nil {
		return err
	}

	for _, arg := range []Argument{
		id,
		opts.Permission,
		opts.Tags,
		opts.DelTags,
		opts.Zero,
		opts.RequestConfirm,
		opts.ActivationState,
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
