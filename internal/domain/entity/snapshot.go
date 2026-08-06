package entity

import (
	"fmt"
	"strings"
	"time"

)

type Snapshot struct {
	ID              SnapshotID
	VolumeID        VolumeID
	Name            string
	Description     string
	CreatedAt       time.Time
	ParentSnapshotID *SnapshotID
	SizeGiB         int64
}

func NewSnapshot(name string, volumeID VolumeID, sizeGiB int64) (Snapshot, error) {
	if strings.TrimSpace(name) == "" {
		return Snapshot{}, fmt.Errorf("%w: snapshot name", ErrEmptyName)
	}
	if volumeID == "" {
		return Snapshot{}, fmt.Errorf("%w: snapshot volume", ErrInvalidRelation)
	}
	if sizeGiB < 0 {
		return Snapshot{}, fmt.Errorf("%w: snapshot size", ErrInvalidResource)
	}
	return Snapshot{
		ID:        NewSnapshotID(),
		VolumeID:  volumeID,
		Name:      strings.TrimSpace(name),
		CreatedAt: time.Now().UTC(),
		SizeGiB:   sizeGiB,
	}, nil
}

func (s *Snapshot) Rename(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("%w: snapshot name", ErrEmptyName)
	}
	s.Name = strings.TrimSpace(name)
	return nil
}

func (s *Snapshot) SetDescription(description string) {
	s.Description = strings.TrimSpace(description)
}

func (s *Snapshot) SetParentSnapshot(snapshotID SnapshotID) error {
	if snapshotID == "" {
		return fmt.Errorf("%w: parent snapshot", ErrInvalidRelation)
	}
	s.ParentSnapshotID = &snapshotID
	return nil
}

func (s *Snapshot) ClearParentSnapshot() {
	s.ParentSnapshotID = nil
}

func (s *Snapshot) UpdateSize(sizeGiB int64) error {
	if sizeGiB < 0 {
		return fmt.Errorf("%w: snapshot size", ErrInvalidResource)
	}
	s.SizeGiB = sizeGiB
	return nil
}
