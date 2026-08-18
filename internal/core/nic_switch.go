package core

import (
	"context"
	"errors"
	"fmt"

	"github.com/nakamuraitsuki/Girder/infrastructure/ovn"
	"github.com/nakamuraitsuki/Girder/infrastructure/ovs"
	"libvirt.org/go/libvirtxml"
)

// ConnectNICtoSwitch connects a libvirt NIC (Tap Device) to an OVN switch
func (c *Core) ConnectNICtoSwitch(ctx context.Context, vmName, nicLogicalName, switchName string) error {
	// Check if the NIC exists in the VM
	iface, err := c.findNICByLogicalName(vmName, nicLogicalName)
	if err != nil {
		return err
	}
	if iface.Target == nil || iface.Target.Dev == "" {
		return fmt.Errorf("nic %q on vm %q does not have a target device", nicLogicalName, vmName)
	}
	tapDevice := iface.Target.Dev

	// Check if the OVN switch exists
	sw, err := c.ovn.GetSwitch(ctx, switchName)
	if err != nil {
		return fmt.Errorf("failed to get switch %q: %w", switchName, err)
	}

	// OVN: creates Logical_Switch_Port(type="")
	lsp, err := c.ovn.CreateLogicalSwitchPort(ctx, sw, &ovn.LogicalSwitchPort{
		Name: nicLogicalName,
		Type: "",
		Addresses: []string{iface.MAC.Address},
	})
	if err != nil {
		return fmt.Errorf("failed to create logical switch port %q on switch %q: %w", nicLogicalName, switchName, err)
	}

	// OVS: creates Port/Interface on br-int and binds iface-id
	bridge, err := c.ovs.GetIntegrationBridge(ctx)
	if err != nil {
		// Rollback OVN logical switch port creation
		if delErr := c.ovn.DeleteLogicalSwitchPort(ctx, sw, lsp.Name); delErr != nil {
			return fmt.Errorf("failed to get integration bridge and failed to delete logical switch port %q: %v, %v", lsp.Name, err, delErr)
		}
		return fmt.Errorf("failed to get integration bridge (ovs setup incomplete?) and rolled back logical switch port %q: %w", lsp.Name, err)
	}
	
	if _, err := c.ovs.AddPort(ctx,
		bridge,
		&ovs.Port{Name: tapDevice},
		&ovs.Interface{
			Name: tapDevice,
			ExternalIDs: map[string]string{
				"iface-id": lsp.Name,
			},
		},
	); err != nil {
		if delErr := c.ovn.DeleteLogicalSwitchPort(ctx, sw, lsp.Name); delErr != nil {
			return fmt.Errorf("failed to add port %q to OVS bridge and failed to delete logical switch port %q: %v, %v", tapDevice, lsp.Name, err, delErr)
		}
		return fmt.Errorf("failed to bind ovs port for nic %q (logical switch port %q rolled back): %w", nicLogicalName, lsp.Name, err)
	}

	return nil
}

// DisconnectNICfromSwitch disconnects a libvirt NIC (Tap Device) from an OVN switch
func (c *Core) DisconnectNICfromSwitch(ctx context.Context, vmName, nicLogicalName, switchName string) error {
	iface, err := c.findNICByLogicalName(vmName, nicLogicalName)
	if err != nil {
		return err
	}
	if iface.Target == nil || iface.Target.Dev == "" {
		return fmt.Errorf("nic %q on vm %q has no tap device", nicLogicalName, vmName)
	}
	tapDevice := iface.Target.Dev

	sw, err := c.ovn.GetSwitch(ctx, switchName)
	if err != nil {
		return fmt.Errorf("failed to get switch %q: %w", switchName, err)
	}

	if _, err := c.ovn.GetLogicalSwitchPort(ctx, nicLogicalName); err != nil {
		if errors.Is(err, ovn.ErrLogicalSwitchPortNotFound) {
			return fmt.Errorf("nic %q is not connected to switch %q", nicLogicalName, switchName)
		}
		return fmt.Errorf("failed to get logical switch port %q: %w", nicLogicalName, err)
	}

	// First, remove the OVS port to avoid leaving a dangling port if OVN deletion fails
	if err := c.ovs.RemovePort(ctx, tapDevice); err != nil {
		return fmt.Errorf("failed to remove ovs port for nic %q: %w", nicLogicalName, err)
	}

	// Then, remove the logical switch port from OVN
	if err := c.ovn.DeleteLogicalSwitchPort(ctx, sw, nicLogicalName); err != nil {
		return fmt.Errorf("ovs port removed but failed to delete logical switch port %q (inconsistent state, manual check required): %w", nicLogicalName, err)
	}

	return nil
}

// findNICByLogicalName searches for a NIC in the specified VM by its alias (logical name).
// The libvirt package does not have a dedicated type for NICs, so we use libvirtxml.DomainInterface directly (Principle 2: Use driver types directly).
func (c *Core) findNICByLogicalName(vmName, logicalName string) (*libvirtxml.DomainInterface, error) {
	ifaces, err := c.libvirt.ListNICs(vmName)
	if err != nil {
		return nil, fmt.Errorf("failed to list NICs for vm %q: %w", vmName, err)
	}

	for i := range ifaces {
		if ifaces[i].Alias != nil && ifaces[i].Alias.Name == logicalName {
			return &ifaces[i], nil
		}
	}

	return nil, fmt.Errorf("nic %q not found on vm %q", logicalName, vmName)
}
