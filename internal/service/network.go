package service

import (
	"context"

	app "girder/internal/app"
	"girder/internal/domain/entity"
)

type NetworkServiceImpl struct {
	SwitchDriver any
	RouterDriver any
	PortDriver   any
	RouteDriver  any
	NATDriver    any
	DHCPDriver   any
	DNSDriver    any
}

func (NetworkServiceImpl) CreateSwitch(context.Context, app.CreateSwitchRequest) (entity.Switch, error) {
	return entity.Switch{}, nil
}

func (NetworkServiceImpl) DeleteSwitch(context.Context, app.DeleteSwitchRequest) (entity.SwitchID, error) {
	return "", nil
}

func (NetworkServiceImpl) CreateRouter(context.Context, app.CreateRouterRequest) (entity.Router, error) {
	return entity.Router{}, nil
}

func (NetworkServiceImpl) DeleteRouter(context.Context, app.DeleteRouterRequest) (entity.RouterID, error) {
	return "", nil
}

func (NetworkServiceImpl) ConnectPort(context.Context, app.ConnectPortRequest) (entity.Port, error) {
	return entity.Port{}, nil
}

func (NetworkServiceImpl) DisconnectPort(context.Context, app.DisconnectPortRequest) (entity.Port, error) {
	return entity.Port{}, nil
}

func (NetworkServiceImpl) AddRoute(context.Context, app.AddRouteRequest) (entity.Route, error) {
	return entity.Route{}, nil
}

func (NetworkServiceImpl) DeleteRoute(context.Context, app.DeleteRouteRequest) (entity.RouteID, error) {
	return "", nil
}

func (NetworkServiceImpl) ConfigureNAT(context.Context, app.ConfigureNATRequest) (entity.NAT, error) {
	return entity.NAT{}, nil
}

func (NetworkServiceImpl) ConfigureDHCP(context.Context, app.ConfigureDHCPRequest) (entity.DHCP, error) {
	return entity.DHCP{}, nil
}

func (NetworkServiceImpl) CreateDNSZone(context.Context, app.CreateDNSZoneRequest) (entity.DNSZone, error) {
	return entity.DNSZone{}, nil
}

func (NetworkServiceImpl) DeleteDNSZone(context.Context, app.DeleteDNSZoneRequest) (entity.DNSZoneID, error) {
	return "", nil
}

func (NetworkServiceImpl) CreateDNSRecord(context.Context, app.CreateDNSRecordRequest) (entity.DNSRecord, error) {
	return entity.DNSRecord{}, nil
}

func (NetworkServiceImpl) DeleteDNSRecord(context.Context, app.DeleteDNSRecordRequest) (entity.DNSRecordID, error) {
	return "", nil
}
