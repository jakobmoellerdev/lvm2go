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
	"fmt"
)

type (
	PVResizeOptions struct {
		PhysicalVolumeName
		CommonOptions
	}
	PVResizeOption interface {
		ApplyToPVResizeOptions(opts *PVResizeOptions)
	}
	PVResizeOptionsList []PVResizeOption
)

var (
	_ ArgumentGenerator = PVResizeOptionsList{}
	_ Argument          = (*PVResizeOptions)(nil)
)

func (c *client) PVResize(ctx context.Context, opts ...PVResizeOption) error {
	args, err := PVResizeOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunLVM(ctx, append([]string{"pvresize"}, args.GetRaw()...)...)
}

func (opts *PVResizeOptions) ApplyToPVResizeOptions(new *PVResizeOptions) {
	*new = *opts
}

func (list PVResizeOptionsList) AsArgs() (Arguments, error) {
	args := NewArgs(ArgsTypeGeneric)
	options := PVResizeOptions{}
	for _, opt := range list {
		opt.ApplyToPVResizeOptions(&options)
	}
	if err := options.ApplyToArgs(args); err != nil {
		return nil, err
	}
	return args, nil
}

func (opts *PVResizeOptions) ApplyToArgs(args Arguments) error {
	if opts.PhysicalVolumeName == "" {
		return fmt.Errorf("PhysicalVolumeName is required for resizing a physical volume")
	}

	for _, arg := range []Argument{
		opts.PhysicalVolumeName,
		opts.CommonOptions,
	} {
		if err := arg.ApplyToArgs(args); err != nil {
			return err
		}
	}

	return nil
}
