package lvm2go

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
)

type DeviceList []DeviceListEntry

func (opt DeviceList) ToDevices() Devices {
	var devices Devices
	for _, entry := range opt {
		devices = append(devices, entry.Device)
	}
	return devices
}

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
	Device  string       `json:"device"`
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
			if len(fields) < 5 {
				return fmt.Errorf("invalid device list line: %q", scanner.Text())
			}
			devList = append(devList, DeviceListEntry{
				Device:  fields[0],
				IDType:  DeviceIDType(fields[1]),
				IDName:  fields[2],
				DevName: fields[3],
				PVID:    fields[4],
			})
		}
		return scanner.Err()
	})

	if err := c.RunLVMRaw(ctx, devListProcessor, append([]string{"version"}, args.GetRaw()...)...); err != nil {
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
