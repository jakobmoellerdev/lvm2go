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
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
)

type DeviceList []DeviceListEntry

type DeviceIDType string

const (
	DeviceIDTypeSysWWID   DeviceIDType = "sys_wwid"
	DeviceIDTypeWWIDNAA   DeviceIDType = "wwid_naa"
	DeviceIDTypeWWIDT10   DeviceIDType = "wwid_t10"
	DeviceIDTypeWWIDEUI   DeviceIDType = "wwid_eui"
	DeviceIDTypeSysSerial DeviceIDType = "sys_serial"
	DeviceIDTypeMPathUUID DeviceIDType = "mpath_uuid"
	DeviceIDTypeCryptUUID DeviceIDType = "crypt_uuid"
	DeviceIDTypeMDUUID    DeviceIDType = "md_uuid"
	DeviceIDTypeLVMLVUUID DeviceIDType = "lvmlv_uuid"
	DeviceIDTypeLoopFile  DeviceIDType = "loop_file"
	DeviceIDTypeDevname   DeviceIDType = "devname"
)

func (opt DeviceIDType) ApplyToDevModifyOptions(opts *DevModifyOptions) {
	opts.DeviceIDType = opt
}

func (opt DeviceIDType) ApplyToArgs(args Arguments) error {
	if opt == "" {
		return nil
	}
	args.AddOrReplaceAll([]string{"--deviceidtype", string(opt)})
	return nil
}

type DeviceListEntry struct {
	IDType  DeviceIDType `json:"id_type"`
	IDName  string       `json:"id_name"`
	DevName string       `json:"dev_name"`
	PVID    string       `json:"pvid"`
}

type (
	DevListOptions struct {
		DevicesFile
	}
	DevListOption interface {
		ApplyToDevListOptions(opts *DevListOptions)
	}
	DevListOptionsList []DevListOption
)

var (
	_ ArgumentGenerator = DevListOptionsList{}
	_ Argument          = (*DevListOptions)(nil)
)

func (c *client) DevList(ctx context.Context, opts ...DevListOption) ([]DeviceListEntry, error) {
	args, err := DevListOptionsList(opts).AsArgs()
	if err != nil {
		return nil, err
	}

	var devList []DeviceListEntry
	devListProcessor := RawOutputProcessor(func(line io.Reader) error {
		scanner := bufio.NewScanner(line)
		for scanner.Scan() {
			fields := strings.Fields(strings.TrimSpace(scanner.Text()))
			if fields[0] != "Device" {
				return fmt.Errorf("invalid device list header: %q", scanner.Text())
			}
			dev := fields[1]
			kvs := make(map[string]string, len(fields)-2)
			for _, field := range fields[2:] {
				kv := strings.Split(field, "=")
				if len(kv) != 2 {
					return fmt.Errorf("invalid device list field: %q", field)
				}
				kvs[kv[0]] = kv[1]
			}
			entry := DeviceListEntry{}
			for k, v := range kvs {
				switch k {
				case "IDTYPE":
					entry.IDType = DeviceIDType(v)
				case "IDNAME":
					entry.IDName = v
				case "DEVNAME":
					if dev != v {
						return fmt.Errorf("invalid device list entry: %q", scanner.Text())
					}
					entry.DevName = v
				case "PVID":
					entry.PVID = v
				}
			}
			devList = append(devList, entry)
		}
		return scanner.Err()
	})

	if err := c.RunLVMRaw(ctx, devListProcessor, append([]string{"lvmdevices"}, args.GetRaw()...)...); err != nil {
		return nil, fmt.Errorf("failed to get version: %v", err)
	}

	return devList, nil
}

func (list DevListOptionsList) AsArgs() (Arguments, error) {
	args := NewArgs(ArgsTypeGeneric)
	options := DevListOptions{}
	for _, opt := range list {
		opt.ApplyToDevListOptions(&options)
	}
	if err := options.ApplyToArgs(args); err != nil {
		return nil, err
	}
	return args, nil
}

func (opts *DevListOptions) ApplyToDevListOptions(new *DevListOptions) {
	*new = *opts
}

func (opts *DevListOptions) ApplyToArgs(args Arguments) error {
	if err := opts.DevicesFile.ApplyToArgs(args); err != nil {
		return err
	}
	return nil
}
