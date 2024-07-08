package lvm2go

import (
	"context"
	"errors"
	"fmt"
)

type (
	PVCreateOptions struct {
		PhysicalVolumeName
		CommonOptions
	}
	PVCreateOption interface {
		ApplyToPVCreateOptions(opts *PVCreateOptions)
	}
	PVCreateOptionsList []PVCreateOption
)

var (
	_ ArgumentGenerator = PVCreateOptionsList{}
	_ Argument          = (*PVCreateOptions)(nil)
)

func (c *client) PVCreate(ctx context.Context, opts ...PVCreateOption) error {
	args, err := PVCreateOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunLVM(ctx, append([]string{"pvcreate"}, args.GetRaw()...)...)
}

func (L PVCreateOptionsList) AsArgs() (Arguments, error) {
	return nil, fmt.Errorf("not implemented: %w", errors.ErrUnsupported)
}

func (opts *PVCreateOptions) ApplyToArgs(args Arguments) error {
	return fmt.Errorf("not implemented: %w", errors.ErrUnsupported)
}
