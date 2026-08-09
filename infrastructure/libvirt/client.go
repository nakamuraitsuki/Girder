package libvirt

import (
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

// Client provides direct access to a libvirt connection.
//
// Libvirt remains the source of truth for VM state. Client does not maintain
// a separate representation of libvirt resources.
type Client struct {
	conn *libvirt.Connect
}

// NewClient connects to the system libvirt instance.
func NewClient() (*Client, error) {
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
	}, nil
}

// Close closes the libvirt connection.
func (c *Client) Close() error {
	_, err := c.conn.Close()
	return err
}

// CreateVM defines and starts a new virtual machine from the given libvirt domain object.
func (c *Client) CreateVM(domain *libvirtxml.Domain) (*libvirt.Domain, error) {
	xml, err := domain.Marshal()
	if err != nil {
		return nil, err
	}

	vm, err := c.conn.DomainDefineXML(xml)
	if err != nil {
		return nil, err
	}

	if err := vm.Create(); err != nil {
		return nil, err
	}

	return vm, nil
}