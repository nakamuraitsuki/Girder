package ovn

import "github.com/ovn-kubernetes/libovsdb/model"

// LogicalSwitch is the model for the OVN_Northbound Logical_Switch table.
//
// ref: https://www.ovn.org/support/dist-docs/ovn-nb.5.html (Logical_Switch TABLE)
type LogicalSwitch struct {
	UUID string `ovsdb:"_uuid"`

	// Naming
	Name string `ovsdb:"name"`

	// Relational Entities
	Ports             []string `ovsdb:"ports"`
	LoadBalancer      []string `ovsdb:"load_balancer"`
	LoadBalancerGroup []string `ovsdb:"load_balancer_group"`
	ACLs              []string `ovsdb:"acls"`
	QOSRules          []string `ovsdb:"qos_rules"`
	DNSRecords        []string `ovsdb:"dns_records"`
	ForwardingGroups  []string `ovsdb:"forwarding_groups"`

	// IPAM / Multi-Cast / InterConnect
	OtherConfig map[string]string `ovsdb:"other_config"`

	// Control Plane Protection
	Copp *string `ovsdb:"copp"`

	// Common Columns
	ExternalIDs map[string]string `ovsdb:"external_ids"`
}

func DatabaseModel() (model.ClientDBModel, error) {
	return model.NewClientDBModel("OVN_Northbound", map[string]model.Model{
		"Logical_Switch": &LogicalSwitch{},
	})
}
