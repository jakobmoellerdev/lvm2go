package lvm2go

import (
	"context"
	"errors"
	"fmt"
)

type (
	PVRemoveOptions struct {
		PhysicalVolumeName
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
)

func (c *client) PVRemove(ctx context.Context, opts ...PVRemoveOption) error {
	args, err := PVRemoveOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunLVM(ctx, append([]string{"pvremove"}, args.GetRaw()...)...)
}

func (L PVRemoveOptionsList) AsArgs() (Arguments, error) {
	return nil, fmt.Errorf("not implemented: %w", errors.ErrUnsupported)
}

func (opts *PVRemoveOptions) ApplyToArgs(args Arguments) error {
	return fmt.Errorf("not implemented: %w", errors.ErrUnsupported)
}
