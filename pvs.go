/*
 Copyright 2024 The lvm2go Authors.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package lvm2go

import (
	"context"
)

type (
	PVsOptions struct {
		Unit
		Tags
		Select

		ColumnOptions
		CommonOptions
	}
	PVsOption interface {
		ApplyToPVsOptions(opts *PVsOptions)
	}
	PVsOptionsList []PVsOption
)

var (
	_ ArgumentGenerator = PVsOptionsList{}
	_ Argument          = (*PVsOptions)(nil)
)

// PVs returns a list of logical volumes that match the given options.
// If no logical volumes are found, nil is returned.
// It is really just a wrapper around the `lvs --reportformat json` command.
func (c *client) PVs(ctx context.Context, opts ...PVsOption) ([]*PhysicalVolume, error) {
	type lvReport struct {
		Report []struct {
			PV []*PhysicalVolume `json:"pv"`
		} `json:"report"`
	}

	var res = new(lvReport)

	args := []string{
		"pvs", "--reportformat", "json",
	}
	argsFromOpts, err := PVsOptionsList(opts).AsArgs()
	if err != nil {
		return nil, err
	}

	err = c.RunLVMInto(ctx, res, append(args, argsFromOpts.GetRaw()...)...)

	if IsNotFound(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if len(res.Report) == 0 {
		return nil, nil
	}

	pvs := res.Report[0].PV

	if len(pvs) == 0 {
		return nil, nil
	}

	return pvs, nil
}

func (opts *PVsOptions) ApplyToArgs(args Arguments) error {
	for _, arg := range []Argument{
		opts.Unit,
		opts.Tags,
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

func (list PVsOptionsList) AsArgs() (Arguments, error) {
	args := NewArgs(ArgsTypePVs)
	options := PVsOptions{}
	for _, opt := range list {
		opt.ApplyToPVsOptions(&options)
	}
	if err := options.ApplyToArgs(args); err != nil {
		return nil, err
	}
	return args, nil
}

func (opts *PVsOptions) ApplyToPVsOptions(new *PVsOptions) {
	*new = *opts
}
