package entity

import (
	"fmt"
	"strings"
)

type VMState string

const (
	VMStateStopped   VMState = "stopped"
	VMStateRunning   VMState = "running"
	VMStatePaused    VMState = "paused"
	VMStateSuspended VMState = "suspended"
)

type VM struct {
	ID          VMID
	Name        string
	Description string
	ImageID     ImageID
	NodeID      *NodeID
	State       VMState
	CPU         int
	MemoryMiB   int64
	DiskGiB     int64
	PortIDs     map[PortID]struct{}
	VolumeIDs   map[VolumeID]struct{}
	Metadata    map[string]string
}

func NewVM(name string, imageID ImageID, cpu int, memoryMiB int64, diskGiB int64) (VM, error) {
	if strings.TrimSpace(name) == "" {
		return VM{}, fmt.Errorf("%w: vm name", ErrEmptyName)
	}
	if imageID == "" {
		return VM{}, fmt.Errorf("%w: vm image", ErrInvalidRelation)
	}
	if cpu < 0 || memoryMiB < 0 || diskGiB < 0 {
		return VM{}, fmt.Errorf("%w: vm resources", ErrInvalidResource)
	}
	return VM{
		ID:        NewVMID(),
		Name:      strings.TrimSpace(name),
		ImageID:   imageID,
		State:     VMStateStopped,
		CPU:       cpu,
		MemoryMiB: memoryMiB,
		DiskGiB:   diskGiB,
		PortIDs:   map[PortID]struct{}{},
		VolumeIDs: map[VolumeID]struct{}{},
		Metadata:  map[string]string{},
	}, nil
}

func (v *VM) Rename(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("%w: vm name", ErrEmptyName)
	}
	v.Name = strings.TrimSpace(name)
	return nil
}

func (v *VM) SetDescription(description string) {
	v.Description = strings.TrimSpace(description)
}

func (v *VM) SetImage(imageID ImageID) error {
	if imageID == "" {
		return fmt.Errorf("%w: vm image", ErrInvalidRelation)
	}
	v.ImageID = imageID
	return nil
}

func (v *VM) AssignNode(nodeID NodeID) error {
	if nodeID == "" {
		return fmt.Errorf("%w: vm node", ErrInvalidRelation)
	}
	v.NodeID = &nodeID
	return nil
}

func (v *VM) DetachNode() {
	v.NodeID = nil
}

func (v *VM) AttachPort(portID PortID) error {
	if portID == "" {
		return fmt.Errorf("%w: vm port", ErrInvalidRelation)
	}
	if v.PortIDs == nil {
		v.PortIDs = map[PortID]struct{}{}
	}
	v.PortIDs[portID] = struct{}{}
	return nil
}

func (v *VM) DetachPort(portID PortID) {
	if v.PortIDs == nil {
		return
	}
	delete(v.PortIDs, portID)
}

func (v *VM) AttachVolume(volumeID VolumeID) error {
	if volumeID == "" {
		return fmt.Errorf("%w: vm volume", ErrInvalidRelation)
	}
	if v.VolumeIDs == nil {
		v.VolumeIDs = map[VolumeID]struct{}{}
	}
	v.VolumeIDs[volumeID] = struct{}{}
	return nil
}

func (v *VM) DetachVolume(volumeID VolumeID) {
	if v.VolumeIDs == nil {
		return
	}
	delete(v.VolumeIDs, volumeID)
}

func (v *VM) SetState(state VMState) {
	v.State = state
}

func (v VM) ResourceRequest() NodeAllocation {
	return NodeAllocation{CPU: v.CPU, MemoryMiB: v.MemoryMiB, DiskGiB: v.DiskGiB}
}

func (v *VM) SetResources(cpu int, memoryMiB int64, diskGiB int64) error {
	if cpu < 0 || memoryMiB < 0 || diskGiB < 0 {
		return fmt.Errorf("%w: vm resources", ErrInvalidResource)
	}
	v.CPU = cpu
	v.MemoryMiB = memoryMiB
	v.DiskGiB = diskGiB
	return nil
}

func (v *VM) SetMetadata(key, val string) error {
	key = strings.TrimSpace(key)
	if key == "" {
		return fmt.Errorf("%w: metadata key", ErrEmptyName)
	}
	if v.Metadata == nil {
		v.Metadata = map[string]string{}
	}
	v.Metadata[key] = strings.TrimSpace(val)
	return nil
}

func (v *VM) RemoveMetadata(key string) {
	if v.Metadata == nil {
		return
	}
	delete(v.Metadata, strings.TrimSpace(key))
}

func (v *VM) Start() { v.State = VMStateRunning }
func (v *VM) Stop() { v.State = VMStateStopped }
func (v *VM) Pause() { v.State = VMStatePaused }
func (v *VM) Suspend() { v.State = VMStateSuspended }