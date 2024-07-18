package lvm2go

import (
	"context"
)

type DeleteNotFound bool

func (opt DeleteNotFound) ApplyToDevUpdateOptions(opts *DevUpdateOptions) {
	opts.DeleteNotFound = opt
}

func (opt DeleteNotFound) ApplyToArgs(args Arguments) error {
	if opt {
		args.AddOrReplace("--delete-not-found")
	}
	return nil
}

type (
	DevUpdateOptions struct {
		DevicesFile

		DeleteNotFound
		RefreshDevices
	}
	DevUpdateOption interface {
		ApplyToDevUpdateOptions(opts *DevUpdateOptions)
	}
	DevUpdateOptionsList []DevUpdateOption
)

var (
	_ ArgumentGenerator = DevUpdateOptionsList{}
	_ Argument          = (*DevUpdateOptions)(nil)
)

func (c *client) DevUpdate(ctx context.Context, opts ...DevUpdateOption) error {
	args, err := DevUpdateOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunRaw(
		ctx,
		NoOpRawOutputProcessor(false),
		append([]string{"lvmdevices", "--update"}, args.GetRaw()...)...,
	)
}

func (list DevUpdateOptionsList) AsArgs() (Arguments, error) {
	args := NewArgs(ArgsTypeGeneric)
	options := DevUpdateOptions{}
	for _, opt := range list {
		opt.ApplyToDevUpdateOptions(&options)
	}
	if err := options.ApplyToArgs(args); err != nil {
		return nil, err
	}
	return args, nil
}

func (opts *DevUpdateOptions) ApplyToArgs(args Arguments) error {
	if err := opts.DevicesFile.ApplyToArgs(args); err != nil {
		return err
	}
	if err := opts.DeleteNotFound.ApplyToArgs(args); err != nil {
		return err
	}
	if err := opts.RefreshDevices.ApplyToArgs(args); err != nil {
		return err
	}
	return nil
}
