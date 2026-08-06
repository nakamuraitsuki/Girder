package entity

import (
	"fmt"
	"strings"
	
	"girder/internal/domain/value"
)

type Port struct {
	ID               PortID
	Name             string
	Description      string
	MACAddress       value.MACAddress
	IPv4Addresses    map[value.IPv4]struct{}
	IPv6Addresses    map[value.IPv6]struct{}
	AttachedVMID     *VMID
	AttachedSwitchID *SwitchID
	AttachedRouterID *RouterID
	State            PortState
}

func NewPort(name string, mac value.MACAddress) (Port, error) {
	if strings.TrimSpace(name) == "" {
		return Port{}, fmt.Errorf("%w: port name", ErrEmptyName)
	}
	if mac.IsZero() {
		return Port{}, fmt.Errorf("%w: port mac", ErrInvalidResource)
	}
	return Port{
		ID:            NewPortID(),
		Name:          strings.TrimSpace(name),
		MACAddress:    mac,
		IPv4Addresses: map[value.IPv4]struct{}{},
		IPv6Addresses: map[value.IPv6]struct{}{},
		State:         PortStateDetached,
	}, nil
}

func (p *Port) Rename(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("%w: port name", ErrEmptyName)
	}
	p.Name = strings.TrimSpace(name)
	return nil
}

func (p *Port) SetDescription(description string) {
	p.Description = strings.TrimSpace(description)
}

func (p *Port) attachExclusive(vmID *VMID, switchID *SwitchID, routerID *RouterID) error {
	if p.AttachedVMID != nil || p.AttachedSwitchID != nil || p.AttachedRouterID != nil {
		return fmt.Errorf("%w: port already attached", ErrAlreadyAttached)
	}
	p.AttachedVMID = vmID
	p.AttachedSwitchID = switchID
	p.AttachedRouterID = routerID
	p.State = PortStateAttached
	return nil
}

func (p *Port) AttachToVM(vmID VMID) error {
	if vmID == "" {
		return fmt.Errorf("%w: vm attachment", ErrInvalidRelation)
	}
	return p.attachExclusive(&vmID, nil, nil)
}

func (p *Port) AttachToSwitch(switchID SwitchID) error {
	if switchID == "" {
		return fmt.Errorf("%w: switch attachment", ErrInvalidRelation)
	}
	return p.attachExclusive(nil, &switchID, nil)
}

func (p *Port) AttachToRouter(routerID RouterID) error {
	if routerID == "" {
		return fmt.Errorf("%w: router attachment", ErrInvalidRelation)
	}
	return p.attachExclusive(nil, nil, &routerID)
}

func (p *Port) Detach() {
	p.AttachedVMID = nil
	p.AttachedSwitchID = nil
	p.AttachedRouterID = nil
	p.State = PortStateDetached
}

func (p *Port) AddIPv4(address value.IPv4) error {
	if address.IsZero() {
		return fmt.Errorf("%w: ipv4 address", ErrInvalidResource)
	}
	if p.IPv4Addresses == nil {
		p.IPv4Addresses = map[value.IPv4]struct{}{}
	}
	p.IPv4Addresses[address] = struct{}{}
	return nil
}

func (p *Port) RemoveIPv4(address value.IPv4) {
	if p.IPv4Addresses == nil {
		return
	}
	delete(p.IPv4Addresses, address)
}

func (p *Port) AddIPv6(address value.IPv6) error {
	if address.IsZero() {
		return fmt.Errorf("%w: ipv6 address", ErrInvalidResource)
	}
	if p.IPv6Addresses == nil {
		p.IPv6Addresses = map[value.IPv6]struct{}{}
	}
	p.IPv6Addresses[address] = struct{}{}
	return nil
}

func (p *Port) RemoveIPv6(address value.IPv6) {
	if p.IPv6Addresses == nil {
		return
	}
	delete(p.IPv6Addresses, address)
}
