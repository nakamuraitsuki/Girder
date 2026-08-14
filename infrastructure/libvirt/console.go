package libvirt

import (
	"fmt"

	"libvirt.org/go/libvirt"
)

// OpenConsole opens a console connection to the specified virtual machine using the provided protocol (e.g., "vnc", "spice").
func (c *Client) OpenConsole(name string) (*libvirt.Stream, error) {
	domain, err := c.conn.LookupDomainByName(name)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup domain: %w", err)
	}

	stream, err := c.conn.NewStream(0)
	if err != nil {
		return nil, fmt.Errorf("failed to create new stream: %w", err)
	}

	if err := domain.OpenConsole("", stream, 0); err != nil {
		stream.Abort()
		stream.Free()
		domain.Free()
		return nil, fmt.Errorf("failed to open console: %w", err)
	}

	domain.Free()

	return stream, nil
}
