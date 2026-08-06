package entity

import (
	"fmt"
	"strings"
)

type User struct {
	ID          UserID
	Name        string
	Email       string
	Description string
	RoleIDs     map[string]struct{}
	ProjectIDs  map[ProjectID]struct{}
	Enabled     bool
}

func NewUser(name, email string) (User, error) {
	if strings.TrimSpace(name) == "" {
		return User{}, fmt.Errorf("%w: user name", ErrEmptyName)
	}
	if strings.TrimSpace(email) == "" {
		return User{}, fmt.Errorf("%w: user email", ErrEmptyName)
	}
	return User{ID: NewUserID(), Name: strings.TrimSpace(name), Email: strings.TrimSpace(email), RoleIDs: map[string]struct{}{}, ProjectIDs: map[ProjectID]struct{}{}, Enabled: true}, nil
}

func (u *User) Rename(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("%w: user name", ErrEmptyName)
	}
	u.Name = strings.TrimSpace(name)
	return nil
}

func (u *User) SetEmail(email string) error {
	if strings.TrimSpace(email) == "" {
		return fmt.Errorf("%w: user email", ErrEmptyName)
	}
	u.Email = strings.TrimSpace(email)
	return nil
}

func (u *User) SetDescription(description string) {
	u.Description = strings.TrimSpace(description)
}

func (u *User) AddRole(role string) error {
	role = strings.TrimSpace(role)
	if role == "" {
		return fmt.Errorf("%w: user role", ErrEmptyName)
	}
	if u.RoleIDs == nil {
		u.RoleIDs = map[string]struct{}{}
	}
	u.RoleIDs[role] = struct{}{}
	return nil
}

func (u *User) RemoveRole(role string) {
	if u.RoleIDs == nil {
		return
	}
	delete(u.RoleIDs, strings.TrimSpace(role))
}

func (u *User) JoinProject(projectID ProjectID) error {
	if projectID == "" {
		return fmt.Errorf("%w: project", ErrInvalidRelation)
	}
	if u.ProjectIDs == nil {
		u.ProjectIDs = map[ProjectID]struct{}{}
	}
	u.ProjectIDs[projectID] = struct{}{}
	return nil
}

func (u *User) LeaveProject(projectID ProjectID) {
	if u.ProjectIDs == nil {
		return
	}
	delete(u.ProjectIDs, projectID)
}

func (u *User) Enable()  { u.Enabled = true }
func (u *User) Disable() { u.Enabled = false }

type Project struct {
	ID           ProjectID
	Name         string
	Description  string
	OwnerUserID  UserID
	MemberIDs    map[UserID]struct{}
	BlueprintIDs map[BlueprintID]struct{}
	Enabled      bool
}

func NewProject(name string, ownerUserID UserID) (Project, error) {
	if strings.TrimSpace(name) == "" {
		return Project{}, fmt.Errorf("%w: project name", ErrEmptyName)
	}
	if ownerUserID == "" {
		return Project{}, fmt.Errorf("%w: project owner", ErrInvalidRelation)
	}
	return Project{ID: NewProjectID(), Name: strings.TrimSpace(name), OwnerUserID: ownerUserID, MemberIDs: map[UserID]struct{}{}, BlueprintIDs: map[BlueprintID]struct{}{}, Enabled: true}, nil
}

func (p *Project) Rename(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("%w: project name", ErrEmptyName)
	}
	p.Name = strings.TrimSpace(name)
	return nil
}

func (p *Project) SetDescription(description string) {
	p.Description = strings.TrimSpace(description)
}

func (p *Project) AddMember(userID UserID) error {
	if userID == "" {
		return fmt.Errorf("%w: user", ErrInvalidRelation)
	}
	if p.MemberIDs == nil {
		p.MemberIDs = map[UserID]struct{}{}
	}
	p.MemberIDs[userID] = struct{}{}
	return nil
}

func (p *Project) RemoveMember(userID UserID) {
	if p.MemberIDs == nil {
		return
	}
	delete(p.MemberIDs, userID)
}

func (p *Project) AddBlueprint(blueprintID BlueprintID) error {
	if blueprintID == "" {
		return fmt.Errorf("%w: blueprint", ErrInvalidRelation)
	}
	if p.BlueprintIDs == nil {
		p.BlueprintIDs = map[BlueprintID]struct{}{}
	}
	p.BlueprintIDs[blueprintID] = struct{}{}
	return nil
}

func (p *Project) RemoveBlueprint(blueprintID BlueprintID) {
	if p.BlueprintIDs == nil {
		return
	}
	delete(p.BlueprintIDs, blueprintID)
}

func (p *Project) Enable()  { p.Enabled = true }
func (p *Project) Disable() { p.Enabled = false }

type Blueprint struct {
	ID           BlueprintID
	Name         string
	Description  string
	ProjectID    ProjectID
	OwnerUserID  UserID
	NodeIDs      map[NodeID]struct{}
	ImageIDs     map[ImageID]struct{}
	VMIDs        map[VMID]struct{}
	VolumeIDs    map[VolumeID]struct{}
	SnapshotIDs  map[SnapshotID]struct{}
	SwitchIDs    map[SwitchID]struct{}
	RouterIDs    map[RouterID]struct{}
	PortIDs      map[PortID]struct{}
	RouteIDs     map[RouteID]struct{}
	ACLIDs       map[ACLID]struct{}
	NATIDs       map[NATID]struct{}
	DHCPIDs      map[DHCPID]struct{}
	DNSZoneIDs   map[DNSZoneID]struct{}
	DNSRecordIDs map[DNSRecordID]struct{}
}

func NewBlueprint(name string, projectID ProjectID, ownerUserID UserID) (Blueprint, error) {
	if strings.TrimSpace(name) == "" {
		return Blueprint{}, fmt.Errorf("%w: blueprint name", ErrEmptyName)
	}
	if projectID == "" || ownerUserID == "" {
		return Blueprint{}, fmt.Errorf("%w: blueprint relation", ErrInvalidRelation)
	}
	return Blueprint{ID: NewBlueprintID(), Name: strings.TrimSpace(name), ProjectID: projectID, OwnerUserID: ownerUserID, NodeIDs: map[NodeID]struct{}{}, ImageIDs: map[ImageID]struct{}{}, VMIDs: map[VMID]struct{}{}, VolumeIDs: map[VolumeID]struct{}{}, SnapshotIDs: map[SnapshotID]struct{}{}, SwitchIDs: map[SwitchID]struct{}{}, RouterIDs: map[RouterID]struct{}{}, PortIDs: map[PortID]struct{}{}, RouteIDs: map[RouteID]struct{}{}, ACLIDs: map[ACLID]struct{}{}, NATIDs: map[NATID]struct{}{}, DHCPIDs: map[DHCPID]struct{}{}, DNSZoneIDs: map[DNSZoneID]struct{}{}, DNSRecordIDs: map[DNSRecordID]struct{}{}}, nil
}

func (b *Blueprint) Rename(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("%w: blueprint name", ErrEmptyName)
	}
	b.Name = strings.TrimSpace(name)
	return nil
}

func (b *Blueprint) SetDescription(description string) {
	b.Description = strings.TrimSpace(description)
}

func (b *Blueprint) AddNode(id NodeID) error   { return addBlueprintIDSet(&b.NodeIDs, id, "node") }
func (b *Blueprint) RemoveNode(id NodeID)      { delete(b.NodeIDs, id) }
func (b *Blueprint) AddImage(id ImageID) error { return addBlueprintIDSet(&b.ImageIDs, id, "image") }
func (b *Blueprint) RemoveImage(id ImageID)    { delete(b.ImageIDs, id) }
func (b *Blueprint) AddVM(id VMID) error       { return addBlueprintIDSet(&b.VMIDs, id, "vm") }
func (b *Blueprint) RemoveVM(id VMID)          { delete(b.VMIDs, id) }
func (b *Blueprint) AddVolume(id VolumeID) error {
	return addBlueprintIDSet(&b.VolumeIDs, id, "volume")
}
func (b *Blueprint) RemoveVolume(id VolumeID) { delete(b.VolumeIDs, id) }
func (b *Blueprint) AddSnapshot(id SnapshotID) error {
	return addBlueprintIDSet(&b.SnapshotIDs, id, "snapshot")
}
func (b *Blueprint) RemoveSnapshot(id SnapshotID) { delete(b.SnapshotIDs, id) }
func (b *Blueprint) AddSwitch(id SwitchID) error {
	return addBlueprintIDSet(&b.SwitchIDs, id, "switch")
}
func (b *Blueprint) RemoveSwitch(id SwitchID) { delete(b.SwitchIDs, id) }
func (b *Blueprint) AddRouter(id RouterID) error {
	return addBlueprintIDSet(&b.RouterIDs, id, "router")
}
func (b *Blueprint) RemoveRouter(id RouterID)  { delete(b.RouterIDs, id) }
func (b *Blueprint) AddPort(id PortID) error   { return addBlueprintIDSet(&b.PortIDs, id, "port") }
func (b *Blueprint) RemovePort(id PortID)      { delete(b.PortIDs, id) }
func (b *Blueprint) AddRoute(id RouteID) error { return addBlueprintIDSet(&b.RouteIDs, id, "route") }
func (b *Blueprint) RemoveRoute(id RouteID)    { delete(b.RouteIDs, id) }
func (b *Blueprint) AddACL(id ACLID) error     { return addBlueprintIDSet(&b.ACLIDs, id, "acl") }
func (b *Blueprint) RemoveACL(id ACLID)        { delete(b.ACLIDs, id) }
func (b *Blueprint) AddNAT(id NATID) error     { return addBlueprintIDSet(&b.NATIDs, id, "nat") }
func (b *Blueprint) RemoveNAT(id NATID)        { delete(b.NATIDs, id) }
func (b *Blueprint) AddDHCP(id DHCPID) error   { return addBlueprintIDSet(&b.DHCPIDs, id, "dhcp") }
func (b *Blueprint) RemoveDHCP(id DHCPID)      { delete(b.DHCPIDs, id) }
func (b *Blueprint) AddDNSZone(id DNSZoneID) error {
	return addBlueprintIDSet(&b.DNSZoneIDs, id, "dns zone")
}
func (b *Blueprint) RemoveDNSZone(id DNSZoneID) { delete(b.DNSZoneIDs, id) }
func (b *Blueprint) AddDNSRecord(id DNSRecordID) error {
	return addBlueprintIDSet(&b.DNSRecordIDs, id, "dns record")
}
func (b *Blueprint) RemoveDNSRecord(id DNSRecordID) { delete(b.DNSRecordIDs, id) }

func addBlueprintIDSet[T comparable](set *map[T]struct{}, id T, label string) error {
	var zero T
	if id == zero {
		return fmt.Errorf("%w: %s", ErrInvalidRelation, label)
	}
	if *set == nil {
		*set = map[T]struct{}{}
	}
	(*set)[id] = struct{}{}
	return nil
}
