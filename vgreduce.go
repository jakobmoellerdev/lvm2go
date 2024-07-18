package lvm2go

import (
	"context"
	"errors"
	"fmt"
)

type (
	VGReduceOptions struct {
		VolumeGroupName
		CommonOptions
	}
	VGReduceOption interface {
		ApplyToVGReduceOptions(opts *VGReduceOptions)
	}
	VGReduceOptionsList []VGReduceOption
)

var (
	_ ArgumentGenerator = VGReduceOptionsList{}
	_ Argument          = (*VGReduceOptions)(nil)
)

func (c *client) VGReduce(ctx context.Context, opts ...VGReduceOption) error {
	args, err := VGReduceOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunLVM(ctx, append([]string{"vgreduce"}, args.GetRaw()...)...)
}

func (list VGReduceOptionsList) AsArgs() (Arguments, error) {
	return nil, fmt.Errorf("not implemented: %w", errors.ErrUnsupported)
}

func (opts *VGReduceOptions) ApplyToArgs(args Arguments) error {
	return fmt.Errorf("not implemented: %w", errors.ErrUnsupported)
}
