package entity

import (
	"fmt"
	"strings"
)

type VolumeState string

const (
	VolumeStateAvailable VolumeState = "available"
	VolumeStateInUse     VolumeState = "in-use"
	VolumeStateArchived  VolumeState = "archived"
)

type Volume struct {
	ID          VolumeID
	Name        string
	Description string
	SizeGiB     int64
	AttachedVMID *VMID
	State       VolumeState
	Encrypted   bool
	SnapshotIDs map[SnapshotID]struct{}
}

func NewVolume(name string, sizeGiB int64) (Volume, error) {
	if strings.TrimSpace(name) == "" {
		return Volume{}, fmt.Errorf("%w: volume name", ErrEmptyName)
	}
	if sizeGiB < 0 {
		return Volume{}, fmt.Errorf("%w: volume size", ErrInvalidResource)
	}
	return Volume{
		ID:          NewVolumeID(),
		Name:        strings.TrimSpace(name),
		SizeGiB:     sizeGiB,
		State:       VolumeStateAvailable,
		SnapshotIDs: map[SnapshotID]struct{}{},
	}, nil
}

func (v *Volume) Rename(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("%w: volume name", ErrEmptyName)
	}
	v.Name = strings.TrimSpace(name)
	return nil
}

func (v *Volume) SetDescription(description string) {
	v.Description = strings.TrimSpace(description)
}

func (v *Volume) Resize(sizeGiB int64) error {
	if sizeGiB < 0 {
		return fmt.Errorf("%w: volume size", ErrInvalidResource)
	}
	v.SizeGiB = sizeGiB
	return nil
}

func (v *Volume) AttachVM(vmID VMID) error {
	if vmID == "" {
		return fmt.Errorf("%w: volume vm", ErrInvalidRelation)
	}
	v.AttachedVMID = &vmID
	v.State = VolumeStateInUse
	return nil
}

func (v *Volume) DetachVM() {
	v.AttachedVMID = nil
	v.State = VolumeStateAvailable
}

func (v *Volume) EnableEncryption() {
	v.Encrypted = true
}

func (v *Volume) DisableEncryption() {
	v.Encrypted = false
}

func (v *Volume) AddSnapshot(snapshotID SnapshotID) error {
	if snapshotID == "" {
		return fmt.Errorf("%w: snapshot id", ErrInvalidRelation)
	}
	if v.SnapshotIDs == nil {
		v.SnapshotIDs = map[SnapshotID]struct{}{}
	}
	v.SnapshotIDs[snapshotID] = struct{}{}
	return nil
}

func (v *Volume) RemoveSnapshot(snapshotID SnapshotID) {
	if v.SnapshotIDs == nil {
		return
	}
	delete(v.SnapshotIDs, snapshotID)
}

func (v *Volume) Archive() {
	v.State = VolumeStateArchived
}