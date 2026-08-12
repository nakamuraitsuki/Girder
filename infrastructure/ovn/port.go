package ovn

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var ErrLogicalSwitchPortAlreadyExists = fmt.Errorf("logical switch port already exists")
var ErrLogicalSwitchPortNotFound = fmt.Errorf("logical switch port not found")

// CreateLogicalSwitchPort creates a new OVN Logical_Switch_Port.
func (c *Client) CreateLogicalSwitchPort(
	ctx context.Context,
	port *LogicalSwitchPort,
) (*LogicalSwitchPort, error) {
	// validate input.
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

	ops, err := c.nb.Create(ctx, port)
	if err != nil {
		return nil, fmt.Errorf("failed to create logical switch port %q: %w", port.Name, err)
	}

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
	// TODO: list Logical_Switch_Port records.

	// TODO: handle not found / ambiguous name.

	return nil, nil
}

// DeleteLogicalSwitchPort deletes the Logical_Switch_Port identified by name.
func (c *Client) DeleteLogicalSwitchPort(
	ctx context.Context,
	name string,
) error {
	// TODO: get the target port.

	// TODO: build delete operation.

	// TODO: transact delete operation.

	return nil
}
