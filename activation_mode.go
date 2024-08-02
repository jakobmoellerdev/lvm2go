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

type ActivationMode string

const (
	ActivationModePartial  ActivationMode = "partial"
	ActivationModeDegraded ActivationMode = "degraded"
	ActivationModeComplete ActivationMode = "complete"
)

func (opt ActivationMode) ApplyToArgs(args Arguments) error {
	if opt == "" {
		return nil
	}
	args.AddOrReplace(fmt.Sprintf("--activationmode=%s", string(opt)))
	return nil
}

func (opt ActivationMode) ApplyToLVChangeOptions(opts *LVChangeOptions) {
	opts.ActivationMode = opt
}
