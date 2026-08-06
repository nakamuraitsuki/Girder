package entity

import (
	"fmt"
	"strings"
)

type DNSZone struct {
	ID          DNSZoneID
	Name        string
	Origin      string
	Description string
	TTLSeconds  uint32
	RecordIDs   map[DNSRecordID]struct{}
	Enabled     bool
}

func NewDNSZone(name, origin string, ttlSeconds uint32) (DNSZone, error) {
	if strings.TrimSpace(name) == "" {
		return DNSZone{}, fmt.Errorf("%w: dns zone name", ErrEmptyName)
	}
	if strings.TrimSpace(origin) == "" {
		return DNSZone{}, fmt.Errorf("%w: dns zone origin", ErrEmptyName)
	}
	return DNSZone{ID: NewDNSZoneID(), Name: strings.TrimSpace(name), Origin: strings.TrimSpace(origin), TTLSeconds: ttlSeconds, RecordIDs: map[DNSRecordID]struct{}{}, Enabled: true}, nil
}

func (z *DNSZone) Rename(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("%w: dns zone name", ErrEmptyName)
	}
	z.Name = strings.TrimSpace(name)
	return nil
}

func (z *DNSZone) SetDescription(description string) {
	z.Description = strings.TrimSpace(description)
}

func (z *DNSZone) SetTTL(ttlSeconds uint32) {
	z.TTLSeconds = ttlSeconds
}

func (z *DNSZone) AddRecord(recordID DNSRecordID) error {
	if recordID == "" {
		return fmt.Errorf("%w: dns record", ErrInvalidRelation)
	}
	if z.RecordIDs == nil {
		z.RecordIDs = map[DNSRecordID]struct{}{}
	}
	z.RecordIDs[recordID] = struct{}{}
	return nil
}

func (z *DNSZone) RemoveRecord(recordID DNSRecordID) {
	if z.RecordIDs == nil {
		return
	}
	delete(z.RecordIDs, recordID)
}

type DNSRecord struct {
	ID         DNSRecordID
	ZoneID     DNSZoneID
	Name       string
	Type       DNSRecordType
	Value      string
	TTLSeconds uint32
	Priority   *uint16
}

func NewDNSRecord(zoneID DNSZoneID, name string, recordType DNSRecordType, valueText string, ttlSeconds uint32) (DNSRecord, error) {
	if zoneID == "" {
		return DNSRecord{}, fmt.Errorf("%w: dns record zone", ErrInvalidRelation)
	}
	if strings.TrimSpace(name) == "" {
		return DNSRecord{}, fmt.Errorf("%w: dns record name", ErrEmptyName)
	}
	if strings.TrimSpace(valueText) == "" {
		return DNSRecord{}, fmt.Errorf("%w: dns record value", ErrEmptyName)
	}
	return DNSRecord{ID: NewDNSRecordID(), ZoneID: zoneID, Name: strings.TrimSpace(name), Type: recordType, Value: strings.TrimSpace(valueText), TTLSeconds: ttlSeconds}, nil
}

func (r *DNSRecord) Rename(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("%w: dns record name", ErrEmptyName)
	}
	r.Name = strings.TrimSpace(name)
	return nil
}

func (r *DNSRecord) UpdateValue(valueText string) error {
	if strings.TrimSpace(valueText) == "" {
		return fmt.Errorf("%w: dns record value", ErrEmptyName)
	}
	r.Value = strings.TrimSpace(valueText)
	return nil
}

func (r *DNSRecord) UpdateTTL(ttlSeconds uint32) {
	r.TTLSeconds = ttlSeconds
}

func (r *DNSRecord) UpdatePriority(priority *uint16) {
	r.Priority = priority
}
