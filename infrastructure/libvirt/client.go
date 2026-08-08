package libvirt

import (
    "libvirt.org/go/libvirt"
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

// CreateVM は、指定されたXML定義に基づいて仮想マシンを作成する。
func (c *Client) CreateVM(xml string) (*libvirt.Domain, error) {
	domain, err := c.conn.DomainDefineXML(xml)
	if err != nil {
		return nil, err
	}

	if err := domain.Create(); err != nil {
		return nil, err
	}

	return domain, nil
}