package lvm2go

import (
	"context"
	"errors"
	"fmt"
)

type (
	LVReduceOptions struct {
		VolumeGroupName
		CommonOptions
	}
	LVReduceOption interface {
		ApplyToLVReduceOptions(opts *LVReduceOptions)
	}
	LVReduceOptionsList []LVReduceOption
)

var (
	_ ArgumentGenerator = LVReduceOptionsList{}
	_ Argument          = (*LVReduceOptions)(nil)
)

func (c *client) LVReduce(ctx context.Context, opts ...LVReduceOption) error {
	args, err := LVReduceOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunLVM(ctx, append([]string{"lvreduce"}, args.GetRaw()...)...)
}

func (L LVReduceOptionsList) AsArgs() (Arguments, error) {
	return nil, fmt.Errorf("not implemented: %w", errors.ErrUnsupported)
}

func (opts *LVReduceOptions) ApplyToArgs(args Arguments) error {
	return fmt.Errorf("not implemented: %w", errors.ErrUnsupported)
}
