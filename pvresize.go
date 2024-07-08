package lvm2go

import (
	"context"
	"errors"
	"fmt"
)

type (
	PVResizeOptions struct {
		PhysicalVolumeName
		CommonOptions
	}
	PVResizeOption interface {
		ApplyToPVResizeOptions(opts *PVResizeOptions)
	}
	PVResizeOptionsList []PVResizeOption
)

var (
	_ ArgumentGenerator = PVResizeOptionsList{}
	_ Argument          = (*PVResizeOptions)(nil)
)

func (c *client) PVResize(ctx context.Context, opts ...PVResizeOption) error {
	args, err := PVResizeOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunLVM(ctx, append([]string{"pvresize"}, args.GetRaw()...)...)
}

func (L PVResizeOptionsList) AsArgs() (Arguments, error) {
	return nil, fmt.Errorf("not implemented: %w", errors.ErrUnsupported)
}

func (opts *PVResizeOptions) ApplyToArgs(args Arguments) error {
	return fmt.Errorf("not implemented: %w", errors.ErrUnsupported)
}
