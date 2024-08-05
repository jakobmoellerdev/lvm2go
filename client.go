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
	"io"
)

var (
	ErrVolumeGroupNotFound   = errors.New("volume group not found")
	ErrLogicalVolumeNotFound = errors.New("logical volume not found")
)

type client struct{}

var _ Client = (*client)(nil)

func NewClient() Client {
	return &client{}
}

// Client provides operations on lvm2 logical volumes, volume groups, and physical volumes as well as the hosts lvm2
// subsystem.
type Client interface {
	LogicalVolumeClient
	VolumeGroupClient
	PhysicalVolumeClient
	DevicesClient
	MetaClient
	DevicesClient
}

// MetaClient is a client that provides metadata information about the LVM2 library.
// This includes the version of the library and the raw configuration on the host system.
type MetaClient interface {
	// Version returns the version of the LVM2 library.
	// If the version cannot be determined, an error is returned.
	//
	// See man lvm version for more information.
	Version(ctx context.Context, opts ...VersionOption) (Version, error)
	// RawConfig returns the raw configuration of the LVM2 library.
	// If the configuration cannot be determined, an error is returned.
	// The configuration is returned as an unstructured map.
	//
	// See man lvm config for more information.
	RawConfig(ctx context.Context, opts ...ConfigOption) (RawConfig, error)

	// ReadAndDecodeConfig requests and decodes configuration values from lvm2 formatted files.
	// The configuration values are decoded into the given value v.
	// If the configuration cannot be determined, an error is returned.
	// The configuration v has to be formatted in a very specific way.
	//
	// An inner struct field represents each Config Block in an lvm configuration file.
	// The inner struct field must be tagged with the tag of key LVMConfigStructTag set to the value
	// corresponding to the config block.
	// The config block struct field must be exported.
	//
	// A field represents each Value in a Config Block in the inner struct.
	// The field must be tagged with the LVMConfigStructTag and its value set to the key of the value in the lvm config.
	// The field must be exported.
	//
	// Example:
	// type LVMConfig struct {
	//     Devices struct {
	//		 Dir string `lvm:"dir"`
	//     } `lvm:"devices"`
	// }
	//
	// The above struct will result in the following query:
	// "lvm config devices/dir"
	//
	// Note that the query can be extended and changed similarly to RawConfig.
	// E.g., to query the full merged configuration, use ConfigTypeFull.
	// Otherwise, the default configuration is queried, which might not result in a key being found.
	//
	// Possible value types for configuration keys are:
	// - string
	// - int(8,16,32,64)
	ReadAndDecodeConfig(ctx context.Context, v any, opts ...ConfigOption) error

	// WriteAndEncodeConfig writes configuration values to the given writer.
	// The configuration values are encoded from the given value v.
	// If the configuration cannot be written, an error is returned.
	// The configuration v has to be formatted in a very specific way that is equivalent to the format of ReadAndDecodeConfig.
	//
	// Example:
	// type LVMConfig struct {
	//     Devices struct {
	//		 Dir string `lvm:"dir"`
	//     } `lvm:"devices"`
	// }
	//
	// The above struct will result in the following write
	// devices {
	//     dir="/dev"
	// }
	//
	// Possible value types for configuration keys are:
	// - string
	// - int64
	//
	// Note that in lvm2, int64 is used for all integer values as well as for boolean values (0 = false, 1 = true).
	WriteAndEncodeConfig(ctx context.Context, v any, writer io.Writer) error

	// UpdateGlobalConfig updates the global configuration with the given values from v.
	// If the configuration cannot be updated, an error is returned.
	// For more information on the written config file, see LVMGlobalConfiguration
	// For more information on v and its structure, see WriteAndEncodeConfig.
	UpdateGlobalConfig(ctx context.Context, v any) error

	// UpdateLocalConfig updates the local configuration with the given values from v.
	// If the configuration cannot be updated, an error is returned.
	// For more information on the written config file, see LVMLocalConfiguration
	// For more information on v and its structure, see WriteAndEncodeConfig.
	UpdateLocalConfig(ctx context.Context, v any) error

	// UpdateProfileConfig updates the profile configuration with the given values from v.
	// The profile is expected to be resolvable to a valid path (for more information see GetProfilePath).
	// If the configuration cannot be updated, an error is returned.
	// For more information on v and its structure, see WriteAndEncodeConfig.
	UpdateProfileConfig(ctx context.Context, v any, profile Profile) error

	// CreateProfile creates a profile with the given profileName and value.
	// The Profile is encoded from the given value v.
	// The Profile is expected to be resolvable to a valid path (for more information see GetProfilePath).
	//
	// Note that although all keys can be used in the profile, lvm2 might error on unknown keys or fail
	// on unsupported keys.
	// To avoid this, make sure the keys in v are one of the keys reported by lvm2
	// when running "lvmconfig --typeconfig profilable" (or use ConfigTypeProfilable with RawConfig).
	CreateProfile(ctx context.Context, v any, profile Profile) (string, error)

	// RemoveProfile removes a profile with the given profileName.
	// The Profile is expected to be resolvable to a valid path (for more information see GetProfilePath).
	RemoveProfile(ctx context.Context, profile Profile) error

	// GetProfilePath returns the path to the profile within the profile directory as configured on the host.
	//
	// Example:
	// for a configured profile directory /etc/lvm/profile on the host, the following result will be returned:
	// - GetProfilePath(ctx, "test") -> "/etc/lvm/profile/test.profile", nil
	// - GetProfilePath(ctx, "test.profile") -> "/etc/lvm/profile/test.profile", nil
	// - GetProfilePath(ctx, "/etc/lvm/profile/test") -> "/etc/lvm/profile/test.profile", nil
	// - GetProfilePath(ctx, "/etc/lvm/profile/test.profile") -> "/etc/lvm/profile/test.profile", nil
	// - GetProfilePath(ctx, "/var/test") -> "", error
	// - GetProfilePath(ctx, "/var/") -> "", error
	// - GetProfilePath(ctx, "test.something") -> "", error
	//
	// For more information on the profile directory, check the lvm2 configuration.
	// Usually, the directory is set to /etc/lvm/profile as per config key config/profile_dir,
	// but it can be changed to any other directory based on the host.
	// For getting the current profile directory, see GetProfileDirectory.
	GetProfilePath(ctx context.Context, profile Profile) (string, error)

	// GetProfileDirectory returns the profile directory as configured on the host.
	// If the profile directory cannot be determined, an error is returned.
	// The profile directory is the directory where lvm2 profiles are stored.
	// The directory is expected to be resolvable to a valid path.
	//
	// See man lvm and man lvmconfig for more information.
	GetProfileDirectory(ctx context.Context) (string, error)
}

// VolumeGroupClient is a client that provides operations on lvm2 volume groups.
type VolumeGroupClient interface {
	// VG returns a volume group that matches the given options.
	//
	// If no VolumeGroupName is defined, ErrVolumeGroupNameRequired is returned.
	// If no volume group is found, ErrVolumeGroupNotFound is returned.
	//
	// It is equivalent to calling VGs with the same options and returning the first volume group in the list.
	// see VGs for more information.
	VG(ctx context.Context, opts ...VGsOption) (*VolumeGroup, error)

	// VGs return a list of volume groups that match the given options.
	//
	// If no volume groups are found, an empty slice is returned.
	// If options limit the number of volume groups returned,
	// the slice may be shorter than the total number of volume groups.
	//
	// See man lvm vgs for more information.
	VGs(ctx context.Context, opts ...VGsOption) ([]*VolumeGroup, error)

	// VGCreate creates a new volume group with the given options.
	//
	// See man lvm vgcreate for more information.
	VGCreate(ctx context.Context, opts ...VGCreateOption) error

	// VGRemove removes a volume group with the given options.
	//
	// See man lvm vgremove for more information.
	VGRemove(ctx context.Context, opts ...VGRemoveOption) error

	// VGExtend extends a volume group with the given options.
	//
	// See man lvm vgextend for more information.
	VGExtend(ctx context.Context, opts ...VGExtendOption) error

	// VGReduce reduces a volume group with the given options.
	//
	// See man lvm vgreduce for more information.
	VGReduce(ctx context.Context, opts ...VGReduceOption) error

	// VGRename renames a volume group with the given options.
	//
	// See man lvm vgrename for more information.
	VGRename(ctx context.Context, opts ...VGRenameOption) error

	// VGChange changes a volume group with the given options.
	//
	// See man lvm vgchange for more information.
	VGChange(ctx context.Context, opts ...VGChangeOption) error
}

// LogicalVolumeClient is a client that provides operations on lvm2 logical volumes.
type LogicalVolumeClient interface {
	// LV returns a logical volume that matches the given options.
	//
	// If no LogicalVolumeName is defined, ErrLogicalVolumeNameRequired is returned.
	// If no VolumeGroupName is defined, ErrVolumeGroupNameRequired is returned.
	// If no logical volume is found in the volume group, ErrLogicalVolumeNotFound is returned.
	//
	// It is equivalent to calling LVs with the same options and returning the first logical volume in the list.
	// see LVs for more information.
	LV(ctx context.Context, opts ...LVsOption) (*LogicalVolume, error)

	// LVs return a list of logical volumes that match the given options.
	//
	// If no logical volumes are found, an empty slice is returned.
	// If options limit the number of volume groups returned,
	// the slice may be shorter than the total number of logical volumes.
	//
	// See man lvm lvs for more information.
	LVs(ctx context.Context, opts ...LVsOption) ([]*LogicalVolume, error)

	// LVCreate creates a new logical volume with the given options.
	//
	// See man lvm lvcreate for more information.
	LVCreate(ctx context.Context, opts ...LVCreateOption) error

	// LVRemove removes a logical volume with the given options.
	//
	// See man lvm lvremove for more information.
	LVRemove(ctx context.Context, opts ...LVRemoveOption) error

	// LVResize resizes a logical volume with the given options.
	//
	// See man lvm lvresize for more information.
	LVResize(ctx context.Context, opts ...LVResizeOption) error

	// LVExtend extends a logical volume with the given options.
	//
	// See man lvm lvextend for more information.
	LVExtend(ctx context.Context, opts ...LVExtendOption) error

	// LVReduce reduces a logical volume with the given options.
	//
	// See man lvm lvreduce for more information.
	LVReduce(ctx context.Context, opts ...LVReduceOption) error

	// LVRename renames a logical volume with the given options.
	//
	// See man lvm lvrename for more information.
	LVRename(ctx context.Context, opts ...LVRenameOption) error

	// LVChange changes a logical volume with the given options.
	//
	// See man lvm lvchange for more information.
	LVChange(ctx context.Context, opts ...LVChangeOption) error
}

// PhysicalVolumeClient is a client that provides operations on lvm2 physical volumes.
type PhysicalVolumeClient interface {
	// PVs return a list of physical volumes that match the given options.
	//
	// If no physical volumes are found, an empty slice is returned.
	// If options limit the number of physical volumes returned,
	// the slice may be shorter than the total number of physical volumes.
	//
	// See man lvm pvs for more information.
	PVs(ctx context.Context, opts ...PVsOption) ([]*PhysicalVolume, error)

	// PVCreate creates a new physical volume with the given options.
	//
	// See man lvm pvcreate for more information.
	PVCreate(ctx context.Context, opts ...PVCreateOption) error

	// PVRemove removes a physical volume with the given options.
	//
	// See man lvm pvremove for more information.
	PVRemove(ctx context.Context, opts ...PVRemoveOption) error

	// PVResize removes a physical volume with the given options.
	//
	// See man lvm pvresize for more information.
	PVResize(ctx context.Context, opts ...PVResizeOption) error

	// PVChange changes a physical volume with the given options.
	//
	// See man lvm pvchange for more information.
	PVChange(ctx context.Context, opts ...PVChangeOption) error

	// PVMove moves extents between physical volumes with the given options.
	//
	// see man lvm pvmove for more information.
	PVMove(ctx context.Context, opts ...PVMoveOption) error
}

// DevicesClient is a client that provides operations on lvm2 device files.
type DevicesClient interface {
	// DevList returns a list of devices that match the given options.
	//
	// Replicates lvmdevices
	// See man lvmdevices for more information.
	DevList(ctx context.Context, opts ...DevListOption) ([]DeviceListEntry, error)

	// DevCheck checks the device files and returns an error if any inconsistencies are found.
	//
	// Replicates lvmdevices --check
	// See man lvmdevices for more information.
	DevCheck(ctx context.Context, opts ...DevCheckOption) error

	// DevUpdate updates the device files through attempted automatic corrections.
	//
	// Replicates lvmdevices --update
	// See man lvmdevices for more information.
	DevUpdate(ctx context.Context, opts ...DevUpdateOption) error

	// DevModify adds and removes devices in device files with the given options.
	//
	// Replicates lvmdevices --adddev, --addpvid, --deldev and --delpvid
	// See man lvmdevices for more information.
	DevModify(ctx context.Context, opts ...DevModifyOption) error
}
