package ovs

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ovn-kubernetes/libovsdb/model"
)

var ErrBridgeAlreadyExists = fmt.Errorf("bridge already exists")
var ErrBridgeNotFound = fmt.Errorf("bridge not found")

// getIntegrationBridge returns the integration bridge of the OVS DB.
// No Cache is used, so it always queries the OVS DB.
func (c *Client) GetIntegrationBridge(ctx context.Context) (*Bridge, error) {
	var bridges []Bridge
	if err := c.ovsdb.List(ctx, &bridges); err != nil {
		return nil, fmt.Errorf("failed to list bridges: %w", err)
	}

	var matches []Bridge
	for _, bridge := range bridges {
		if bridge.Name == integrationBridgeName {
			matches = append(matches, bridge)
		}
	}

	switch len(matches) {
	case 0:
		return nil, fmt.Errorf("integration bridge %q not found (ovs setup incomplete?)", integrationBridgeName)
	case 1:
		return &matches[0], nil
	default:
		return nil, fmt.Errorf("multiple integration bridges found: %v", matches)
	}
}

// CreateBridge creates a new, empty top-level OVS bridge and registers it
// in the root Open_vSwitch row's bridges column. No port is attached here.
func (c *Client) CreateBridge(ctx context.Context, name string) (*Bridge, error) {
	if name == "" {
		return nil, fmt.Errorf("bridge name cannot be empty")
	}

	_, err := c.GetBridge(ctx, name)
	switch {
	case err == nil:
		return nil, fmt.Errorf("bridge %q already exists: %w", name, ErrBridgeAlreadyExists)
	case errors.Is(err, ErrBridgeNotFound):
		// processd
	default:
		return nil, fmt.Errorf("failed to get bridge: %w", err)
	}

	ovsRow, err := c.getRootOpenVSwitch(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get root Open_vSwitch row: %w", err)
	}

	bridge := &Bridge{
		UUID: uuid.NewString(),
		Name: name,
	}
	bridgeOps, err := c.ovsdb.Create(bridge)
	if err != nil {
		return nil, fmt.Errorf("failed to create bridge: %w", err)
	}

	mutateOps, err := c.ovsdb.Where(ovsRow).Mutate(ovsRow,
		model.Mutation{
			Field:   &ovsRow.Bridges,
			Mutator: "insert",
			Value:   []string{bridge.UUID},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create mutation for Open_vSwitch row: %w", err)
	}

	ops := append(bridgeOps, mutateOps...)
	if _, err := c.ovsdb.Transact(ctx, ops...); err != nil {
		return nil, fmt.Errorf("failed to transact: %w", err)
	}

	return bridge, nil
}

// GetBridge returns the Bridge with the given name.
func (c *Client) GetBridge(ctx context.Context, name string) (*Bridge, error) {
	var bridges []Bridge
	if err := c.ovsdb.List(ctx, &bridges); err != nil {
		return nil, fmt.Errorf("failed to list bridges: %w", err)
	}

	for i := range bridges {
		if bridges[i].Name == name {
			return &bridges[i], nil
		}
	}

	return nil, fmt.Errorf("bridge %q not found: %w", name, ErrBridgeNotFound)
}

// getRootOpenVSwitch returns the single root row of the Open_vSwitch table.
func (c *Client) getRootOpenVSwitch(ctx context.Context) (*OpenVSwitch, error) {
	var rows []OpenVSwitch
	if err := c.ovsdb.List(ctx, &rows); err != nil {
		return nil, fmt.Errorf("failed to list Open_vSwitch rows: %w", err)
	}

	switch len(rows) {
	case 0:
		return nil, errors.New("no Open_vSwitch row found (ovs setup incomplete?)")
	case 1:
		return &rows[0], nil
	default:
		return nil, fmt.Errorf("multiple Open_vSwitch rows found: %d", len(rows))
	}
}
