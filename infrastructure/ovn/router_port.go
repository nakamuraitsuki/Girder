package ovn

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ovn-kubernetes/libovsdb/model"
)

var ErrLogicalRouterPortAlreadyExists = fmt.Errorf("logical router port already exists")
var ErrLogicalRouterPortNotFound = fmt.Errorf("logical router port not found")

// CreateLogicalRouterPort creates a new OVN Logical_Router_Port.
//
// MAC is intentionally left empty unless port.MAC is set by the caller;
// OVN assigns it automatically when empty. Callers that need the actual
// MAC should call GetLogicalRouterPort afterwards.
func (c *Client) CreateLogicalRouterPort(
	ctx context.Context,
	router *LogicalRouter,
	port *LogicalRouterPort,
) (*LogicalRouterPort, error) {
	if router == nil {
		return nil, errors.New("logical router is nil")
	}
	if port == nil {
		return nil, errors.New("logical router port is nil")
	}
	if port.Name == "" {
		return nil, errors.New("logical router port name is required")
	}

	_, err := c.GetLogicalRouterPort(ctx, port.Name)
	switch {
	case err == nil:
		return nil, fmt.Errorf(
			"failed to create logical router port %q: %w",
			port.Name,
			ErrLogicalRouterPortAlreadyExists,
		)
	case errors.Is(err, ErrLogicalRouterPortNotFound):
		// proceed
	default:
		return nil, fmt.Errorf(
			"failed to create logical router port %q: %w",
			port.Name,
			err,
		)
	}

	port.UUID = uuid.NewString()

	createOps, err := c.nb.Create(port)
	if err != nil {
		return nil, fmt.Errorf("failed to create logical router port %q: %w", port.Name, err)
	}

	mutateOps, err := c.nb.Where(router).Mutate(router,
		model.Mutation{
			Field:   &router.Ports,
			Mutator: "insert",
			Value:   []string{port.UUID},
		})
	if err != nil {
		return nil, fmt.Errorf(
			"failed to build mutate op attaching port %q to router %q: %w",
			port.Name, router.Name, err,
		)
	}

	ops := append(createOps, mutateOps...)

	if _, err := c.nb.Transact(ctx, ops...); err != nil {
		return nil, fmt.Errorf("failed to create logical router port %q: %w", port.Name, err)
	}

	return port, nil
}

// GetLogicalRouterPort returns the Logical_Router_Port identified by name.
func (c *Client) GetLogicalRouterPort(
	ctx context.Context,
	name string,
) (*LogicalRouterPort, error) {
	var ports []LogicalRouterPort
	if err := c.nb.List(ctx, &ports); err != nil {
		return nil, fmt.Errorf("failed to list logical router ports: %w", err)
	}

	var found *LogicalRouterPort
	for _, port := range ports {
		if port.Name != name {
			continue
		}

		if found != nil {
			return nil, fmt.Errorf(
				"logical router port %q is ambiguous: multiple ports share this name (environment may have been modified outside Girder)",
				name,
			)
		}

		found = &port
	}

	if found == nil {
		return nil, fmt.Errorf("logical router port %q not found: %w", name, ErrLogicalRouterPortNotFound)
	}
	return found, nil
}

// DeleteLogicalRouterPort deletes the Logical_Router_Port identified by name.
func (c *Client) DeleteLogicalRouterPort(
	ctx context.Context,
	router *LogicalRouter,
	name string,
) error {
	if router == nil {
		return errors.New("logical router is nil")
	}

	port, err := c.GetLogicalRouterPort(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to get logical router port %q: %w", name, err)
	}

	deleteOps, err := c.nb.Where(port).Delete()
	if err != nil {
		return fmt.Errorf("failed to build delete op for logical router port %q: %w", name, err)
	}

	mutateOps, err := c.nb.Where(router).Mutate(router,
		model.Mutation{
			Field:   &router.Ports,
			Mutator: "delete",
			Value:   []string{port.UUID},
		})
	if err != nil {
		return fmt.Errorf(
			"failed to build mutate op detaching port %q from router %q: %w",
			port.Name, router.Name, err,
		)
	}

	ops := append(deleteOps, mutateOps...)

	if _, err := c.nb.Transact(ctx, ops...); err != nil {
		return fmt.Errorf("failed to delete logical router port %q: %w", name, err)
	}

	return nil
}