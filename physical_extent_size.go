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

type PhysicalExtentSize Size

func (opt PhysicalExtentSize) ApplyToVGCreateOptions(opts *VGCreateOptions) {
	opts.PhysicalExtentSize = opt
}

func (opt PhysicalExtentSize) ApplyToArgs(args Arguments) error {
	if opt.Val == 0 {
		return nil
	}

	size := Size(opt)

	if err := size.Validate(); err != nil {
		return err
	}

	args.AddOrReplaceAll([]string{"--physicalextentsize", size.String()})
	return nil
}
