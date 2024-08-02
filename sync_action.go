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
	"fmt"
)

type SyncAction string

const (
	SyncActionCheck  SyncAction = "check"
	SyncActionRepair SyncAction = "repair"
)

func (opt SyncAction) ApplyToArgs(args Arguments) error {
	if opt == "" {
		return nil
	}
	args.AddOrReplace(fmt.Sprintf("--syncaction=%s", string(opt)))
	return nil
}

func (opt SyncAction) ApplyToLVChangeOptions(opts *LVChangeOptions) {
	opts.SyncAction = opt
}
