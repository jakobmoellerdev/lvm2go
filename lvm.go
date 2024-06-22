package lvm2go

import (
	"context"
)

type VolumeGroupClient interface {
	VGs(ctx context.Context, opts ...VGsOption) ([]VolumeGroup, error)
	VGCreate(ctx context.Context, opts ...VGCreateOption) error
	VGRemove(ctx context.Context, opts ...VGRemoveOption) error
}

type LogicalVolumeClient interface {
	LVs(ctx context.Context, opts ...LVsOption) ([]LogicalVolume, error)
	LVCreate(ctx context.Context, opts ...LVCreateOption) error
	LVRemove(ctx context.Context, opts ...LVRemoveOption) error
}

type Client interface {
	LogicalVolumeClient
	VolumeGroupClient
}

type client struct{}

func NewClient() Client {
	return &client{}
}

func (c *client) RunLVM(ctx context.Context, args ...string) error {
	return RunLVM(ctx, args...)
}
