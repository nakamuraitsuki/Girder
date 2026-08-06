package entity

import "girder/internal/domain/value"

type NodeID string
type ImageID string
type VMID string
type VolumeID string
type SnapshotID string
type SwitchID string
type RouterID string
type PortID string
type RouteID string
type ACLID string
type NATID string
type DHCPID string
type DNSZoneID string
type DNSRecordID string
type BlueprintID string
type UserID string
type ProjectID string

func NewNodeID() NodeID           { return NodeID(value.NewPrefixedULID("node")) }
func NewImageID() ImageID         { return ImageID(value.NewPrefixedULID("img")) }
func NewVMID() VMID               { return VMID(value.NewPrefixedULID("vm")) }
func NewVolumeID() VolumeID       { return VolumeID(value.NewPrefixedULID("vol")) }
func NewSnapshotID() SnapshotID   { return SnapshotID(value.NewPrefixedULID("snap")) }
func NewSwitchID() SwitchID       { return SwitchID(value.NewPrefixedULID("sw")) }
func NewRouterID() RouterID       { return RouterID(value.NewPrefixedULID("rt")) }
func NewPortID() PortID           { return PortID(value.NewPrefixedULID("port")) }
func NewRouteID() RouteID         { return RouteID(value.NewPrefixedULID("route")) }
func NewACLID() ACLID             { return ACLID(value.NewPrefixedULID("acl")) }
func NewNATID() NATID             { return NATID(value.NewPrefixedULID("nat")) }
func NewDHCPID() DHCPID           { return DHCPID(value.NewPrefixedULID("dhcp")) }
func NewDNSZoneID() DNSZoneID     { return DNSZoneID(value.NewPrefixedULID("dnszone")) }
func NewDNSRecordID() DNSRecordID { return DNSRecordID(value.NewPrefixedULID("dnsrec")) }
func NewBlueprintID() BlueprintID { return BlueprintID(value.NewPrefixedULID("bp")) }
func NewUserID() UserID           { return UserID(value.NewPrefixedULID("user")) }
func NewProjectID() ProjectID     { return ProjectID(value.NewPrefixedULID("proj")) }
