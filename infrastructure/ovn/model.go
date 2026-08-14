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
	Up      *bool `ovsdb:"up"`
	Enabled *bool `ovsdb:"enabled"`

	// Addressing
	Addresses        []string `ovsdb:"addresses"`
	DynamicAddresses *string  `ovsdb:"dynamic_addresses"`
	PortSecurity     []string `ovsdb:"port_security"`

	// DHCP
	DHCPv4Options *string `ovsdb:"dhcpv4_options"`
	DHCPv6Options *string `ovsdb:"dhcpv6_options"`

	// Mirror / Health Check
	MirrorRules []string `ovsdb:"mirror_rules"`

	// HA
	HAChassisGroup *string `ovsdb:"ha_chassis_group"`

	// Common Columns
	ExternalIDs map[string]string `ovsdb:"external_ids"`
}

type LogicalRouter struct {
	UUID string `ovsdb:"_uuid"`

	// Naming
	Name string `ovsdb:"name"`

	Ports             []string `ovsdb:"ports"`
	StaticRoutes      []string `ovsdb:"static_routes"`
	Policies          []string `ovsdb:"policies"`
	Nat               []string `ovsdb:"nat"`
	LoadBalancer      []string `ovsdb:"load_balancer"`
	LoadBalancerGroup []string `ovsdb:"load_balancer_group"`

	Copp *string `ovsdb:"copp"`

	Enabled *bool `ovsdb:"enabled"`

	Options map[string]string `ovsdb:"options"`

	ExternalIDs map[string]string `ovsdb:"external_ids"`
}

type LogicalRouterPort struct {
	UUID string `ovsdb:"_uuid"`

	// Naming
	Name string `ovsdb:"name"`

	GatewayChassis []string `ovsdb:"gateway_chassis"`
	HaChassisGroup *string  `ovsdb:"ha_chassis_group"`

	Networks []string `ovsdb:"networks"`

	MAC string `ovsdb:"mac"`

	Enabled *bool `ovsdb:"enabled"`

	Ipv6RaConfigs        []string `ovsdb:"ipv6_ra_configs"`
	Ipv6Prefix           []string `ovsdb:"ipv6_prefix"`
	Ipv6PrefixDelegation bool     `ovsdb:"ipv6_prefix_delegation"`

	Peer *string `ovsdb:"peer"`

	Options map[string]string `ovsdb:"options"`

	Status map[string]string `ovsdb:"status"`

	ExternalIDs map[string]string `ovsdb:"external_ids"`
}

func DatabaseModel() (model.ClientDBModel, error) {
	return model.NewClientDBModel("OVN_Northbound", map[string]model.Model{
		"Logical_Switch":      &LogicalSwitch{},
		"Logical_Switch_Port": &LogicalSwitchPort{},
		"Logical_Router":      &LogicalRouter{},
		"Logical_Router_Port": &LogicalRouterPort{},
	})
}
