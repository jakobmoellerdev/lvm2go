package lvm2go

import (
	"context"
)

type (
	LVsOptions struct {
		VolumeGroupName
		Tags
		Select

		ColumnOptions
		CommonOptions
	}
	LVsOption interface {
		ApplyToLVsOptions(opts *LVsOptions)
	}
	LVsOptionsList []LVsOption
)

var (
	_ ArgumentGenerator = LVsOptionsList{}
	_ Argument          = (*LVsOptions)(nil)
)

// LVs returns a list of logical volumes that match the given options.
// If no logical volumes are found, nil is returned.
// It is really just a wrapper around the `lvs --reportformat json` command.
func (c *client) LVs(ctx context.Context, opts ...LVsOption) ([]LogicalVolume, error) {
	type lvReport struct {
		Report []struct {
			LV []LogicalVolume `json:"lv"`
		} `json:"report"`
	}

	var res = new(lvReport)

	args := []string{
		"lvs", "--reportformat", "json",
	}
	argsFromOpts, err := LVsOptionsList(opts).AsArgs()
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

	lvs := res.Report[0].LV

	if len(lvs) == 0 {
		return nil, nil
	}

	return lvs, nil
}

func (opts *LVsOptions) ApplyToArgs(args Arguments) error {
	for _, arg := range []Argument{
		opts.VolumeGroupName,
		opts.Tags,
		opts.CommonOptions,
		opts.ColumnOptions,
	} {
		if err := arg.ApplyToArgs(args); err != nil {
			return err
		}
	}
	return nil
}

func (list LVsOptionsList) AsArgs() (Arguments, error) {
	args := NewArgs(ArgsTypeLVs)
	options := LVsOptions{}
	for _, opt := range list {
		opt.ApplyToLVsOptions(&options)
	}
	if err := options.ApplyToArgs(args); err != nil {
		return nil, err
	}
	return args, nil
}

func (opts *LVsOptions) ApplyToLVsOptions(new *LVsOptions) {
	*new = *opts
}
