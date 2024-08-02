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

type ActivationState string

const (
	Activate     ActivationState = "y"
	Deactivate   ActivationState = "n"
	AutoActivate ActivationState = "ay"
)

func (opt ActivationState) ApplyToLVCreateOptions(opts *LVCreateOptions) {
	opts.ActivationState = opt
}

func (opt ActivationState) ApplyToLVChangeOptions(opts *LVChangeOptions) {
	opts.ActivationState = opt
}

func (opt ActivationState) ApplyToArgs(args Arguments) error {
	if opt == "" {
		return nil
	}
	args.AddOrReplaceAll([]string{"--activate", string(opt)})
	return nil
}
