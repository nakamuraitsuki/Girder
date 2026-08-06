package mock

import (
	"context"

	"girder/internal/app"
	"girder/internal/domain/entity"
)

type MockComputeService struct{}

func (MockComputeService) CreateVM(context.Context, app.CreateVMRequest) (entity.VM, error) {
	return entity.VM{}, nil
}

func (MockComputeService) DeleteVM(context.Context, app.DeleteVMRequest) (entity.VMID, error) {
	return "", nil
}

func (MockComputeService) StartVM(context.Context, app.StartVMRequest) (entity.VM, error) {
	return entity.VM{}, nil
}

func (MockComputeService) StopVM(context.Context, app.StopVMRequest) (entity.VM, error) {
	return entity.VM{}, nil
}

func (MockComputeService) RestartVM(context.Context, app.RestartVMRequest) (entity.VM, error) {
	return entity.VM{}, nil
}

func (MockComputeService) MigrateVM(context.Context, app.MigrateVMRequest) (entity.VM, error) {
	return entity.VM{}, nil
}

func (MockComputeService) ImportImage(context.Context, app.ImportImageRequest) (entity.Image, error) {
	return entity.Image{}, nil
}

func (MockComputeService) DeleteImage(context.Context, app.DeleteImageRequest) (entity.ImageID, error) {
	return "", nil
}

func (MockComputeService) CreateVolume(context.Context, app.CreateVolumeRequest) (entity.Volume, error) {
	return entity.Volume{}, nil
}

func (MockComputeService) AttachVolume(context.Context, app.AttachVolumeRequest) (entity.Volume, error) {
	return entity.Volume{}, nil
}

func (MockComputeService) DetachVolume(context.Context, app.DetachVolumeRequest) (entity.Volume, error) {
	return entity.Volume{}, nil
}

func (MockComputeService) CreateSnapshot(context.Context, app.CreateSnapshotRequest) (entity.Snapshot, error) {
	return entity.Snapshot{}, nil
}

func (MockComputeService) RestoreSnapshot(context.Context, app.RestoreSnapshotRequest) (entity.Volume, error) {
	return entity.Volume{}, nil
}
