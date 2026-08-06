package entity

import (
	"fmt"
	"strings"
)

type Router struct {
	ID          RouterID
	Name        string
	Description string
	PortIDs     map[PortID]struct{}
	RouteIDs    map[RouteID]struct{}
	ACLIDs      map[ACLID]struct{}
	NATIDs      map[NATID]struct{}
	DHCPIDs     map[DHCPID]struct{}
	DNSZoneIDs  map[DNSZoneID]struct{}
	State       RouterState
}

func NewRouter(name string) (Router, error) {
	if strings.TrimSpace(name) == "" {
		return Router{}, fmt.Errorf("%w: router name", ErrEmptyName)
	}
	return Router{ID: NewRouterID(), Name: strings.TrimSpace(name), PortIDs: map[PortID]struct{}{}, RouteIDs: map[RouteID]struct{}{}, ACLIDs: map[ACLID]struct{}{}, NATIDs: map[NATID]struct{}{}, DHCPIDs: map[DHCPID]struct{}{}, DNSZoneIDs: map[DNSZoneID]struct{}{}, State: RouterStateUp}, nil
}

func (r *Router) Rename(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("%w: router name", ErrEmptyName)
	}
	r.Name = strings.TrimSpace(name)
	return nil
}

func (r *Router) SetDescription(description string) {
	r.Description = strings.TrimSpace(description)
}

func (r *Router) AttachPort(portID PortID) error {
	if portID == "" {
		return fmt.Errorf("%w: port", ErrInvalidRelation)
	}
	if r.PortIDs == nil {
		r.PortIDs = map[PortID]struct{}{}
	}
	r.PortIDs[portID] = struct{}{}
	return nil
}

func (r *Router) DetachPort(portID PortID) {
	if r.PortIDs == nil {
		return
	}
	delete(r.PortIDs, portID)
}

func (r *Router) AddRoute(routeID RouteID) error {
	if routeID == "" {
		return fmt.Errorf("%w: route", ErrInvalidRelation)
	}
	if r.RouteIDs == nil {
		r.RouteIDs = map[RouteID]struct{}{}
	}
	r.RouteIDs[routeID] = struct{}{}
	return nil
}

func (r *Router) DeleteRoute(routeID RouteID) {
	if r.RouteIDs == nil {
		return
	}
	delete(r.RouteIDs, routeID)
}

func (r *Router) AttachACL(aclID ACLID) error {
	if aclID == "" {
		return fmt.Errorf("%w: acl", ErrInvalidRelation)
	}
	if r.ACLIDs == nil {
		r.ACLIDs = map[ACLID]struct{}{}
	}
	r.ACLIDs[aclID] = struct{}{}
	return nil
}

func (r *Router) DetachACL(aclID ACLID) {
	if r.ACLIDs == nil {
		return
	}
	delete(r.ACLIDs, aclID)
}

func (r *Router) AddNAT(natID NATID) error {
	if natID == "" {
		return fmt.Errorf("%w: nat", ErrInvalidRelation)
	}
	if r.NATIDs == nil {
		r.NATIDs = map[NATID]struct{}{}
	}
	r.NATIDs[natID] = struct{}{}
	return nil
}

func (r *Router) DeleteNAT(natID NATID) {
	if r.NATIDs == nil {
		return
	}
	delete(r.NATIDs, natID)
}

func (r *Router) AddDHCP(dhcpID DHCPID) error {
	if dhcpID == "" {
		return fmt.Errorf("%w: dhcp", ErrInvalidRelation)
	}
	if r.DHCPIDs == nil {
		r.DHCPIDs = map[DHCPID]struct{}{}
	}
	r.DHCPIDs[dhcpID] = struct{}{}
	return nil
}

func (r *Router) DeleteDHCP(dhcpID DHCPID) {
	if r.DHCPIDs == nil {
		return
	}
	delete(r.DHCPIDs, dhcpID)
}

func (r *Router) AddDNSZone(zoneID DNSZoneID) error {
	if zoneID == "" {
		return fmt.Errorf("%w: dns zone", ErrInvalidRelation)
	}
	if r.DNSZoneIDs == nil {
		r.DNSZoneIDs = map[DNSZoneID]struct{}{}
	}
	r.DNSZoneIDs[zoneID] = struct{}{}
	return nil
}

func (r *Router) DeleteDNSZone(zoneID DNSZoneID) {
	if r.DNSZoneIDs == nil {
		return
	}
	delete(r.DNSZoneIDs, zoneID)
}

func (r *Router) SetState(state RouterState) {
	r.State = state
}
