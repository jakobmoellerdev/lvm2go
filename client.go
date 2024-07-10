package lvm2go

import (
	"context"
)

type ClientOption interface {
	ApplyToClientOptions(opts *clientOptions)
}

type clientOptions struct{}

type client struct {
	opts clientOptions
}

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
	MetaClient
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
}

// VolumeGroupClient is a client that provides operations on lvm2 volume groups.
type VolumeGroupClient interface {
	// VGs return a list of volume groups that match the given options.
	//
	// If no volume groups are found, an empty slice is returned.
	// If options limit the number of volume groups returned,
	// the slice may be shorter than the total number of volume groups.
	//
	// See man lvm vgs for more information.
	VGs(ctx context.Context, opts ...VGsOption) ([]VolumeGroup, error)

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
	// LVs return a list of logical volumes that match the given options.
	//
	// If no logical volumes are found, an empty slice is returned.
	// If options limit the number of volume groups returned,
	// the slice may be shorter than the total number of logical volumes.
	//
	// See man lvm lvs for more information.
	LVs(ctx context.Context, opts ...LVsOption) ([]LogicalVolume, error)

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
	PVs(ctx context.Context, opts ...PVsOption) ([]PhysicalVolume, error)

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
}
