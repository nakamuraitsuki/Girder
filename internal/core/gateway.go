package core

import (
	"context"
	"fmt"

	"github.com/nakamuraitsuki/Girder/infrastructure/ovn"
	"github.com/nakamuraitsuki/Girder/infrastructure/ovs"
)

func (c *Core) CreateGateway(
	ctx context.Context,
	gwName string,
	physicalNICName string,
) error {
	ovsBridge, err := c.ovs.CreateBridge(ctx, physicalNICName)
	if err != nil {
		return fmt.Errorf("failed to create OVS bridge: %w", err)
	}

	if _, err := c.ovs.AddPort(
		ctx,
		ovsBridge,
		&ovs.Port{
			Name: physicalNICName,
		},
		&ovs.Interface{
			Name: physicalNICName,
		},
	); err != nil {
		return fmt.Errorf("failed to add port to OVS bridge: %w", err)
	}

	sw, err := c.ovn.CreateSwitch(
		ctx,
		&ovn.LogicalSwitch{
			Name: gwName,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create logical switch: %w", err)
	}

	if err := c.ovs.AddBridgeMapping(ctx, sw.Name, ovsBridge.Name); err != nil {
		return fmt.Errorf("failed to add bridge mapping: %w", err)
	}

	// LocalnetPort
	if _, err := c.ovn.CreateLogicalSwitchPort(
		ctx,
		sw,
		&ovn.LogicalSwitchPort{
			Name: sw.Name,
			Type: "localnet",
			Options: map[string]string{
				"network_name": gwName,
			},
		},
	); err != nil {
		return fmt.Errorf("failed to create logical switch port: %w", err)
	}

	return nil
}
