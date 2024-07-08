package lvm2go

import (
	"context"
	"errors"
	"fmt"
)

type (
	LVExtendOptions struct {
		LogicalVolumeName
		VolumeGroupName

		CommonOptions
	}
	LVExtendOption interface {
		ApplyToLVExtendOptions(opts *LVExtendOptions)
	}
	LVExtendOptionsList []LVExtendOption
)

var (
	_ ArgumentGenerator = LVExtendOptionsList{}
	_ Argument          = (*LVExtendOptions)(nil)
)

func (c *client) LVExtend(ctx context.Context, opts ...LVExtendOption) error {
	args, err := LVExtendOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunLVM(ctx, append([]string{"lvextend"}, args.GetRaw()...)...)
}

func (L LVExtendOptionsList) AsArgs() (Arguments, error) {
	return nil, fmt.Errorf("not implemented: %w", errors.ErrUnsupported)
}

func (opts *LVExtendOptions) ApplyToArgs(args Arguments) error {
	return fmt.Errorf("not implemented: %w", errors.ErrUnsupported)
}
