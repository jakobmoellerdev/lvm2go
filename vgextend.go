package lvm2go

import (
	"context"
	"errors"
	"fmt"
)

type (
	VGExtendOptions struct {
		VolumeGroupName
		CommonOptions
	}
	VGExtendOption interface {
		ApplyToVGExtendOptions(opts *VGExtendOptions)
	}
	VGExtendOptionsList []VGExtendOption
)

var (
	_ ArgumentGenerator = VGExtendOptionsList{}
	_ Argument          = (*VGExtendOptions)(nil)
	_ VGExtendOption    = (*VGExtendOptions)(nil)
)

func (c *client) VGExtend(ctx context.Context, opts ...VGExtendOption) error {
	args, err := VGExtendOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunLVM(ctx, append([]string{"vgextend"}, args.GetRaw()...)...)
}

func (opts *VGExtendOptions) ApplyToVGExtendOptions(new *VGExtendOptions) {
	*new = *opts
}

func (list VGExtendOptionsList) AsArgs() (Arguments, error) {
	return nil, fmt.Errorf("not implemented: %w", errors.ErrUnsupported)
}

func (opts *VGExtendOptions) ApplyToArgs(args Arguments) error {
	if opts.VolumeGroupName == "" {
		return fmt.Errorf("VolumeGroupName is required for extension of a volume group")
	}

	if err := opts.VolumeGroupName.ApplyToArgs(args); err != nil {
		return err
	}

	if err := opts.CommonOptions.ApplyToArgs(args); err != nil {
		return err
	}

	return nil
}
