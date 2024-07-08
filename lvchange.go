package lvm2go

import (
	"context"
	"errors"
	"fmt"
)

type (
	LVChangeOptions struct {
		VolumeGroupName
		CommonOptions
	}
	LVChangeOption interface {
		ApplyToVGRemoveOptions(opts *LVChangeOptions)
	}
	LVChangeOptionsList []LVChangeOption
)

var (
	_ ArgumentGenerator = LVChangeOptionsList{}
	_ Argument          = (*LVChangeOptions)(nil)
)

func (c *client) LVChange(ctx context.Context, opts ...LVChangeOption) error {
	args, err := LVChangeOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunLVM(ctx, append([]string{"lvchange"}, args.GetRaw()...)...)
}

func (L LVChangeOptionsList) AsArgs() (Arguments, error) {
	return nil, fmt.Errorf("not implemented: %w", errors.ErrUnsupported)
}

func (opts *LVChangeOptions) ApplyToArgs(args Arguments) error {
	return fmt.Errorf("not implemented: %w", errors.ErrUnsupported)
}
