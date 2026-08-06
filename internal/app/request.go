package app

import (
	"girder/internal/domain/entity"
	"girder/internal/domain/value"
)

type RegisterNodeRequest struct {
	Name     string
	Capacity entity.NodeCapacity
	Labels   map[string]string
}

type RemoveNodeRequest struct {
	NodeID entity.NodeID
}

type GetNodeRequest struct {
	NodeID entity.NodeID
}

type CreateVMRequest struct {
	Name      string
	ImageID   entity.ImageID
	Resources entity.NodeAllocation
	NodeID    *entity.NodeID
	Metadata  map[string]string
}

type DeleteVMRequest struct {
	VMID entity.VMID
}

type StartVMRequest struct {
	VMID entity.VMID
}

type StopVMRequest struct {
	VMID entity.VMID
}

type RestartVMRequest struct {
	VMID entity.VMID
}

type MigrateVMRequest struct {
	VMID   entity.VMID
	NodeID entity.NodeID
}

type ImportImageRequest struct {
	Name         string
	OS           string
	Architecture string
	Version      string
	SizeBytes    int64
	Checksum     string
	Tags         map[string]string
}

type DeleteImageRequest struct {
	ImageID entity.ImageID
}

type CreateVolumeRequest struct {
	Name      string
	SizeGiB   int64
	Encrypted bool
}

type AttachVolumeRequest struct {
	VolumeID entity.VolumeID
	VMID     entity.VMID
}

type DetachVolumeRequest struct {
	VolumeID entity.VolumeID
}

type CreateSnapshotRequest struct {
	Name     string
	VolumeID entity.VolumeID
	SizeGiB  int64
}

type RestoreSnapshotRequest struct {
	SnapshotID entity.SnapshotID
	VolumeID   entity.VolumeID
}

type CreateSwitchRequest struct {
	Name string
}

type DeleteSwitchRequest struct {
	SwitchID entity.SwitchID
}

type CreateRouterRequest struct {
	Name string
}

type DeleteRouterRequest struct {
	RouterID entity.RouterID
}

type ConnectPortRequest struct {
	PortID   entity.PortID
	VMID     *entity.VMID
	SwitchID *entity.SwitchID
	RouterID *entity.RouterID
}

type DisconnectPortRequest struct {
	PortID entity.PortID
}

type AddRouteRequest struct {
	Route entity.Route
}

type DeleteRouteRequest struct {
	RouteID entity.RouteID
}

type ConfigureNATRequest struct {
	Name        string
	Type        entity.NATType
	Source      value.CIDR
	Translation value.IPv4
	Enabled     bool
}

type ConfigureDHCPRequest struct {
	Name         string
	Network      value.CIDR
	RangeStart   value.IPv4
	RangeEnd     value.IPv4
	Gateway      value.IPv4
	DNSServers   []value.IPv4
	LeaseSeconds int
	Enabled      bool
}

type CreateDNSZoneRequest struct {
	Name       string
	Origin     string
	TTLSeconds uint32
}

type DeleteDNSZoneRequest struct {
	ZoneID entity.DNSZoneID
}

type CreateDNSRecordRequest struct {
	ZoneID     entity.DNSZoneID
	Name       string
	Type       entity.DNSRecordType
	Value      string
	TTLSeconds uint32
	Priority   *uint16
}

type DeleteDNSRecordRequest struct {
	RecordID entity.DNSRecordID
}

type ApplyBlueprintRequest struct {
	Blueprint entity.Blueprint
}

type ValidateBlueprintRequest struct {
	Blueprint entity.Blueprint
}

type ImportBlueprintRequest struct {
	Blueprint entity.Blueprint
}

type ExportBlueprintRequest struct {
	BlueprintID entity.BlueprintID
}

type DestroyBlueprintRequest struct {
	BlueprintID entity.BlueprintID
}
