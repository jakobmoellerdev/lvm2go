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

type MetaClient interface {
	Version(ctx context.Context, opts ...VersionOption) (Version, error)
	RawConfig(ctx context.Context, opts ...ConfigOption) (RawConfig, error)
}

type VolumeGroupClient interface {
	VGs(ctx context.Context, opts ...VGsOption) ([]VolumeGroup, error)
	VGCreate(ctx context.Context, opts ...VGCreateOption) error
	VGRemove(ctx context.Context, opts ...VGRemoveOption) error
	VGExtend(ctx context.Context, opts ...VGExtendOption) error
	VGReduce(ctx context.Context, opts ...VGReduceOption) error
	VGRename(ctx context.Context, opts ...VGRenameOption) error
	VGChange(ctx context.Context, opts ...VGChangeOption) error
}

type LogicalVolumeClient interface {
	LVs(ctx context.Context, opts ...LVsOption) ([]LogicalVolume, error)
	LVCreate(ctx context.Context, opts ...LVCreateOption) error
	LVRemove(ctx context.Context, opts ...LVRemoveOption) error
	LVResize(ctx context.Context, opts ...LVResizeOption) error
	LVExtend(ctx context.Context, opts ...LVExtendOption) error
	LVReduce(ctx context.Context, opts ...LVReduceOption) error
	LVRename(ctx context.Context, opts ...LVRenameOption) error
	LVChange(ctx context.Context, opts ...LVChangeOption) error
}

type PhysicalVolumeClient interface {
	PVs(ctx context.Context, opts ...PVsOption) ([]PhysicalVolume, error)
	PVCreate(ctx context.Context, opts ...PVCreateOption) error
	PVRemove(ctx context.Context, opts ...PVRemoveOption) error
	PVResize(ctx context.Context, opts ...PVResizeOption) error
	PVChange(ctx context.Context, opts ...PVChangeOption) error
}

type Client interface {
	LogicalVolumeClient
	VolumeGroupClient
	PhysicalVolumeClient
	MetaClient
}
