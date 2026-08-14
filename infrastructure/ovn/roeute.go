package ovn

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var ErrLogicalRouterNotFound = fmt.Errorf("logical router not found")

var ErrLogicalRouterAlreadyExists = fmt.Errorf("logical router already exists")

// CreateRouter creates a new OVN Logical_Router.
func (c *Client) CreateRouter(ctx context.Context, lr *LogicalRouter) (*LogicalRouter, error) {
	_, err := c.GetRouter(ctx, lr.Name)
	switch {
	case err == nil:
		return nil, fmt.Errorf("failed to create logical router %q: %w", lr.Name, ErrLogicalRouterAlreadyExists)
	case errors.Is(err, ErrLogicalRouterNotFound):
		// proceed
	default:
		return nil, fmt.Errorf("failed to create logical router %q: %w", lr.Name, err)
	}

	if lr.UUID == "" {
		lr.UUID = uuid.NewString()
	}

	ops, err := c.nb.Create(lr)
	if err != nil {
		return nil, fmt.Errorf("failed to build create op for logical router %q: %w", lr.Name, err)
	}

	if _, err := c.nb.Transact(ctx, ops...); err != nil {
		return nil, fmt.Errorf("failed to create logical router %q: %w", lr.Name, err)
	}

	return lr, nil
}

// GetLogicalRouter returns the logical router identified by name.
//
// Name uniqueness is not enforced by ONV itself, but Girder enforce it
// at creation time.
func (c *Client) GetRouter(ctx context.Context, name string) (*LogicalRouter, error) {
	var routers []LogicalRouter
	if err := c.nb.List(ctx, &routers); err != nil {
		return nil, fmt.Errorf("failed to list logical routers: %w", err)
	}

	var found *LogicalRouter
	for _, r := range routers {
		if r.Name != name {
			continue
		}
		if found != nil {
			return nil, fmt.Errorf("logical router %q is ambiguous: multiple routers share this name", name)
		}
		found = &r
	}

	if found == nil {
		return nil, fmt.Errorf("logical router %q not found: %w", name, ErrLogicalRouterNotFound)
	}
	return found, nil
}

// DeleteLogicalRouter deletes the logical router identified by name.
func (c *Client) DeleteRouter(ctx context.Context, name string) error {
	lr, err := c.GetRouter(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to get logical router %q: %w", name, err)
	}

	ops, err := c.nb.Where(lr).Delete()
	if err != nil {
		return fmt.Errorf("failed to build delete op for logical router %q: %w", name, err)
	}

	if _, err := c.nb.Transact(ctx, ops...); err != nil {
		return fmt.Errorf("failed to delete logical router %q: %w", name, err)
	}

	return nil
}
