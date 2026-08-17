package ovs

import (
	"context"
	"fmt"
)

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
