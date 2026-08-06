package mock

import (
	"context"

	"girder/internal/app"
	"girder/internal/domain/entity"
)

type MockNetworkService struct{}

func (MockNetworkService) CreateSwitch(context.Context, app.CreateSwitchRequest) (entity.Switch, error) {
	return entity.Switch{}, nil
}

func (MockNetworkService) DeleteSwitch(context.Context, app.DeleteSwitchRequest) (entity.SwitchID, error) {
	return "", nil
}

func (MockNetworkService) CreateRouter(context.Context, app.CreateRouterRequest) (entity.Router, error) {
	return entity.Router{}, nil
}

func (MockNetworkService) DeleteRouter(context.Context, app.DeleteRouterRequest) (entity.RouterID, error) {
	return "", nil
}

func (MockNetworkService) ConnectPort(context.Context, app.ConnectPortRequest) (entity.Port, error) {
	return entity.Port{}, nil
}

func (MockNetworkService) DisconnectPort(context.Context, app.DisconnectPortRequest) (entity.Port, error) {
	return entity.Port{}, nil
}

func (MockNetworkService) AddRoute(context.Context, app.AddRouteRequest) (entity.Route, error) {
	return entity.Route{}, nil
}

func (MockNetworkService) DeleteRoute(context.Context, app.DeleteRouteRequest) (entity.RouteID, error) {
	return "", nil
}

func (MockNetworkService) ConfigureNAT(context.Context, app.ConfigureNATRequest) (entity.NAT, error) {
	return entity.NAT{}, nil
}

func (MockNetworkService) ConfigureDHCP(context.Context, app.ConfigureDHCPRequest) (entity.DHCP, error) {
	return entity.DHCP{}, nil
}

func (MockNetworkService) CreateDNSZone(context.Context, app.CreateDNSZoneRequest) (entity.DNSZone, error) {
	return entity.DNSZone{}, nil
}

func (MockNetworkService) DeleteDNSZone(context.Context, app.DeleteDNSZoneRequest) (entity.DNSZoneID, error) {
	return "", nil
}

func (MockNetworkService) CreateDNSRecord(context.Context, app.CreateDNSRecordRequest) (entity.DNSRecord, error) {
	return entity.DNSRecord{}, nil
}

func (MockNetworkService) DeleteDNSRecord(context.Context, app.DeleteDNSRecordRequest) (entity.DNSRecordID, error) {
	return "", nil
}
