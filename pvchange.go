package lvm2go

import (
	"context"
	"errors"
	"fmt"
)

type (
	PVChangeOptions struct {
		PhysicalVolumeName
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

func (list PVChangeOptionsList) AsArgs() (Arguments, error) {
	return nil, fmt.Errorf("not implemented: %w", errors.ErrUnsupported)
}

func (opts *PVChangeOptions) ApplyToArgs(args Arguments) error {
	return fmt.Errorf("not implemented: %w", errors.ErrUnsupported)
}
