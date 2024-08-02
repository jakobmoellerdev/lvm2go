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
	"errors"
	"fmt"
)

type (
	LVReduceOptions struct {
		VolumeGroupName
		LogicalVolumeName
		CommonOptions
	}
	LVReduceOption interface {
		ApplyToLVReduceOptions(opts *LVReduceOptions)
	}
	LVReduceOptionsList []LVReduceOption
)

var (
	_ ArgumentGenerator = LVReduceOptionsList{}
	_ Argument          = (*LVReduceOptions)(nil)
)

func (c *client) LVReduce(ctx context.Context, opts ...LVReduceOption) error {
	args, err := LVReduceOptionsList(opts).AsArgs()
	if err != nil {
		return err
	}

	return c.RunLVM(ctx, append([]string{"lvreduce"}, args.GetRaw()...)...)
}

func (list LVReduceOptionsList) AsArgs() (Arguments, error) {
	return nil, fmt.Errorf("not implemented: %w", errors.ErrUnsupported)
}

func (opts *LVReduceOptions) ApplyToArgs(_ Arguments) error {
	return fmt.Errorf("not implemented: %w", errors.ErrUnsupported)
}
