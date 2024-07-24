package lvm2go

import (
	"context"
	"fmt"
)

type (
	VGRemoveOptions struct {
		VolumeGroupName
		Tags
		Select

		Force

		CommonOptions
	}
	VGRemoveOption interface {
		ApplyToVGRemoveOptions(opts *VGRemoveOptions)
	}
	VGRemoveOptionsList []VGRemoveOption
)

var (
	_ ArgumentGenerator = VGRemoveOptionsList{}
	_ Argument          = (*VGRemoveOptions)(nil)
)

func (c *client) VGRemove(ctx context.Context, opts ...VGRemoveOption) error {
	args, err := VGRemoveOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunLVM(ctx, append([]string{"vgremove"}, args.GetRaw()...)...)
}

func (opts *VGRemoveOptions) ApplyToArgs(args Arguments) error {
	if opts.VolumeGroupName == "" {
		return fmt.Errorf("VolumeGroupName is required for removal of a volume group")
	}

	for _, arg := range []Argument{
		opts.VolumeGroupName,
		opts.Tags,
		opts.Force,
		opts.CommonOptions,
	} {
		if err := arg.ApplyToArgs(args); err != nil {
			return err
		}
	}

	return nil
}

func (opts *VGRemoveOptions) ApplyToVGRemoveOptions(new *VGRemoveOptions) {
	*new = *opts
}

var (
	_ ArgumentGenerator = VGRemoveOptionsList{}
	_ Argument          = (*VGRemoveOptions)(nil)
)

func (list VGRemoveOptionsList) AsArgs() (Arguments, error) {
	args := NewArgs(ArgsTypeGeneric)
	options := VGRemoveOptions{}
	for _, opt := range list {
		opt.ApplyToVGRemoveOptions(&options)
	}
	if err := options.ApplyToArgs(args); err != nil {
		return nil, err
	}
	return args, nil
}
