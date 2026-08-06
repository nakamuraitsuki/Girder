package entity

import (
	"errors"
	"sort"
	"fmt"
	"strings"

	"girder/internal/domain/value"
)

var (
	ErrEmptyName            = errors.New("name must not be empty")
	ErrInvalidResource      = errors.New("invalid resource")
	ErrAlreadyAttached      = errors.New("already attached")
	ErrNotAttached          = errors.New("not attached")
	ErrDuplicateItem        = errors.New("duplicate item")
	ErrNotFound             = errors.New("not found")
	ErrInsufficientCapacity = errors.New("insufficient capacity")
	ErrInvalidRelation      = errors.New("invalid relation")
	ErrInvalidState         = errors.New("invalid state")
)

type PortState string

const (
	PortStateDetached PortState = "detached"
	PortStateAttached PortState = "attached"
)

type SwitchState string

const (
	SwitchStateUp   SwitchState = "up"
	SwitchStateDown SwitchState = "down"
)

type RouterState string

const (
	RouterStateUp   RouterState = "up"
	RouterStateDown RouterState = "down"
)

type NATType string

const (
	NATTypeSource NATType = "source"
	NATTypeDest   NATType = "destination"
	NATTypeMasq   NATType = "masquerade"
)

type ACLAction string

const (
	ACLActionAllow ACLAction = "allow"
	ACLActionDeny  ACLAction = "deny"
)

type DNSRecordType string

const (
	DNSRecordTypeA     DNSRecordType = "A"
	DNSRecordTypeAAAA  DNSRecordType = "AAAA"
	DNSRecordTypeCNAME DNSRecordType = "CNAME"
	DNSRecordTypeMX    DNSRecordType = "MX"
	DNSRecordTypeNS    DNSRecordType = "NS"
	DNSRecordTypeTXT   DNSRecordType = "TXT"
)

type NodeCapacity struct {
	CPU       int
	MemoryMiB int64
	DiskGiB   int64
}

type NodeAllocation struct {
	CPU       int
	MemoryMiB int64
	DiskGiB   int64
}

func (c NodeCapacity) CanHost(request NodeAllocation) bool {
	return request.CPU >= 0 && request.MemoryMiB >= 0 && request.DiskGiB >= 0 &&
		request.CPU <= c.CPU && request.MemoryMiB <= c.MemoryMiB && request.DiskGiB <= c.DiskGiB
}

func (c NodeCapacity) Sub(allocation NodeAllocation) NodeCapacity {
	return NodeCapacity{
		CPU:       c.CPU - allocation.CPU,
		MemoryMiB: c.MemoryMiB - allocation.MemoryMiB,
		DiskGiB:   c.DiskGiB - allocation.DiskGiB,
	}
}


type ACLRule struct {
	Priority    int
	Action      ACLAction
	Source      value.CIDR
	Destination value.CIDR
	Protocol    string
	PortRange   string
}

type ACL struct {
	ID            ACLID
	Name          string
	Description   string
	DefaultAction ACLAction
	Rules         []ACLRule
}

func NewACL(name string, defaultAction ACLAction) (ACL, error) {
	if strings.TrimSpace(name) == "" {
		return ACL{}, fmt.Errorf("%w: acl name", ErrEmptyName)
	}
	if defaultAction != ACLActionAllow && defaultAction != ACLActionDeny {
		return ACL{}, fmt.Errorf("%w: acl action", ErrInvalidResource)
	}
	return ACL{ID: NewACLID(), Name: strings.TrimSpace(name), DefaultAction: defaultAction, Rules: []ACLRule{}}, nil
}

func (a *ACL) Rename(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("%w: acl name", ErrEmptyName)
	}
	a.Name = strings.TrimSpace(name)
	return nil
}

func (a *ACL) SetDescription(description string) {
	a.Description = strings.TrimSpace(description)
}

func (a *ACL) SetDefaultAction(action ACLAction) error {
	if action != ACLActionAllow && action != ACLActionDeny {
		return fmt.Errorf("%w: acl action", ErrInvalidResource)
	}
	a.DefaultAction = action
	return nil
}

func (a *ACL) AddRule(rule ACLRule) error {
	if rule.Priority < 0 {
		return fmt.Errorf("%w: acl rule", ErrInvalidResource)
	}
	for _, existing := range a.Rules {
		if existing.Priority == rule.Priority {
			return fmt.Errorf("%w: acl rule priority", ErrDuplicateItem)
		}
	}
	a.Rules = append(a.Rules, rule)
	sort.Slice(a.Rules, func(i, j int) bool { return a.Rules[i].Priority < a.Rules[j].Priority })
	return nil
}

func (a *ACL) DeleteRule(priority int) {
	filtered := a.Rules[:0]
	for _, rule := range a.Rules {
		if rule.Priority != priority {
			filtered = append(filtered, rule)
		}
	}
	a.Rules = filtered
}