package lvm2go

import (
	"context"
	"errors"
)

type (
	LVRenameOptions struct {
		VolumeGroupName
		Old LogicalVolumeName
		New LogicalVolumeName
		CommonOptions
	}
	LVRenameOption interface {
		ApplyToLVRenameOptions(opts *LVRenameOptions)
	}
	LVRenameOptionsList []LVRenameOption
)

func (opts *LVRenameOptions) SetOldOrNew(name LogicalVolumeName) {
	if opts.Old == "" {
		opts.Old = name
	} else if opts.New == "" {
		opts.New = name
	} else {
		opts.Old = opts.New
		opts.New = name
	}
}

var (
	_ ArgumentGenerator = LVRenameOptionsList{}
	_ Argument          = (*LVRenameOptions)(nil)
)

func (c *client) LVRename(ctx context.Context, opts ...LVRenameOption) error {
	args, err := LVRenameOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunLVM(ctx, append([]string{"lvrename"}, args.GetRaw()...)...)
}

func (list LVRenameOptionsList) AsArgs() (Arguments, error) {
	args := NewArgs(ArgsTypeLVRename)
	options := LVRenameOptions{}
	for _, opt := range list {
		opt.ApplyToLVRenameOptions(&options)
	}
	if err := options.ApplyToArgs(args); err != nil {
		return nil, err
	}
	return args, nil
}

func (opts *LVRenameOptions) ApplyToArgs(args Arguments) error {
	if opts.VolumeGroupName == "" {
		return errors.New("VolumeGroupName is required")
	}
	if opts.Old == "" {
		return errors.New("old LogicalVolumeName is required")
	}
	if opts.New == "" {
		return errors.New("new LogicalVolumeName is required")
	}

	for _, arg := range []Argument{
		opts.VolumeGroupName,
		opts.Old,
		opts.New,
		opts.CommonOptions,
	} {
		if err := arg.ApplyToArgs(args); err != nil {
			return err
		}
	}

	return nil
}
