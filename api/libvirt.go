package api

import "context"

type DomainState string

const (
	DomainStateUndefined DomainState = "undefined"
	DomainStateRunning   DomainState = "running"
	DomainStatePaused    DomainState = "paused"
	DomainStateShutoff   DomainState = "shutoff"
)

type Connection struct {
	ID       string
	URI      string
	ReadOnly bool
}

type Domain struct {
	ID             string
	Name           string
	UUID           string
	State          DomainState
	CPUs           int
	MemoryMiB      int64
	Disks          []Device
	NetworkDevices []Device
	Metadata       Metadata
}

type Device struct {
	Type       string
	Source     string
	Target     string
	Model      string
	MACAddress string
	BootOrder  int
	ReadOnly   bool
	Metadata   Metadata
}

type DomainSnapshot struct {
	ID          string
	DomainID    string
	Name        string
	Description string
	Metadata    Metadata
}

type Network struct {
	ID          string
	Name        string
	UUID        string
	Bridge      string
	CIDR        string
	Gateway     string
	DHCPEnabled bool
	Metadata    Metadata
}

type ConnectRequest struct {
	URI      string
	Username string
	Password string
	ReadOnly bool
}

type DisconnectRequest struct {
	ConnectionID string
}

type ListDomainsRequest struct {
	ConnectionID string
	All          bool
}

type GetDomainRequest struct {
	ConnectionID string
	DomainID     string
}

type DefineDomainRequest struct {
	ConnectionID string
	Domain       Domain
}

type UndefineDomainRequest struct {
	ConnectionID string
	DomainID     string
}

type StartDomainRequest struct {
	ConnectionID string
	DomainID     string
}

type StopDomainRequest struct {
	ConnectionID string
	DomainID     string
}

type RestartDomainRequest struct {
	ConnectionID string
	DomainID     string
}

type DestroyDomainRequest struct {
	ConnectionID string
	DomainID     string
}

type AttachDeviceRequest struct {
	ConnectionID string
	DomainID     string
	Device       Device
}

type DetachDeviceRequest struct {
	ConnectionID string
	DomainID     string
	Device       Device
}

type CreateSnapshotRequest struct {
	ConnectionID string
	DomainID     string
	Name         string
	Description  string
}

type DeleteSnapshotRequest struct {
	ConnectionID string
	DomainID     string
	SnapshotID   string
}

type ListNetworksRequest struct {
	ConnectionID string
	All          bool
}

type DefineNetworkRequest struct {
	ConnectionID string
	Network      Network
}

type UndefineNetworkRequest struct {
	ConnectionID string
	NetworkID    string
}

type LibvirtAPI interface {
	Connect(context.Context, ConnectRequest) (Connection, error)
	Disconnect(context.Context, DisconnectRequest) error
	ListDomains(context.Context, ListDomainsRequest) ([]Domain, error)
	GetDomain(context.Context, GetDomainRequest) (Domain, error)
	DefineDomain(context.Context, DefineDomainRequest) (Domain, error)
	UndefineDomain(context.Context, UndefineDomainRequest) (string, error)
	StartDomain(context.Context, StartDomainRequest) (Domain, error)
	StopDomain(context.Context, StopDomainRequest) (Domain, error)
	RestartDomain(context.Context, RestartDomainRequest) (Domain, error)
	DestroyDomain(context.Context, DestroyDomainRequest) (string, error)
	AttachDevice(context.Context, AttachDeviceRequest) (Domain, error)
	DetachDevice(context.Context, DetachDeviceRequest) (Domain, error)
	CreateSnapshot(context.Context, CreateSnapshotRequest) (DomainSnapshot, error)
	DeleteSnapshot(context.Context, DeleteSnapshotRequest) (string, error)
	ListNetworks(context.Context, ListNetworksRequest) ([]Network, error)
	DefineNetwork(context.Context, DefineNetworkRequest) (Network, error)
	UndefineNetwork(context.Context, UndefineNetworkRequest) (string, error)
}
