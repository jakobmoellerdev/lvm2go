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
		ApplyToVGRemoveOptions(opts *VGExtendOptions)
	}
	VGExtendOptionsList []VGExtendOption
)

var (
	_ ArgumentGenerator = VGExtendOptionsList{}
	_ Argument          = (*VGExtendOptions)(nil)
)

func (c *client) VGExtend(ctx context.Context, opts ...VGExtendOption) error {
	args, err := VGExtendOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunLVM(ctx, append([]string{"vgextend"}, args.GetRaw()...)...)
}

func (L VGExtendOptionsList) AsArgs() (Arguments, error) {
	return nil, fmt.Errorf("not implemented: %w", errors.ErrUnsupported)
}

func (opts *VGExtendOptions) ApplyToArgs(args Arguments) error {
	return fmt.Errorf("not implemented: %w", errors.ErrUnsupported)
}
