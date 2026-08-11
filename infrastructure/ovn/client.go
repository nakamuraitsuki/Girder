package ovn

import (
	"context"
	"fmt"

	"github.com/ovn-kubernetes/libovsdb/client"
)

// Client provides connection to OVN Northbound DB.
type Client struct {
	nb client.Client
}

// Connect connects to the OVN Northbound DB at the given endpoint
// (e.g. unix:/var/run/ovn/ovnnb_db.sock or tcp:127.0.0.1:6641)
// and returns a Client instance.
func Connect(ctx context.Context, endpoint string) (*Client, error) {
	dbModel, err := DatabaseModel()
	if err != nil {
		return nil, fmt.Errorf("failed to get database model: %w", err)
	}

	nb, err := client.NewOVSDBClient(dbModel, client.WithEndpoint(endpoint))
	if err != nil {
		return nil, fmt.Errorf("failed to create OVN Northbound DB client: %w", err)
	}

	if err := nb.Connect(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to OVN Northbound DB: %w", err)
	}

	// Subscribe to all tables, including Logical_Switch.
	// Teh schema is small enough for now that monitoring everything is fine;
	// we can narrow this down to specific tables later if needed.
	// 
	// We don't need the returned MonitorCookie yet. It would only be needed
	// for finer-grained control, such as canceling a specific monitor or
	// updating it without reconnecting.
	if _, err := nb.MonitorAll(ctx); err != nil {
		return nil, fmt.Errorf("failed to subscribe to OVN Northbound DB tables: %w", err)
	}

	return &Client{nb: nb}, nil
}

// Close disconnects from the OVN Northbound DB.
func (c *Client) Close() {
	c.nb.Disconnect()
}
