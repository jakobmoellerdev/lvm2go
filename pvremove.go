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
	return nil, fmt.Errorf("not implemented: %w", errors.ErrUnsupported)
}

func (opts *PVRemoveOptions) ApplyToArgs(args Arguments) error {
	if opts.PhysicalVolumeName == "" {
		return fmt.Errorf("PhysicalVolumeName is required for removal of a physical volume")
	}

	for _, arg := range []Argument{
		opts.PhysicalVolumeName,
		opts.CommonOptions,
	} {
		if err := arg.ApplyToArgs(args); err != nil {
			return err
		}
	}

	return nil
}
