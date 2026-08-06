package app

import (
	"context"

	"girder/internal/domain/entity"
)

type ComputeService interface {
	CreateVM(context.Context, CreateVMRequest) (entity.VM, error)
	DeleteVM(context.Context, DeleteVMRequest) (entity.VMID, error)
	StartVM(context.Context, StartVMRequest) (entity.VM, error)
	StopVM(context.Context, StopVMRequest) (entity.VM, error)
	RestartVM(context.Context, RestartVMRequest) (entity.VM, error)
	MigrateVM(context.Context, MigrateVMRequest) (entity.VM, error)
	ImportImage(context.Context, ImportImageRequest) (entity.Image, error)
	DeleteImage(context.Context, DeleteImageRequest) (entity.ImageID, error)
	CreateVolume(context.Context, CreateVolumeRequest) (entity.Volume, error)
	AttachVolume(context.Context, AttachVolumeRequest) (entity.Volume, error)
	DetachVolume(context.Context, DetachVolumeRequest) (entity.Volume, error)
	CreateSnapshot(context.Context, CreateSnapshotRequest) (entity.Snapshot, error)
	RestoreSnapshot(context.Context, RestoreSnapshotRequest) (entity.Volume, error)
}
