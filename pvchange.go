package lvm2go

import (
	"context"
)

type (
	PVChangeOptions struct {
		PhysicalVolumeName
		Tags
		DelTags
		CommonOptions
	}
	PVChangeOption interface {
		ApplyToPVChangeOptions(opts *PVChangeOptions)
	}
	PVChangeOptionsList []PVChangeOption
)

var (
	_ ArgumentGenerator = PVChangeOptionsList{}
	_ Argument          = (*PVChangeOptions)(nil)
)

func (c *client) PVChange(ctx context.Context, opts ...PVChangeOption) error {
	args, err := PVChangeOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunLVM(ctx, append([]string{"pvchange"}, args.GetRaw()...)...)
}

func (opts *PVChangeOptions) ApplyToPVChangeOptions(new *PVChangeOptions) {
	*new = *opts
}

func (list PVChangeOptionsList) AsArgs() (Arguments, error) {
	args := NewArgs(ArgsTypeGeneric)
	options := PVChangeOptions{}
	for _, opt := range list {
		opt.ApplyToPVChangeOptions(&options)
	}
	if err := options.ApplyToArgs(args); err != nil {
		return nil, err
	}
	return args, nil
}

func (opts *PVChangeOptions) ApplyToArgs(args Arguments) error {
	if opts.PhysicalVolumeName == "" {
		return ErrPhysicalVolumeNameRequired
	}

	for _, arg := range []Argument{
		opts.PhysicalVolumeName,
		opts.Tags,
		opts.DelTags,
		opts.CommonOptions,
	} {
		if err := arg.ApplyToArgs(args); err != nil {
			return err
		}
	}

	return nil
}
