package lvm2go

import (
	"context"
)

type (
	VGsOptions struct {
		VolumeGroupName
		Tags
		Unit
		Select

		ColumnOptions
		CommonOptions
	}
	VGsOption interface {
		ApplyToVGsOptions(opts *VGsOptions)
	}
	VGsOptionsList []VGsOption
)

var (
	_ ArgumentGenerator = VGsOptionsList{}
	_ Argument          = (*VGsOptions)(nil)
)

func (c *client) VGs(ctx context.Context, opts ...VGsOption) ([]*VolumeGroup, error) {
	type vgReport struct {
		Report []struct {
			VG []*VolumeGroup `json:"vg"`
		} `json:"report"`
	}
	res := new(vgReport)

	args := []string{
		"vgs", "--reportformat", "json",
	}
	argsFromOpts, err := VGsOptionsList(opts).AsArgs()
	if err != nil {
		return nil, err
	}

	err = c.RunLVMInto(ctx, res, append(args, argsFromOpts.GetRaw()...)...)

	if IsLVMNotFound(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if len(res.Report) == 0 {
		return nil, nil
	}

	if len(res.Report[0].VG) == 0 {
		return nil, nil
	}

	return res.Report[0].VG, nil
}

func (c *client) VG(ctx context.Context, opts ...VGsOption) (*VolumeGroup, error) {
	found := false
	for _, opt := range opts {
		if _, ok := opt.(VolumeGroupName); ok {
			found = true
			break
		}
	}
	if !found {
		return nil, ErrVolumeGroupNameRequired
	}

	vgs, err := c.VGs(ctx, opts...)
	if err != nil {
		return nil, err
	}

	if len(vgs) == 0 {
		return nil, ErrVolumeGroupNotFound
	}

	return vgs[0], nil
}

func (opts *VGsOptions) ApplyToArgs(args Arguments) error {
	for _, arg := range []Argument{
		opts.VolumeGroupName,
		opts.Tags,
		opts.Unit,
		opts.CommonOptions,
		opts.ColumnOptions,
		opts.Select,
	} {
		if err := arg.ApplyToArgs(args); err != nil {
			return err
		}
	}

	return nil
}

func (opts *VGsOptions) ApplyToVGsOptions(new *VGsOptions) {
	*new = *opts
}

func (list VGsOptionsList) AsArgs() (Arguments, error) {
	args := NewArgs(ArgsTypeVGs)
	var options VGsOptions
	for _, opt := range list {
		opt.ApplyToVGsOptions(&options)
	}
	if err := options.ApplyToArgs(args); err != nil {
		return nil, err
	}
	return args, nil
}
