package service

import (
	"fmt"
	"net/netip"

	"girder/internal/domain/entity"
	"girder/internal/domain/value"
)

type IPAddressAllocationService struct{}

func (IPAddressAllocationService) NextIPv4(network value.CIDR, used map[value.IPv4]struct{}) (value.IPv4, error) {
	if !network.IsIPv4() {
		return value.IPv4{}, fmt.Errorf("%w: ipv4 network required", entity.ErrInvalidResource)
	}
	prefix := network.Prefix()
	addr := prefix.Masked().Addr()
	if !addr.Is4() {
		return value.IPv4{}, fmt.Errorf("%w: ipv4 network required", entity.ErrInvalidResource)
	}
	current := addr.As4()
	for {
		ipAddr := netip.AddrFrom4(current)
		ip, err := value.NewIPv4(ipAddr.String())
		if err != nil {
			return value.IPv4{}, err
		}
		if used == nil {
			return ip, nil
		}
		if _, exists := used[ip]; !exists {
			return ip, nil
		}
		incrementIPv4(&current)
		if !prefix.Contains(netip.AddrFrom4(current)) {
			break
		}
	}
	return value.IPv4{}, fmt.Errorf("%w: no available ip", entity.ErrNotFound)
}

func incrementIPv4(addr *[4]byte) {
	for i := len(addr) - 1; i >= 0; i-- {
		addr[i]++
		if addr[i] != 0 {
			return
		}
	}
}

