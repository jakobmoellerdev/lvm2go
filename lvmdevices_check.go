package lvm2go

import (
	"context"
)

type (
	DevCheckOptions struct {
		DevicesFile
		RefreshDevices
	}
	DevCheckOption interface {
		ApplyToDevCheckOptions(opts *DevCheckOptions)
	}
	DevCheckOptionsList []DevCheckOption
)

var (
	_ ArgumentGenerator = DevCheckOptionsList{}
	_ Argument          = (*DevCheckOptions)(nil)
)

func (c *client) DevCheck(ctx context.Context, opts ...DevCheckOption) error {
	args, err := DevCheckOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunRaw(
		ctx,
		NoOpRawOutputProcessor(false),
		append([]string{"lvmdevices", "--check"}, args.GetRaw()...)...,
	)
}

func (list DevCheckOptionsList) AsArgs() (Arguments, error) {
	args := NewArgs(ArgsTypeGeneric)
	options := DevCheckOptions{}
	for _, opt := range list {
		opt.ApplyToDevCheckOptions(&options)
	}
	if err := options.ApplyToArgs(args); err != nil {
		return nil, err
	}
	return args, nil
}

func (opts *DevCheckOptions) ApplyToDevCheckOptions(new *DevCheckOptions) {
	*new = *opts
}

func (opts *DevCheckOptions) ApplyToArgs(args Arguments) error {
	if err := opts.DevicesFile.ApplyToArgs(args); err != nil {
		return err
	}
	if err := opts.RefreshDevices.ApplyToArgs(args); err != nil {
		return err
	}
	return nil
}
