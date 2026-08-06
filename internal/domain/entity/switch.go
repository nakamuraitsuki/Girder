package entity

import (
	"fmt"
	"strings"
)

type Switch struct {
	ID          SwitchID
	Name        string
	Description string
	PortIDs     map[PortID]struct{}
	ACLIDs      map[ACLID]struct{}
	State       SwitchState
}

func NewSwitch(name string) (Switch, error) {
	if strings.TrimSpace(name) == "" {
		return Switch{}, fmt.Errorf("%w: switch name", ErrEmptyName)
	}
	return Switch{ID: NewSwitchID(), Name: strings.TrimSpace(name), PortIDs: map[PortID]struct{}{}, ACLIDs: map[ACLID]struct{}{}, State: SwitchStateUp}, nil
}

func (s *Switch) Rename(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("%w: switch name", ErrEmptyName)
	}
	s.Name = strings.TrimSpace(name)
	return nil
}

func (s *Switch) SetDescription(description string) {
	s.Description = strings.TrimSpace(description)
}

func (s *Switch) AttachPort(portID PortID) error {
	if portID == "" {
		return fmt.Errorf("%w: port", ErrInvalidRelation)
	}
	if s.PortIDs == nil {
		s.PortIDs = map[PortID]struct{}{}
	}
	s.PortIDs[portID] = struct{}{}
	return nil
}

func (s *Switch) DetachPort(portID PortID) {
	if s.PortIDs == nil {
		return
	}
	delete(s.PortIDs, portID)
}

func (s *Switch) AttachACL(aclID ACLID) error {
	if aclID == "" {
		return fmt.Errorf("%w: acl", ErrInvalidRelation)
	}
	if s.ACLIDs == nil {
		s.ACLIDs = map[ACLID]struct{}{}
	}
	s.ACLIDs[aclID] = struct{}{}
	return nil
}

func (s *Switch) DetachACL(aclID ACLID) {
	if s.ACLIDs == nil {
		return
	}
	delete(s.ACLIDs, aclID)
}

func (s *Switch) SetState(state SwitchState) {
	s.State = state
}