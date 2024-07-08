package lvm2go

import (
	"context"
	"errors"
	"fmt"
)

type (
	LVRenameOptions struct {
		VolumeGroupName
		CommonOptions
	}
	LVRenameOption interface {
		ApplyToLVRenameOptions(opts *LVRenameOptions)
	}
	LVRenameOptionsList []LVRenameOption
)

var (
	_ ArgumentGenerator = LVRenameOptionsList{}
	_ Argument          = (*LVRenameOptions)(nil)
)

func (c *client) LVRename(ctx context.Context, opts ...LVRenameOption) error {
	args, err := LVRenameOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunLVM(ctx, append([]string{"lvrename"}, args.GetRaw()...)...)
}

func (L LVRenameOptionsList) AsArgs() (Arguments, error) {
	return nil, fmt.Errorf("not implemented: %w", errors.ErrUnsupported)
}

func (opts *LVRenameOptions) ApplyToArgs(args Arguments) error {
	return fmt.Errorf("not implemented: %w", errors.ErrUnsupported)
}
