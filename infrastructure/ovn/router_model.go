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

	Ipv6RaConfigs        []string `ovsdb:"ipv6_ra_configs"`
	Ipv6Prefix           []string `ovsdb:"ipv6_prefix"`
	Ipv6PrefixDelegation bool     `ovsdb:"ipv6_prefix_delegation"`

	Peer *string `ovsdb:"peer"`

	Options map[string]string `ovsdb:"options"`

	Status map[string]string `ovsdb:"status"`

	ExternalIDs map[string]string `ovsdb:"external_ids"`
}