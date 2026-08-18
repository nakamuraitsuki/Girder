package ovs

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/ovn-kubernetes/libovsdb/model"
)

// GetBridgeMappings returns the current ovn-bridgemappings entries as a
// map of network_name to bridge name. Returns anempty map if the key is not yet set.
func (c *Client) GetBridgeMappings(ctx context.Context) (map[string]string, error) {
	ovsRow, err := c.getRootOpenVSwitch(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get root Open_vSwitch row: %w", err)
	}

	return parseBridgeMappings(ovsRow.ExternalIDs[externalIDKeyBridgeMappings]), nil
}

func (c *Client) AddBridgeMapping(ctx context.Context, networkName, bridgeName string) error {
	if networkName == "" {
		return fmt.Errorf("network name is empty")
	}
	if bridgeName == "" {
		return fmt.Errorf("bridge name is empty")
	}

	ovsRow, err := c.getRootOpenVSwitch(ctx)
	if err != nil {
		return fmt.Errorf("failed to get root Open_vSwitch row: %w", err)
	}

	mappings := parseBridgeMappings(ovsRow.ExternalIDs[externalIDKeyBridgeMappings])
	mappings[networkName] = bridgeName
	newValue := formatBridgeMappings(mappings)

	mutateOps, err := c.ovsdb.Where(ovsRow).Mutate(ovsRow,
		model.Mutation{
			Field:   &ovsRow.ExternalIDs,
			Mutator: "insert",
			Value:   map[string]string{externalIDKeyBridgeMappings: newValue},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to build mutation for bridge mappings: %w", err)
	}

	if _, err := c.ovsdb.Transact(ctx, mutateOps...); err != nil {
		return fmt.Errorf("failed to transact: %w", err)
	}

	return nil
}

// Utility

// parseBridgeMappings parses a value like "office:br-ex,isp:br-ex"
// into a map of network_name to bridge name.
func parseBridgeMappings(value string) map[string]string {
	mappings := make(map[string]string)
	if value == "" {
		return mappings
	}

	for _, entry := range strings.Split(value, ",") {
		parts := strings.SplitN(entry, ":", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			continue
		}
		mappings[parts[0]] = parts[1]
	}

	return mappings
}

// formatBridgeMappings formats a map of network_name to bridge name
// back into the "a:b,c:d" format for storing in the OVSDB.
func formatBridgeMappings(mappings map[string]string) string {
	names := make([]string, 0, len(mappings))
	for name := range mappings {
		names = append(names, name)
	}
	sort.Strings(names)

	entries := make([]string, 0, len(names))
	for _, name := range names {
		entries = append(entries, fmt.Sprintf("%s:%s", name, mappings[name]))
	}
	return strings.Join(entries, ",")
}
