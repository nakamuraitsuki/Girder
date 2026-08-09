package libvirt

import (
	"libvirt.org/go/libvirt"
)

const (
	storagePoolName = "default"

	alpineImageURL = "https://dl-cdn.alpinelinux.org/alpine/v3.22/releases/cloud/generic_alpine-3.22.0-x86_64-bios-cloudinit-metal-r0.qcow2"
	baseImageName  = "alpine-base.qcow2"
)

// Client is a wrapper around libvirt.Connect.
type Client struct {
	conn *libvirt.Connect
}

// NewClient creates a new libvirt client connect to the qemu:///system URI.
func NewClient() (*Client, error) {
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) Close() error {
	_, err := c.conn.Close()
	return err
}
