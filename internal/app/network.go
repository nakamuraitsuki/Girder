package app

import (
	"context"

	"girder/internal/domain/entity"
)

type NetworkService interface {
	CreateSwitch(context.Context, CreateSwitchRequest) (entity.Switch, error)
	DeleteSwitch(context.Context, DeleteSwitchRequest) (entity.SwitchID, error)
	CreateRouter(context.Context, CreateRouterRequest) (entity.Router, error)
	DeleteRouter(context.Context, DeleteRouterRequest) (entity.RouterID, error)
	ConnectPort(context.Context, ConnectPortRequest) (entity.Port, error)
	DisconnectPort(context.Context, DisconnectPortRequest) (entity.Port, error)
	AddRoute(context.Context, AddRouteRequest) (entity.Route, error)
	DeleteRoute(context.Context, DeleteRouteRequest) (entity.RouteID, error)
	ConfigureNAT(context.Context, ConfigureNATRequest) (entity.NAT, error)
	ConfigureDHCP(context.Context, ConfigureDHCPRequest) (entity.DHCP, error)
	CreateDNSZone(context.Context, CreateDNSZoneRequest) (entity.DNSZone, error)
	DeleteDNSZone(context.Context, DeleteDNSZoneRequest) (entity.DNSZoneID, error)
	CreateDNSRecord(context.Context, CreateDNSRecordRequest) (entity.DNSRecord, error)
	DeleteDNSRecord(context.Context, DeleteDNSRecordRequest) (entity.DNSRecordID, error)
}
