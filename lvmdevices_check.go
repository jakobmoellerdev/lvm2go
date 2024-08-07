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
		NoOpRawOutputProcessor(),
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
