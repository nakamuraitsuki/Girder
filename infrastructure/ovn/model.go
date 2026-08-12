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

// LogicalSwitchPort is the model for the OVN_Northbound Logical_Switch_Port table.
//
// ref: https://www.ovn.org/support/dist-docs/ovn-nb.5.html (Logical_Switch_Port TABLE)
type LogicalSwitchPort struct {
	UUID string `ovsdb:"_uuid"`

	// Core Features
	Name string `ovsdb:"name"`
	Type string `ovsdb:"type"`

	// Options
	Options map[string]string `ovsdb:"options"`

	// Containers
	ParentName *string `ovsdb:"parent_name"`
	TagRequest *int    `ovsdb:"tag_request"`
	Tag        *int    `ovsdb:"tag"`

	// Port State
	Up     *bool `ovsdb:"up"`
	Enable *bool `ovsdb:"enable"`

	// Addressing
	Address          []string `ovsdb:"addresses"`
	DynamicAddresses *string  `ovsdb:"dynamic_addresses"`
	PortSecurity     []string `ovsdb:"port_security"`

	// Peer
	Peer *string `ovsdb:"peer"`

	// DHCP
	DHCPv4Options *string `ovsdb:"dhcpv4_options"`
	DHCPv6Options *string `ovsdb:"dhcpv6_options"`

	// Mirror / Health Check
	MirrorRules  []string `ovsdb:"mirror_rules"`
	HealthChecks []string `ovsdb:"health_checks"`

	// HA
	HAChassingGroup *string `ovsdb:"ha_chassis_group"`

	// Common Cloumuns
	ExternalIDs map[string]string `ovsdb:"external_ids"`
}

func DatabaseModel() (model.ClientDBModel, error) {
	return model.NewClientDBModel("OVN_Northbound", map[string]model.Model{
		"Logical_Switch": &LogicalSwitch{},
		"Logical_Switch_Port": &LogicalSwitchPort{},
	})
}
