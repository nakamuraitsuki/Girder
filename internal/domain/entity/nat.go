package entity

import (
	"fmt"
	"strings"
	
	"girder/internal/domain/value"
)

type NAT struct {
	ID          NATID
	Name        string
	Description string
	Type        NATType
	Source      value.CIDR
	Translation value.IPv4
	Enabled     bool
}

func NewNAT(name string, natType NATType, source value.CIDR, translation value.IPv4) (NAT, error) {
	if strings.TrimSpace(name) == "" {
		return NAT{}, fmt.Errorf("%w: nat name", ErrEmptyName)
	}
	if source.IsZero() || translation.IsZero() {
		return NAT{}, fmt.Errorf("%w: nat target", ErrInvalidResource)
	}
	if natType != NATTypeSource && natType != NATTypeDest && natType != NATTypeMasq {
		return NAT{}, fmt.Errorf("%w: nat type", ErrInvalidResource)
	}
	return NAT{ID: NewNATID(), Name: strings.TrimSpace(name), Type: natType, Source: source, Translation: translation, Enabled: true}, nil
}

func (n *NAT) Rename(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("%w: nat name", ErrEmptyName)
	}
	n.Name = strings.TrimSpace(name)
	return nil
}

func (n *NAT) SetDescription(description string) {
	n.Description = strings.TrimSpace(description)
}

func (n *NAT) Enable()  { n.Enabled = true }
func (n *NAT) Disable() { n.Enabled = false }

func (n *NAT) SetTranslation(address value.IPv4) error {
	if address.IsZero() {
		return fmt.Errorf("%w: nat translation", ErrInvalidResource)
	}
	n.Translation = address
	return nil
}