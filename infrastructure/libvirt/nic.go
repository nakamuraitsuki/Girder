package libvirt

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

// AttachNIC attaches a new network interface to the specified domain.
// VM should not be running when this function is called.
func (c *Client) AttachNIC(
	domainName string,
	iface *libvirtxml.DomainInterface,
) (*libvirtxml.DomainInterface, error) {
	domain, err := c.conn.LookupDomainByName(domainName)
	if err != nil {
		return nil, fmt.Errorf("lookup domain: %w", err)
	}
	defer domain.Free()

	if iface.Source == nil {
		iface.Source = &libvirtxml.DomainInterfaceSource{
			Ethernet: &libvirtxml.DomainInterfaceSourceEthernet{},
		}
	}
	if iface.Model == nil {
		iface.Model = &libvirtxml.DomainInterfaceModel{Type: "virtio"}
	}
	if iface.Target == nil {
		iface.Target = &libvirtxml.DomainInterfaceTarget{
			Dev: fmt.Sprintf("tap-%s", randomHex(4)),
		}
	}

	xml, err := iface.Marshal()
	if err != nil {
		return nil, fmt.Errorf("marshal interface XML: %w", err)
	}

	// Check if the domain is active to determine the appropriate flags for attaching the device.
	active, err := domain.IsActive()
	if err != nil {
		return nil, fmt.Errorf("check domain active: %w", err)
	}

	// Determine the appropriate flags for attaching the device based on the domain's active state.
	bit := libvirt.DOMAIN_AFFECT_CONFIG
	if active {
		bit |= libvirt.DOMAIN_AFFECT_LIVE
	}
	flags := libvirt.DomainDeviceModifyFlags(bit)

	if err := domain.AttachDeviceFlags(xml, flags); err != nil {
		return nil, fmt.Errorf("attach device: %w", err)
	}

	actualIface, err := c.findInterfaceByTapDevice(domain, iface.Target.Dev)
	if err != nil {
		return nil, fmt.Errorf("find interface by tap device: %w", err)
	}

	return actualIface, nil
}

// ListNICs returns raw interface information from permanent domain XML.
// It does not include Girder-specific objects.
func (c *Client) ListNICs(domainName string) ([]libvirtxml.DomainInterface, error) {
	dom, err := c.conn.LookupDomainByName(domainName)
	if err != nil {
		return nil, fmt.Errorf("lookup domain: %w", err)
	}
	defer dom.Free()

	xmlDesc, err := dom.GetXMLDesc(libvirt.DOMAIN_XML_INACTIVE)
	if err != nil {
		return nil, fmt.Errorf("get xml desc: %w", err)
	}

	var domainXML libvirtxml.Domain
	if err := domainXML.Unmarshal(xmlDesc); err != nil {
		return nil, fmt.Errorf("unmarshal domain XML: %w", err)
	}

	return domainXML.Devices.Interfaces, nil
}

// findInterfaceByTapDevice is utility function to find a domain interface by its tap device name.
func (c *Client) findInterfaceByTapDevice(dom *libvirt.Domain, tapDev string) (*libvirtxml.DomainInterface, error) {
	xmlDesc, err := dom.GetXMLDesc(0)
	if err != nil {
		return nil, fmt.Errorf("get xml desc: %w", err)
	}

	var domainXML libvirtxml.Domain
	if err := domainXML.Unmarshal(xmlDesc); err != nil {
		return nil, fmt.Errorf("unmarshal domain XML: %w", err)
	}

	for i := range domainXML.Devices.Interfaces {
		iface := &domainXML.Devices.Interfaces[i]
		if iface.Target != nil && iface.Target.Dev == tapDev {
			return iface, nil
		}
	}

	return nil, fmt.Errorf("interface with target dev %s not found", tapDev)
}


// randomHex returns a random hexadecimal string of length 2*n.
func randomHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		// crypto/rand is not expected to fail in normal operation, so panic is acceptable here.
		panic(fmt.Sprintf("randomHex: crypto/rand read failed: %v", err))
	}
	return hex.EncodeToString(b)
}
