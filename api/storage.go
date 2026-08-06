package api

import "context"

type VolumeState string

const (
	VolumeStateAvailable VolumeState = "available"
	VolumeStateInUse     VolumeState = "in-use"
	VolumeStateArchived  VolumeState = "archived"
)

type Volume struct {
	ID         string
	Name       string
	Driver     string
	ExternalID string
	SizeGiB    int64
	State      VolumeState
	AttachedTo string
	Metadata   Metadata
}

type VolumeSnapshot struct {
	ID          string
	VolumeID    string
	Name        string
	Description string
	ExternalID  string
	Metadata    Metadata
}

type ListVolumesRequest struct {
	DriverID string
	All      bool
}

type CreateVolumeRequest struct {
	DriverID string
	Volume   Volume
}

type DeleteVolumeRequest struct {
	DriverID string
	VolumeID string
}

type AttachVolumeRequest struct {
	DriverID string
	VolumeID string
	TargetID string
}

type DetachVolumeRequest struct {
	DriverID string
	VolumeID string
	TargetID string
}

type ResizeVolumeRequest struct {
	DriverID string
	VolumeID string
	SizeGiB  int64
}

type CreateVolumeSnapshotRequest struct {
	DriverID string
	VolumeID string
	Snapshot VolumeSnapshot
}

type DeleteVolumeSnapshotRequest struct {
	DriverID   string
	VolumeID   string
	SnapshotID string
}

type ListVolumeSnapshotsRequest struct {
	DriverID string
	VolumeID string
	All      bool
}

type StorageAPI interface {
	ListVolumes(context.Context, ListVolumesRequest) ([]Volume, error)
	CreateVolume(context.Context, CreateVolumeRequest) (Volume, error)
	DeleteVolume(context.Context, DeleteVolumeRequest) (string, error)
	AttachVolume(context.Context, AttachVolumeRequest) (Volume, error)
	DetachVolume(context.Context, DetachVolumeRequest) (Volume, error)
	ResizeVolume(context.Context, ResizeVolumeRequest) (Volume, error)
	CreateSnapshot(context.Context, CreateVolumeSnapshotRequest) (VolumeSnapshot, error)
	DeleteSnapshot(context.Context, DeleteVolumeSnapshotRequest) (string, error)
	ListSnapshots(context.Context, ListVolumeSnapshotsRequest) ([]VolumeSnapshot, error)
}
