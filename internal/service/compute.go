package service

import (
	"context"

	app "girder/internal/app"
	"girder/internal/domain/entity"
)

type ComputeServiceImpl struct {
	VMDriver       any
	ImageDriver    any
	VolumeDriver   any
	SnapshotDriver any
	NodeDriver     any
}

func (ComputeServiceImpl) CreateVM(context.Context, app.CreateVMRequest) (entity.VM, error) {
	return entity.VM{}, nil
}

func (ComputeServiceImpl) DeleteVM(context.Context, app.DeleteVMRequest) (entity.VMID, error) {
	return "", nil
}

func (ComputeServiceImpl) StartVM(context.Context, app.StartVMRequest) (entity.VM, error) {
	return entity.VM{}, nil
}

func (ComputeServiceImpl) StopVM(context.Context, app.StopVMRequest) (entity.VM, error) {
	return entity.VM{}, nil
}

func (ComputeServiceImpl) RestartVM(context.Context, app.RestartVMRequest) (entity.VM, error) {
	return entity.VM{}, nil
}

func (ComputeServiceImpl) MigrateVM(context.Context, app.MigrateVMRequest) (entity.VM, error) {
	return entity.VM{}, nil
}

func (ComputeServiceImpl) ImportImage(context.Context, app.ImportImageRequest) (entity.Image, error) {
	return entity.Image{}, nil
}

func (ComputeServiceImpl) DeleteImage(context.Context, app.DeleteImageRequest) (entity.ImageID, error) {
	return "", nil
}

func (ComputeServiceImpl) CreateVolume(context.Context, app.CreateVolumeRequest) (entity.Volume, error) {
	return entity.Volume{}, nil
}

func (ComputeServiceImpl) AttachVolume(context.Context, app.AttachVolumeRequest) (entity.Volume, error) {
	return entity.Volume{}, nil
}

func (ComputeServiceImpl) DetachVolume(context.Context, app.DetachVolumeRequest) (entity.Volume, error) {
	return entity.Volume{}, nil
}

func (ComputeServiceImpl) CreateSnapshot(context.Context, app.CreateSnapshotRequest) (entity.Snapshot, error) {
	return entity.Snapshot{}, nil
}

func (ComputeServiceImpl) RestoreSnapshot(context.Context, app.RestoreSnapshotRequest) (entity.Volume, error) {
	return entity.Volume{}, nil
}
