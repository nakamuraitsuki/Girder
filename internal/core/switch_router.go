package core

import (
	"context"
	"errors"
	"fmt"

	"github.com/nakamuraitsuki/Girder/infrastructure/ovn"
)

// ConnectSwitchToRouter connects an OVN logical switch to an OVN logical router.
//
// Naming convention (mirrors the switch-port-name = counterpart-name pattern
// used by ConnectNICtoSwitch):
//   - Logical_Router_Port name = switchName (the router's "view" of the connection)
//   - Logical_Switch_Port name = routerName (the switch's "view" of the connection)
func (c *Core) ConnectSwitchToRouter(ctx context.Context, switchName, routerName string) error {
	sw, err := c.ovn.GetSwitch(ctx, switchName)
	if err != nil {
		return fmt.Errorf("failed to get switch %q: %w", switchName, err)
	}

	router, err := c.ovn.GetRouter(ctx, routerName)
	if err != nil {
		return fmt.Errorf("failed to get logical router %q: %w", routerName, err)
	}

	// OVN: creates Logical_Router_Port on the router side.
	// MAC and networks are left empty here; IP assignment is handled
	// by a separate operation (out of scope for this connect operation).
	lrp, err := c.ovn.CreateLogicalRouterPort(ctx, router, &ovn.LogicalRouterPort{
		Name: switchName,
	})
	if err != nil {
		return fmt.Errorf("failed to create logical router port %q on router %q: %w", switchName, routerName, err)
	}

	// OVN: creates Logical_Switch_Port(type="router") on the switch side,
	// bound to the router port via options["router-port"].
	if _, err := c.ovn.CreateLogicalSwitchPort(ctx, sw, &ovn.LogicalSwitchPort{
		Name: routerName,
		Type: "router",
		Options: map[string]string{
			"router-port": lrp.Name,
		},
	}); err != nil {
		if delErr := c.ovn.DeleteLogicalRouterPort(ctx, router, lrp.Name); delErr != nil {
			return fmt.Errorf("failed to create logical switch port %q on switch %q and failed to delete logical router port %q: %v, %v", routerName, switchName, lrp.Name, err, delErr)
		}
		return fmt.Errorf("failed to create switch-side port for router %q (logical router port %q rolled back): %w", routerName, lrp.Name, err)
	}

	return nil
}
// DisconnectSwitchFromRouter disconnects an OVN logical switch from an OVN logical router.
func (c *Core) DisconnectSwitchFromRouter(ctx context.Context, switchName, routerName string) error {
	sw, err := c.ovn.GetSwitch(ctx, switchName)
	if err != nil {
		return fmt.Errorf("failed to get switch %q: %w", switchName, err)
	}

	router, err := c.ovn.GetRouter(ctx, routerName)
	if err != nil {
		return fmt.Errorf("failed to get logical router %q: %w", routerName, err)
	}

	if _, err := c.ovn.GetLogicalRouterPort(ctx, router, switchName); err != nil {
		if errors.Is(err, ovn.ErrLogicalRouterPortNotFound) {
			return fmt.Errorf("switch %q is not connected to router %q", switchName, routerName)
		}
		return fmt.Errorf("failed to get logical router port %q: %w", switchName, err)
	}

	// First, remove the switch-side port to avoid leaving a dangling
	// options["router-port"] reference if router port deletion fails.
	if err := c.ovn.DeleteLogicalSwitchPort(ctx, sw, routerName); err != nil {
		return fmt.Errorf("failed to delete logical switch port %q on switch %q: %w", routerName, switchName, err)
	}

	// Then, remove the router-side port.
	if err := c.ovn.DeleteLogicalRouterPort(ctx, router, switchName); err != nil {
		return fmt.Errorf("switch port removed but failed to delete logical router port %q (inconsistent state, manual check required): %w", switchName, err)
	}

	return nil
}