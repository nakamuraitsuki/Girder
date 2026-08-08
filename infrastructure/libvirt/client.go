package libvirt

import (
	"github.com/libvirt/libvirt-go"
)

type Client struct {
	conn *libvirt.Connect
}

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
