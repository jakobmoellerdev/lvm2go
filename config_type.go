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

type ConfigType string

const (
	ConfigTypeFull ConfigType = "full"
)

func (c ConfigType) String() string {
	return string(c)
}

func (c ConfigType) AsArgs() []string {
	return []string{"--typeconfig", c.String()}
}

func (c ConfigType) ApplyToConfigOptions(opts *ConfigOptions) {
	opts.ConfigType = c
}

func (c ConfigType) ApplyToArgs(arguments Arguments) error {
	if len(c) == 0 {
		return nil
	}
	arguments.AddOrReplaceAll(c.AsArgs())
	return nil
}
