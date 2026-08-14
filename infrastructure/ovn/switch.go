package ovn

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var ErrSwitchAlreadyExists = fmt.Errorf("switch already exists")
var ErrSwitchNotFound = fmt.Errorf("switch not found")

// CreateSwitch creates a new isolated virtual switch (OVN Logical_Switch).
// No Logical_Switch_Port is created here; port attachment is handled later
// by the NIC connecntion flow (ConnectNIC).
func (c *Client) CreateSwitch(ctx context.Context, sw *LogicalSwitch) (*LogicalSwitch, error) {
	_, err := c.GetSwitch(ctx, sw.Name)
	switch {
	case err == nil:
		// Already exists, return error
		return nil, fmt.Errorf("failed to create switch %q: %w", sw.Name, ErrSwitchAlreadyExists)
	case errors.Is(err, ErrSwitchNotFound):
		// Not found, proceed to create
	default:
		// Unexpected error, return it
		return nil, fmt.Errorf("failed to check existing switch %q: %w", sw.Name, err)
	}

	if sw.UUID == "" {
		sw.UUID = uuid.NewString()
	}

	ops, err := c.nb.Create(sw)
	if err != nil {
		return nil, fmt.Errorf("failed to build create op for switch %q: %w", sw.Name, err)
	}

	if _, err := c.nb.Transact(ctx, ops...); err != nil {
		return nil, fmt.Errorf("failed to create switch %q: %w", sw.Name, err)
	}

	return sw, nil
}

// GetSwitch returns the virtual switch identified by name.
//
// Name uniqueness is not enforced by OVN itself, but Girder enforces it
// at creation time (see CreateSwitch). If multiple switches share the same
// name, it means the environment was modified outside of Girder (e.g. via
// ovn-nbctl directly), and Girder cannot safely resolve which one the
// caller means. In that case, GetSwitch returns an error rather than
// silently picking one — the driver's state is the source of truth, and
// Girder should surface that ambiguity honestly instead of hiding it.
func (c *Client) GetSwitch(ctx context.Context, name string) (*LogicalSwitch, error) {
	var switches []LogicalSwitch
	if err := c.nb.List(ctx, &switches); err != nil {
		return nil, fmt.Errorf("failed to list switches: %w", err)
	}

	var found *LogicalSwitch
	for i := range switches {
		if switches[i].Name != name {
			continue
		}
		if found != nil {
			return nil, fmt.Errorf("switch %q is ambiguous: multiple switches share this name (environment may have been modified outside Girder)", name)
		}
		found = &switches[i]
	}

	if found == nil {
		return nil, fmt.Errorf("switch %q not found: %w", name, ErrSwitchNotFound)
	}
	return found, nil
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
