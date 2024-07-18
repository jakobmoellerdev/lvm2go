package lvm2go

import (
	"context"
	"errors"
	"fmt"
)

type (
	VGChangeOptions struct {
		VolumeGroupName

		Tags
		DelTags

		CommonOptions
	}
	VGChangeOption interface {
		ApplyToVGChangeOptions(opts *VGChangeOptions)
	}
	VGChangeOptionsList []VGChangeOption
)

var (
	_ ArgumentGenerator = VGChangeOptionsList{}
	_ Argument          = (*VGChangeOptions)(nil)
)

func (c *client) VGChange(ctx context.Context, opts ...VGChangeOption) error {
	args, err := VGChangeOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunLVM(ctx, append([]string{"vgchange"}, args.GetRaw()...)...)
}

func (list VGChangeOptionsList) AsArgs() (Arguments, error) {
	return nil, fmt.Errorf("not implemented: %w", errors.ErrUnsupported)
}

func (opts *VGChangeOptions) ApplyToArgs(args Arguments) error {
	return fmt.Errorf("not implemented: %w", errors.ErrUnsupported)
}
