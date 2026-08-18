package ovn

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ovn-kubernetes/libovsdb/model"
)

var ErrLogicalSwitchPortAlreadyExists = fmt.Errorf("logical switch port already exists")
var ErrLogicalSwitchPortNotFound = fmt.Errorf("logical switch port not found")

// CreateLogicalSwitchPort creates a new OVN Logical_Switch_Port.
func (c *Client) CreateLogicalSwitchPort(
	ctx context.Context,
	sw *LogicalSwitch,
	port *LogicalSwitchPort,
) (*LogicalSwitchPort, error) {
	// validate input.
	if sw == nil {
		return nil, errors.New("logical switch is nil")
	}
	if port == nil {
		return nil, errors.New("logical switch port is nil")
	}
	if port.Name == "" {
		return nil, errors.New("logical switch port name is required")
	}

	_, err := c.GetLogicalSwitchPort(ctx, port.Name)
	switch {
	case err == nil:
		return nil, fmt.Errorf(
			"failed to create logical switch port %q: %w",
			port.Name,
			ErrLogicalSwitchPortAlreadyExists,
		)
	case errors.Is(err, ErrLogicalSwitchPortNotFound):
		// proceed
	default:
		return nil, fmt.Errorf(
			"failed to create logical switch port %q: %w",
			port.Name,
			err,
		)
	}

	port.UUID = uuid.NewString()

	createOps, err := c.nb.Create(port)
	if err != nil {
		return nil, fmt.Errorf("failed to create logical switch port %q: %w", port.Name, err)
	}

	mutateOps, err := c.nb.Where(sw).Mutate(sw,
	model.Mutation{
		Field: &sw.Ports,
		Mutator: "insert",
		Value: []string{port.UUID},
	})
	if err != nil {
		return nil, fmt.Errorf(
			"failed to build mutate op attaching port %q to switch %q: %w",
			port.Name, sw.Name, err,
		)
	}

	ops := append(createOps, mutateOps...)

	if _, err := c.nb.Transact(ctx, ops...); err != nil {
		return nil, fmt.Errorf("failed to create logical switch port %q: %w", port.Name, err)
	}

	return port, nil
}

// GetLogicalSwitchPort returns the Logical_Switch_Port identified by name.
func (c *Client) GetLogicalSwitchPort(
	ctx context.Context,
	name string,
) (*LogicalSwitchPort, error) {
	var ports []LogicalSwitchPort
	if err := c.nb.List(ctx, &ports); err != nil {
		return nil, fmt.Errorf("failed to list logical switch ports: %w", err)
	}

	for _, port := range ports {
		if port.Name != name {
			continue
		}
		return &port, nil
	}

	return nil, fmt.Errorf(
		"logical switch port %q not found: %w",
		name,
		ErrLogicalSwitchPortNotFound,
	)
}

// DeleteLogicalSwitchPort deletes the Logical_Switch_Port identified by name.
func (c *Client) DeleteLogicalSwitchPort(
	ctx context.Context,
	sw *LogicalSwitch,
	name string,
) error {
	if sw == nil {
		return errors.New("logical switch is nil")
	}

	port, err := c.GetLogicalSwitchPort(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to get logical switch port %q: %w", name, err)
	}

	deleteOps, err := c.nb.Where(port).Delete()
	if err != nil {
		return fmt.Errorf("failed to build delete op for logical switch port %q: %w", name, err)
	}

	mutateOps, err := c.nb.Where(sw).Mutate(sw,
		model.Mutation{
			Field: &sw.Ports,
			Mutator: "delete",
			Value: []string{port.UUID},
		})
	if err != nil {
		return fmt.Errorf(
			"failed to build mutate op detaching port %q from switch %q: %w",
			port.Name, sw.Name, err,
		)
	}

	ops := append(deleteOps, mutateOps...)

	if _, err := c.nb.Transact(ctx, ops...); err != nil {
		return fmt.Errorf("failed to delete logical switch port %q: %w", name, err)
	}

	return nil
}
