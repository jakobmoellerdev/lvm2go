package lvm2go

import (
	"context"
	"errors"
	"fmt"
)

type (
	LVResizeOptions struct {
		LogicalVolumeName
		VolumeGroupName

		CommonOptions
	}
	LVResizeOption interface {
		ApplyToLVResizeOptions(opts *LVResizeOptions)
	}
	LVResizeOptionsList []LVResizeOption
)

var (
	_ ArgumentGenerator = LVResizeOptionsList{}
	_ Argument          = (*LVResizeOptions)(nil)
)

func (c *client) LVResize(ctx context.Context, opts ...LVResizeOption) error {
	args, err := LVResizeOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunLVM(ctx, append([]string{"lvresize"}, args.GetRaw()...)...)
}

func (L LVResizeOptionsList) AsArgs() (Arguments, error) {
	return nil, fmt.Errorf("not implemented: %w", errors.ErrUnsupported)
}

func (opts *LVResizeOptions) ApplyToArgs(args Arguments) error {
	return fmt.Errorf("not implemented: %w", errors.ErrUnsupported)
}
