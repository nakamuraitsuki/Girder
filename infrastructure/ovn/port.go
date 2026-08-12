package ovn

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

var ErrLogicalSwitchPortAlreadyExists = fmt.Errorf("logical switch port already exists")
var ErrLogicalSwitchPortNotFound = fmt.Errorf("logical switch port not found")

// CreateLogicalSwitchPort creates a new OVN Logical_Switch_Port.
func (c *Client) CreateLogicalSwitchPort(
	ctx context.Context,
	port *LogicalSwitchPort,
) (*LogicalSwitchPort, error) {
	// TODO: validate / resolve values that must not be supplied by caller.

	// TODO: check whether the port already exists.

	// Values controlled by Girder.
	port.UUID = uuid.NewString()

	// TODO: build and transact create operation.

	return port, nil
}

// GetLogicalSwitchPort returns the Logical_Switch_Port identified by name.
func (c *Client) GetLogicalSwitchPort(
	ctx context.Context,
	name string,
) (*LogicalSwitchPort, error) {
	// TODO: list Logical_Switch_Port records.

	// TODO: handle not found / ambiguous name.

	return nil, nil
}

// DeleteLogicalSwitchPort deletes the Logical_Switch_Port identified by name.
func (c *Client) DeleteLogicalSwitchPort(
	ctx context.Context,
	name string,
) error {
	// TODO: get the target port.

	// TODO: build delete operation.

	// TODO: transact delete operation.

	return nil
}
