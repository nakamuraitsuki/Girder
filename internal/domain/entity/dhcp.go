package entity

import (
	"fmt"
	"strings"
	
	"girder/internal/domain/value"
)

type DHCP struct {
	ID           DHCPID
	Name         string
	Description  string
	Network      value.CIDR
	RangeStart   value.IPv4
	RangeEnd     value.IPv4
	Gateway      value.IPv4
	DNSServers   []value.IPv4
	LeaseSeconds int
	Enabled      bool
}

func NewDHCP(name string, network value.CIDR, rangeStart, rangeEnd, gateway value.IPv4, leaseSeconds int) (DHCP, error) {
	if strings.TrimSpace(name) == "" {
		return DHCP{}, fmt.Errorf("%w: dhcp name", ErrEmptyName)
	}
	if network.IsZero() || rangeStart.IsZero() || rangeEnd.IsZero() || gateway.IsZero() {
		return DHCP{}, fmt.Errorf("%w: dhcp range", ErrInvalidResource)
	}
	if leaseSeconds < 0 {
		return DHCP{}, fmt.Errorf("%w: dhcp lease", ErrInvalidResource)
	}
	return DHCP{ID: NewDHCPID(), Name: strings.TrimSpace(name), Network: network, RangeStart: rangeStart, RangeEnd: rangeEnd, Gateway: gateway, LeaseSeconds: leaseSeconds, Enabled: true, DNSServers: []value.IPv4{}}, nil
}

func (d *DHCP) Rename(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("%w: dhcp name", ErrEmptyName)
	}
	d.Name = strings.TrimSpace(name)
	return nil
}

func (d *DHCP) SetDescription(description string) {
	d.Description = strings.TrimSpace(description)
}

func (d *DHCP) Enable()  { d.Enabled = true }
func (d *DHCP) Disable() { d.Enabled = false }

func (d *DHCP) AddDNSServer(address value.IPv4) error {
	if address.IsZero() {
		return fmt.Errorf("%w: dns server", ErrInvalidResource)
	}
	for _, existing := range d.DNSServers {
		if existing == address {
			return nil
		}
	}
	d.DNSServers = append(d.DNSServers, address)
	return nil
}