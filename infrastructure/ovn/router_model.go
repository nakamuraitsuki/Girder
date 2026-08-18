package ovn

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

	Ipv6Prefix           []string `ovsdb:"ipv6_prefix"`

	Peer *string `ovsdb:"peer"`

	Options map[string]string `ovsdb:"options"`

	Status map[string]string `ovsdb:"status"`

	ExternalIDs map[string]string `ovsdb:"external_ids"`
}

// LogicalRouterStaticRoute is the model for the OVN_Northbound
// Logical_Router_Static_Route table.
//
// ref: https://www.ovn.org/support/dist-docs/ovn-nb.5.html (Logical_Router_Static_Route TABLE)
type LogicalRouterStaticRoute struct {
	UUID string `ovsdb:"_uuid"`

	IPPrefix string  `ovsdb:"ip_prefix"`
	Nexthop  string  `ovsdb:"nexthop"`
	Policy   *string `ovsdb:"policy"`

	OutputPort *string `ovsdb:"output_port"`

	ExternalIDs map[string]string `ovsdb:"external_ids"`
}