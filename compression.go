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

type Compression bool

func (opt *Compression) ApplyToArgs(args Arguments) error {
	if opt == nil {
		return nil
	}
	args.AddOrReplaceAll([]string{"--compression", map[bool]string{true: "y", false: "n"}[bool(*opt)]})
	return nil
}

func (opt *Compression) ApplyToLVChangeOptions(opts *LVChangeOptions) {
	opts.Compression = opt
}
