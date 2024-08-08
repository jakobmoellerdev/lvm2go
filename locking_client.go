package lvm2go

import (
	"context"
	"io"
	"sync"
)

type lockingClient struct {
	clnt Client
	mu   sync.RWMutex
}

// NewLockingClient returns a new Client that locks all methods with a read-write mutex.
// This is useful when you want to ensure that only one operation is happening at a time.
// This can however only work if all operations are done through the same client.
// It is a helper for synchronizing dangerous concurrent calls to the same client.
// Note that this can introduce significant performance overhead if the client is used in a highly concurrent environment.
func NewLockingClient(clnt Client) Client {
	lc := &lockingClient{clnt: clnt}
	return lc
}

var _ Client = &lockingClient{}

func (l *lockingClient) LV(ctx context.Context, opts ...LVsOption) (*LogicalVolume, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.clnt.LV(ctx, opts...)
}

func (l *lockingClient) LVs(ctx context.Context, opts ...LVsOption) ([]*LogicalVolume, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.clnt.LVs(ctx, opts...)
}

func (l *lockingClient) LVCreate(ctx context.Context, opts ...LVCreateOption) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.clnt.LVCreate(ctx, opts...)
}

func (l *lockingClient) LVRemove(ctx context.Context, opts ...LVRemoveOption) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.clnt.LVRemove(ctx, opts...)
}

func (l *lockingClient) LVResize(ctx context.Context, opts ...LVResizeOption) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.clnt.LVResize(ctx, opts...)
}

func (l *lockingClient) LVExtend(ctx context.Context, opts ...LVExtendOption) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.clnt.LVExtend(ctx, opts...)
}

func (l *lockingClient) LVReduce(ctx context.Context, opts ...LVReduceOption) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.clnt.LVReduce(ctx, opts...)
}

func (l *lockingClient) LVRename(ctx context.Context, opts ...LVRenameOption) error {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.clnt.LVRename(ctx, opts...)
}

func (l *lockingClient) LVChange(ctx context.Context, opts ...LVChangeOption) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.clnt.LVChange(ctx, opts...)
}

func (l *lockingClient) VG(ctx context.Context, opts ...VGsOption) (*VolumeGroup, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.clnt.VG(ctx, opts...)
}

func (l *lockingClient) VGs(ctx context.Context, opts ...VGsOption) ([]*VolumeGroup, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.clnt.VGs(ctx, opts...)
}

func (l *lockingClient) VGCreate(ctx context.Context, opts ...VGCreateOption) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.clnt.VGCreate(ctx, opts...)
}

func (l *lockingClient) VGRemove(ctx context.Context, opts ...VGRemoveOption) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.clnt.VGRemove(ctx, opts...)
}

func (l *lockingClient) VGExtend(ctx context.Context, opts ...VGExtendOption) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.clnt.VGExtend(ctx, opts...)
}

func (l *lockingClient) VGReduce(ctx context.Context, opts ...VGReduceOption) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.clnt.VGReduce(ctx, opts...)
}

func (l *lockingClient) VGRename(ctx context.Context, opts ...VGRenameOption) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.clnt.VGRename(ctx, opts...)
}

func (l *lockingClient) VGChange(ctx context.Context, opts ...VGChangeOption) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.clnt.VGChange(ctx, opts...)
}

func (l *lockingClient) PVs(ctx context.Context, opts ...PVsOption) ([]*PhysicalVolume, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.clnt.PVs(ctx, opts...)
}

func (l *lockingClient) PVCreate(ctx context.Context, opts ...PVCreateOption) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.clnt.PVCreate(ctx, opts...)
}

func (l *lockingClient) PVRemove(ctx context.Context, opts ...PVRemoveOption) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.clnt.PVRemove(ctx, opts...)
}

func (l *lockingClient) PVResize(ctx context.Context, opts ...PVResizeOption) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.clnt.PVResize(ctx, opts...)
}

func (l *lockingClient) PVChange(ctx context.Context, opts ...PVChangeOption) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.clnt.PVChange(ctx, opts...)
}

func (l *lockingClient) PVMove(ctx context.Context, opts ...PVMoveOption) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.clnt.PVMove(ctx, opts...)
}

func (l *lockingClient) DevList(ctx context.Context, opts ...DevListOption) ([]DeviceListEntry, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.clnt.DevList(ctx, opts...)
}

func (l *lockingClient) DevCheck(ctx context.Context, opts ...DevCheckOption) error {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.clnt.DevCheck(ctx, opts...)
}

func (l *lockingClient) DevUpdate(ctx context.Context, opts ...DevUpdateOption) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.clnt.DevUpdate(ctx, opts...)
}

func (l *lockingClient) DevModify(ctx context.Context, opts ...DevModifyOption) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.clnt.DevModify(ctx, opts...)
}

func (l *lockingClient) Version(ctx context.Context, opts ...VersionOption) (Version, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.clnt.Version(ctx, opts...)
}

func (l *lockingClient) RawConfig(ctx context.Context, opts ...ConfigOption) (RawConfig, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.clnt.RawConfig(ctx, opts...)
}

func (l *lockingClient) ReadAndDecodeConfig(ctx context.Context, v any, opts ...ConfigOption) error {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.clnt.ReadAndDecodeConfig(ctx, v, opts...)
}

func (l *lockingClient) WriteAndEncodeConfig(ctx context.Context, v any, writer io.Writer) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.clnt.WriteAndEncodeConfig(ctx, v, writer)
}

func (l *lockingClient) UpdateGlobalConfig(ctx context.Context, v any) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.clnt.UpdateGlobalConfig(ctx, v)
}

func (l *lockingClient) UpdateLocalConfig(ctx context.Context, v any) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.clnt.UpdateLocalConfig(ctx, v)
}

func (l *lockingClient) UpdateProfileConfig(ctx context.Context, v any, profile Profile) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.clnt.UpdateProfileConfig(ctx, v, profile)
}

func (l *lockingClient) CreateProfile(ctx context.Context, v any, profile Profile) (string, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.clnt.CreateProfile(ctx, v, profile)
}

func (l *lockingClient) RemoveProfile(ctx context.Context, profile Profile) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.clnt.RemoveProfile(ctx, profile)
}

func (l *lockingClient) GetProfilePath(ctx context.Context, profile Profile) (string, error) {
	// no locking needed
	return l.clnt.GetProfilePath(ctx, profile)
}

func (l *lockingClient) GetProfileDirectory(ctx context.Context) (string, error) {
	// no locking needed
	return l.clnt.GetProfileDirectory(ctx)
}
