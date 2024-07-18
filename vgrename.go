package lvm2go

import (
	"context"
	"errors"
	"fmt"
)

type (
	VGRenameOptions struct {
		VolumeGroupName
		CommonOptions
	}
	VGRenameOption interface {
		ApplyToVGRenameOptions(opts *VGRenameOptions)
	}
	VGRenameOptionsList []VGRenameOption
)

var (
	_ ArgumentGenerator = VGRenameOptionsList{}
	_ Argument          = (*VGRenameOptions)(nil)
)

func (c *client) VGRename(ctx context.Context, opts ...VGRenameOption) error {
	args, err := VGRenameOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunLVM(ctx, append([]string{"vgrename"}, args.GetRaw()...)...)
}

func (list VGRenameOptionsList) AsArgs() (Arguments, error) {
	return nil, fmt.Errorf("not implemented: %w", errors.ErrUnsupported)
}

func (opts *VGRenameOptions) ApplyToVGRenameOptions(new *VGRenameOptions) {
	*new = *opts
}

func (opts *VGRenameOptions) ApplyToArgs(args Arguments) error {
	return fmt.Errorf("not implemented: %w", errors.ErrUnsupported)
}
