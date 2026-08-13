package ovs

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ovn-kubernetes/libovsdb/model"
)

var ErrPortAlreadyExists = fmt.Errorf("port already exists")
var ErrPortNotFound = fmt.Errorf("port not found")

// AddPort creates a Port and Interface on br-int and connects it to bridge.
func (c *Client) AddPort(ctx context.Context, port *Port, iface *Interface) (*Port, error) {
	if port == nil {
		return nil, errors.New("port is nil")
	}
	if iface == nil {
		return nil , errors.New("interface is nil")
	}
	if port.Name == "" {
		return nil, errors.New("port name is empty")
	}
	if iface.Name == "" {
		return nil, errors.New("interface name is empty")
	}

	_, err := c.GetPort(ctx, port.Name)
	switch {
	case err == nil:
		return nil, fmt.Errorf("%w: %s", ErrPortAlreadyExists, port.Name)
	case errors.Is(err, ErrPortNotFound):
		// processd
	default:
		return nil, fmt.Errorf("failed to get port: %w", err)
	}

	bridge, err := c.getIntegrationBridge(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration bridge: %w", err)
	}

	iface.UUID = uuid.NewString()
	ifaceOps, err := c.ovsdb.Create(iface)
	if err != nil {
		return nil, fmt.Errorf("failed to create interface: %w", err)
	}

	port.UUID = uuid.NewString()
	port.Interfaces = []string{iface.UUID}
	portOps, err := c.ovsdb.Create(port)
	if err != nil {
		return nil, fmt.Errorf("failed to create port: %w", err)
	}

	mutateOps, err := c.ovsdb.Where(bridge).Mutate(bridge,
		model.Mutation{
			Field: &bridge.Ports,
			Mutator: "insert",
			Value: []string{port.UUID},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create mutation for bridge: %w", err)
	}

	ops := append(ifaceOps, portOps...)
	ops = append(ops, mutateOps...)

	if _, err := c.ovsdb.Transact(ctx, ops...); err != nil {
		return nil, fmt.Errorf("failed to transact: %w", err)
	}

	return port, nil
}

// GetPort returns the Port with the given name. If multiple ports share the same name, an error is returned.
func (c *Client) GetPort(ctx context.Context, name string) (*Port, error) {
	var ports []Port
	if err := c.ovsdb.List(ctx, &ports); err != nil {
		return nil, fmt.Errorf("failed to list ports: %w", err)
	}

	var found *Port
	for i := range ports {
		if ports[i].Name != name {
			continue
		}
		if found != nil {
			return nil, fmt.Errorf(
				"port %q is ambiguous: multiple ports share this name (environment may have been modified outside Girder)",
				name,
			)
		}
		found = &ports[i]
	}

	if found == nil {
		return nil, fmt.Errorf("port %q not found: %w", name, ErrPortNotFound)
	}
	return found, nil
}

// RemovePort removes the Port with the given name. If multiple ports share the same name, an error is returned.
func (c *Client) RemovePort(ctx context.Context, name string) error {
	bridge, err := c.getIntegrationBridge(ctx)
	if err != nil {
		return fmt.Errorf("failed to remove port %q: %w", name, err)
	}

	port, err := c.GetPort(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to get port %q: %w", name, err)
	}

	mutateOps, err := c.ovsdb.Where(bridge).Mutate(bridge,
		model.Mutation{
			Field:   &bridge.Ports,
			Mutator: "delete",
			Value:   []string{port.UUID},
		},
	)
	if err != nil {
		return fmt.Errorf(
			"failed to build mutate op detaching port %q from bridge %q: %w",
			name, bridge.Name, err,
		)
	}

	portDeleteOps, err := c.ovsdb.Where(port).Delete()
	if err != nil {
		return fmt.Errorf("failed to build delete op for port %q: %w", name, err)
	}

	ops := append(mutateOps, portDeleteOps...)
	if _, err := c.ovsdb.Transact(ctx, ops...); err != nil {
		return fmt.Errorf(
			"failed to remove port %q from bridge %q: %w",
			name, bridge.Name, err,
		)
	}

	return nil
}