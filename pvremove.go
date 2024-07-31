package lvm2go

import (
	"context"
)

type (
	PVRemoveOptions struct {
		PhysicalVolumeName
		Force
		CommonOptions
	}
	PVRemoveOption interface {
		ApplyToPVRemoveOptions(opts *PVRemoveOptions)
	}
	PVRemoveOptionsList []PVRemoveOption
)

var (
	_ ArgumentGenerator = PVRemoveOptionsList{}
	_ Argument          = (*PVRemoveOptions)(nil)
	_ PVRemoveOption    = (*PVRemoveOptions)(nil)
)

func (c *client) PVRemove(ctx context.Context, opts ...PVRemoveOption) error {
	args, err := PVRemoveOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunLVM(ctx, append([]string{"pvremove"}, args.GetRaw()...)...)
}

func (opts *PVRemoveOptions) ApplyToPVRemoveOptions(new *PVRemoveOptions) {
	*new = *opts
}

func (list PVRemoveOptionsList) AsArgs() (Arguments, error) {
	args := NewArgs(ArgsTypeGeneric)
	options := PVRemoveOptions{}
	for _, opt := range list {
		opt.ApplyToPVRemoveOptions(&options)
	}
	if err := options.ApplyToArgs(args); err != nil {
		return nil, err
	}
	return args, nil
}

func (opts *PVRemoveOptions) ApplyToArgs(args Arguments) error {
	if opts.PhysicalVolumeName == "" {
		return ErrPhysicalVolumeNameRequired
	}

	for _, arg := range []Argument{
		opts.PhysicalVolumeName,
		opts.Force,
		opts.CommonOptions,
	} {
		if err := arg.ApplyToArgs(args); err != nil {
			return err
		}
	}

	return nil
}
