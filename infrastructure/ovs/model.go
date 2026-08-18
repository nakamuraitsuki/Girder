package ovs

import "github.com/ovn-kubernetes/libovsdb/model"

// OpenVSwitch represents the Open_vSwitch DB root table.
// There must be exactly one record in this table.
// ref: ovs-vswitchd.conf.db(5)
type OpenVSwitch struct {
	UUID            string            `ovsdb:"_uuid"`
	Datapaths       map[string]string `ovsdb:"datapaths"`
	Bridges         []string          `ovsdb:"bridges"`
	SSL             *string           `ovsdb:"ssl"`
	NextCfg         int               `ovsdb:"next_cfg"`
	CurCfg          int               `ovsdb:"cur_cfg"`
	DpdkInitialized bool              `ovsdb:"dpdk_initialized"`
	Statistics      map[string]string `ovsdb:"statistics"`
	OVSVersion      *string           `ovsdb:"ovs_version"`
	DBVersion       *string           `ovsdb:"db_version"`
	SystemType      *string           `ovsdb:"system_type"`
	SystemVersion   *string           `ovsdb:"system_version"`
	DpdkVersion     *string           `ovsdb:"dpdk_version"`
	DatapathTypes   []string          `ovsdb:"datapath_types"`
	IfaceTypes      []string          `ovsdb:"iface_types"`
	ManagerOptions  []string          `ovsdb:"manager_options"`
	OtherConfig     map[string]string `ovsdb:"other_config"`
	ExternalIDs     map[string]string `ovsdb:"external_ids"`
}

// Bridge represents Open_vSwitch DB bridge table.
// ref: ovs-vswitchd.conf.db(5)
type Bridge struct {
	UUID          string            `ovsdb:"_uuid"`
	Name          string            `ovsdb:"name"`
	Ports         []string          `ovsdb:"ports"`
	Mirrors       []string          `ovsdb:"mirrors"`
	NetFlow       *string           `ovsdb:"netflow"`
	SFlow         *string           `ovsdb:"sflow"`
	IPFIX         *string           `ovsdb:"ipfix"`
	FloodVLANs    []int             `ovsdb:"flood_vlans"`
	AutoAttach    *string           `ovsdb:"auto_attach"`
	Controller    []string          `ovsdb:"controller"`
	FlowTables    map[int]string    `ovsdb:"flow_tables"`
	FailMode      *string           `ovsdb:"fail_mode"`
	DatapathID    *string           `ovsdb:"datapath_id"`
	DatapathVersion string          `ovsdb:"datapath_version"`
	Protocols     []string          `ovsdb:"protocols"`
	STPEnable     bool              `ovsdb:"stp_enable"`
	RSTPEnable    bool              `ovsdb:"rstp_enable"`
	McastSnoopingEnable bool        `ovsdb:"mcast_snooping_enable"`
	DatapathType  string            `ovsdb:"datapath_type"`
	OtherConfig   map[string]string `ovsdb:"other_config"`
	ExternalIDs   map[string]string `ovsdb:"external_ids"`
	Status        map[string]string `ovsdb:"status"`
	RSTPStatus    map[string]string `ovsdb:"rstp_status"`
}

// Port represents Open_vSwitch DB port table.
// ref: ovs-vswitchd.conf.db(5)
type Port struct {
	UUID            string            `ovsdb:"_uuid"`
	Name            string            `ovsdb:"name"`
	Interfaces      []string          `ovsdb:"interfaces"`
	VLANMode        *string           `ovsdb:"vlan_mode"`
	Tag             *int              `ovsdb:"tag"`
	Trunks          []int             `ovsdb:"trunks"`
	CVLANs          []int             `ovsdb:"cvlans"`
	BondMode        *string           `ovsdb:"bond_mode"`
	LACP            *string           `ovsdb:"lacp"`
	BondUpdelay     int               `ovsdb:"bond_updelay"`
	BondDowndelay   int               `ovsdb:"bond_downdelay"`
	BondFakeIface   bool              `ovsdb:"bond_fake_iface"`
	OtherConfig     map[string]string `ovsdb:"other_config"`
	ExternalIDs     map[string]string `ovsdb:"external_ids"`
	QOS             *string           `ovsdb:"qos"`
	MAC             *string           `ovsdb:"mac"`
	FakeBridge      bool              `ovsdb:"fake_bridge"`
	Protected       bool              `ovsdb:"protected"`
	BondActiveSlave *string           `ovsdb:"bond_active_slave"`
	Status          map[string]string `ovsdb:"status"`
	Statistics      map[string]int    `ovsdb:"statistics"`
	RSTPStatus      map[string]string `ovsdb:"rstp_status"`
	RSTPStatistics  map[string]int    `ovsdb:"rstp_statistics"`
}

// Interface represents Open_vSwitch DB interface table.
// ref: ovs-vswitchd.conf.db(5)
type Interface struct {
	UUID               string            `ovsdb:"_uuid"`
	Name               string            `ovsdb:"name"`
	IfIndex            *int              `ovsdb:"ifindex"`
	MACInUse           *string           `ovsdb:"mac_in_use"`
	MAC                *string           `ovsdb:"mac"`
	Error              *string           `ovsdb:"error"`
	OFPort             *int              `ovsdb:"ofport"`
	OFPortRequest      *int              `ovsdb:"ofport_request"`
	Type               string            `ovsdb:"type"`
	Options            map[string]string `ovsdb:"options"`
	ExternalIDs        map[string]string `ovsdb:"external_ids"`
	OtherConfig        map[string]string `ovsdb:"other_config"`
	MTU                *int              `ovsdb:"mtu"`
	MTURequest         *int              `ovsdb:"mtu_request"`
	AdminState         *string           `ovsdb:"admin_state"`
	LinkState          *string           `ovsdb:"link_state"`
	LinkResets         *int              `ovsdb:"link_resets"`
	LinkSpeed          *int              `ovsdb:"link_speed"`
	Duplex             *string           `ovsdb:"duplex"`
	LACPCurrent        *bool             `ovsdb:"lacp_current"`
	Status             map[string]string `ovsdb:"status"`
	Statistics         map[string]int    `ovsdb:"statistics"`
	IngressPolicingRate  int             `ovsdb:"ingress_policing_rate"`
	IngressPolicingBurst int             `ovsdb:"ingress_policing_burst"`
	BFD                map[string]string `ovsdb:"bfd"`
	BFDStatus          map[string]string `ovsdb:"bfd_status"`
	CFMMPID            *int              `ovsdb:"cfm_mpid"`
	CFMFlapCount       *int              `ovsdb:"cfm_flap_count"`
	CFMFault           *bool             `ovsdb:"cfm_fault"`
	CFMFaultStatus     []string          `ovsdb:"cfm_fault_status"`
	CFMRemoteOpstate   *string           `ovsdb:"cfm_remote_opstate"`
	CFMHealth          *int              `ovsdb:"cfm_health"`
	CFMRemoteMPIDs     []int             `ovsdb:"cfm_remote_mpids"`
	LLDP               map[string]string `ovsdb:"lldp"`
}

const (
	integrationBridgeName = "br-int"
	externalIDKeyIfaceID  = "iface-id"
	externalIDKeyBridgeMappings = "ovn-bridge-mappings"
)

func DatabaseModel() (model.ClientDBModel, error) {
	return model.NewClientDBModel("Open_vSwitch", map[string]model.Model{
		"Open_vSwitch": &OpenVSwitch{},
		"Bridge": &Bridge{},
		"Port": &Port{},
		"Interface": &Interface{},
	})
}