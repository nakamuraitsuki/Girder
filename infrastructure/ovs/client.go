package ovs

import (
	"context"
	"fmt"

	"github.com/ovn-kubernetes/libovsdb/client"
)

// Client is a client for Open vSwitch database.
type Client struct {
	ovsdb client.Client
}

// Connect connects to the Open vSwitch database.
// (e.g. unix:/var/run/openvswitch/db.sock)
func Connect(ctx context.Context, endpoint string) (*Client, error) {
	dbModel, err := DatabaseModel()
	if err != nil {
		return nil, fmt.Errorf("failed to create database model: %w", err)
	}

	ovsdb, err := client.NewOVSDBClient(dbModel, client.WithEndpoint(endpoint))
	if err != nil {
		return nil, fmt.Errorf("failed to create OVSDB client: %w", err)
	}

	if err := ovsdb.Connect(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to OVS DB: %w", err)
	}

	if _, err := ovsdb.MonitorAll(ctx); err != nil {
		return nil, fmt.Errorf("failed to subscribe to OVS DB tables: %w", err)
	}

	return &Client{ovsdb:ovsdb}, nil
}

// Close disconnect from the OVS DB.
func (c *Client) Close() {
	c.ovsdb.Disconnect()
}

// getIntegrationBridge returns the integration bridge of the OVS DB.
// No Cache is used, so it always queries the OVS DB.
func (c *Client) getIntegrationBridge(ctx context.Context) (*Bridge, error) {
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
