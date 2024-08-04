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
	PVMoveOptions struct {
		From PhysicalVolumeName
		To   PhysicalVolumeNames
		LogicalVolumeName
		AllocationPolicy
		CommonOptions
	}
	PVMoveOption interface {
		ApplyToPVMoveOptions(opts *PVMoveOptions)
	}
	PVMoveOptionsList []PVMoveOption
)

func (opts *PVMoveOptions) SetOldOrNew(name PhysicalVolumeName) {
	if opts.From == "" {
		opts.From = name
	} else {
		opts.To = append(opts.To, name)
	}
}

var (
	_ ArgumentGenerator = PVMoveOptionsList{}
	_ Argument          = (*PVMoveOptions)(nil)
)

func (c *client) PVMove(ctx context.Context, opts ...PVMoveOption) error {
	args, err := PVMoveOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunLVM(ctx, append([]string{"pvmove"}, args.GetRaw()...)...)
}

func (opts *PVMoveOptions) ApplyToPVMoveOptions(new *PVMoveOptions) {
	*new = *opts
}

func (list PVMoveOptionsList) AsArgs() (Arguments, error) {
	args := NewArgs(ArgsTypeGeneric)
	options := PVMoveOptions{}
	for _, opt := range list {
		opt.ApplyToPVMoveOptions(&options)
	}
	if err := options.ApplyToArgs(args); err != nil {
		return nil, err
	}
	return args, nil
}

func (opts *PVMoveOptions) ApplyToArgs(args Arguments) error {
	if opts.From == "" {
		return fmt.Errorf("from is empty: %w", ErrPhysicalVolumeNameRequired)
	}
	if len(opts.To) == 0 {
		return fmt.Errorf("to is empty: %w", ErrPhysicalVolumeNameRequired)
	}

	for _, arg := range []Argument{
		opts.LogicalVolumeName,
		opts.From,
		opts.To,
		opts.AllocationPolicy,
		opts.CommonOptions,
	} {
		if err := arg.ApplyToArgs(args); err != nil {
			return err
		}
	}

	return nil
}
