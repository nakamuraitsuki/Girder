package api

import "context"

type LogicalSwitch struct {
	ID         string
	Name       string
	ExternalID string
	Ports      []LogicalPort
	ACLs       []ACL
	Metadata   Metadata
}

type LogicalRouter struct {
	ID          string
	Name        string
	ExternalID  string
	Ports       []LogicalPort
	Routes      []Route
	NATs        []NAT
	DHCPOptions []DHCPOption
	Metadata    Metadata
}

type LogicalPort struct {
	ID        string
	Name      string
	Type      string
	ParentID  string
	Addresses []string
	Options   map[string]string
	Metadata  Metadata
}

type ACL struct {
	ID        string
	Priority  int
	Direction string
	Match     string
	Action    string
}

type Route struct {
	ID      string
	Prefix  string
	NextHop string
	Metric  int
}

type NAT struct {
	ID         string
	Type       string
	LogicalIP  string
	ExternalIP string
}

type DHCPOption struct {
	ID    string
	Key   string
	Value string
}

type ListLogicalSwitchesRequest struct {
	DriverID string
	All      bool
}

type GetLogicalSwitchRequest struct {
	DriverID string
	SwitchID string
}

type CreateLogicalSwitchRequest struct {
	DriverID string
	Switch   LogicalSwitch
}

type DeleteLogicalSwitchRequest struct {
	DriverID string
	SwitchID string
}

type ListLogicalRoutersRequest struct {
	DriverID string
	All      bool
}

type CreateLogicalRouterRequest struct {
	DriverID string
	Router   LogicalRouter
}

type DeleteLogicalRouterRequest struct {
	DriverID string
	RouterID string
}

type CreateLogicalPortRequest struct {
	DriverID string
	Port     LogicalPort
}

type DeleteLogicalPortRequest struct {
	DriverID string
	PortID   string
}

type ConfigureACLRequest struct {
	DriverID string
	SwitchID string
	ACL      ACL
}

type DeleteACLRequest struct {
	DriverID string
	SwitchID string
	ACLID    string
}

type ConfigureNATRequest struct {
	DriverID string
	RouterID string
	NAT      NAT
}

type DeleteNATRequest struct {
	DriverID string
	RouterID string
	NATID    string
}

type ConfigureDHCPOptionRequest struct {
	DriverID string
	RouterID string
	Option   DHCPOption
}

type DeleteDHCPOptionRequest struct {
	DriverID string
	RouterID string
	OptionID string
}

type OVNAPI interface {
	ListLogicalSwitches(context.Context, ListLogicalSwitchesRequest) ([]LogicalSwitch, error)
	GetLogicalSwitch(context.Context, GetLogicalSwitchRequest) (LogicalSwitch, error)
	CreateLogicalSwitch(context.Context, CreateLogicalSwitchRequest) (LogicalSwitch, error)
	DeleteLogicalSwitch(context.Context, DeleteLogicalSwitchRequest) (string, error)
	ListLogicalRouters(context.Context, ListLogicalRoutersRequest) ([]LogicalRouter, error)
	CreateLogicalRouter(context.Context, CreateLogicalRouterRequest) (LogicalRouter, error)
	DeleteLogicalRouter(context.Context, DeleteLogicalRouterRequest) (string, error)
	CreateLogicalPort(context.Context, CreateLogicalPortRequest) (LogicalPort, error)
	DeleteLogicalPort(context.Context, DeleteLogicalPortRequest) (string, error)
	ConfigureACL(context.Context, ConfigureACLRequest) (ACL, error)
	DeleteACL(context.Context, DeleteACLRequest) (string, error)
	ConfigureNAT(context.Context, ConfigureNATRequest) (NAT, error)
	DeleteNAT(context.Context, DeleteNATRequest) (string, error)
	ConfigureDHCPOption(context.Context, ConfigureDHCPOptionRequest) (DHCPOption, error)
	DeleteDHCPOption(context.Context, DeleteDHCPOptionRequest) (string, error)
}
