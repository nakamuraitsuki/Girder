package ovn

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ovn-kubernetes/libovsdb/client"
)

var ErrSwitchAlreadyExists = fmt.Errorf("switch already exists")
var ErrSwitchNotFound = fmt.Errorf("switch not found")

// CreateSwitch creates a new isolated virtual switch (OVN Logical_Switch).
// No Logical_Switch_Port is created here; port attachment is handled later
// by the NIC connecntion flow (ConnectNIC).
func (c *Client) CreateSwitch(ctx context.Context, name string) (*LogicalSwitch, error) {
	_, err := c.GetSwitch(ctx, name);
	switch {
	case err == nil:
		// Already exists, return error
		return nil, fmt.Errorf("failed to create switch %q: %w", name, ErrSwitchAlreadyExists)
	case errors.Is(err, ErrSwitchNotFound):
		// Not found, proceed to create
	default:
		// Unexpected error, return it
		return nil, fmt.Errorf("failed to check existing switch %q: %w", name, err)
	}

	sw := &LogicalSwitch{
		UUID: uuid.NewString(),
		Name: name,
	}

	ops, err := c.nb.Create(sw)
	if err != nil {
		return nil, fmt.Errorf("failed to build create op for switch %q: %w", name, err)
	}

	if _, err := c.nb.Transact(ctx, ops...); err != nil {
		return nil, fmt.Errorf("failed to create switch %q: %w", name, err)
	}

	return sw, nil
}

// GetSwitch returns the virtual switch identified by name.
func (c *Client) GetSwitch(ctx context.Context, name string) (*LogicalSwitch, error) {
	sw := &LogicalSwitch{Name: name}
	if err := c.nb.Get(ctx, sw); err != nil {
		if errors.Is(err, client.ErrNotFound) {
			return nil, fmt.Errorf("switch %q not found: %w", name, ErrSwitchNotFound)
		}
		return nil, fmt.Errorf("failed to get switch %q: %w", name, err)
	}
	return sw, nil
}

// DeleteSwitch deletes the virtual switch identified by name.
func (c *Client) DeleteSwitch(ctx context.Context, name string) error {
	sw, err := c.GetSwitch(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to get switch %q: %w", name, err)
	}

	ops, err := c.nb.Where(sw).Delete()
	if err != nil {
		return fmt.Errorf("failed to build delete op for switch %q: %w", name, err)
	}

	if _, err := c.nb.Transact(ctx, ops...); err != nil {
		return fmt.Errorf("failed to delete switch %q: %w", name, err)
	}

	return nil
}
