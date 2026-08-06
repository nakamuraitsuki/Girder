package api

import (
	"context"
	"time"
)

type DNSZone struct {
	ID         string
	Name       string
	Origin     string
	TTLSeconds uint32
	ExternalID string
	Metadata   Metadata
}

type DNSRecord struct {
	ID         string
	ZoneID     string
	Name       string
	Type       string
	Value      string
	TTLSeconds uint32
	Priority   *uint16
	Metadata   Metadata
}

type DHCPLease struct {
	ID              string
	NetworkID       string
	Address         string
	HardwareAddress string
	ClientID        string
	Hostname        string
	ExpiresAt       time.Time
	Metadata        Metadata
}

type ListDNSZonesRequest struct {
	DriverID string
	All      bool
}

type CreateDNSZoneRequest struct {
	DriverID string
	Zone     DNSZone
}

type DeleteDNSZoneRequest struct {
	DriverID string
	ZoneID   string
}

type ListDNSRecordsRequest struct {
	DriverID string
	ZoneID   string
	All      bool
}

type CreateDNSRecordRequest struct {
	DriverID string
	Record   DNSRecord
}

type DeleteDNSRecordRequest struct {
	DriverID string
	ZoneID   string
	RecordID string
}

type UpdateDNSRecordRequest struct {
	DriverID string
	ZoneID   string
	Record   DNSRecord
}

type ListDHCPLeasesRequest struct {
	DriverID  string
	NetworkID string
	All       bool
}

type CreateDHCPLeaseRequest struct {
	DriverID string
	Lease    DHCPLease
}

type DeleteDHCPLeaseRequest struct {
	DriverID string
	LeaseID  string
}

type DNSAPI interface {
	ListDNSZones(context.Context, ListDNSZonesRequest) ([]DNSZone, error)
	CreateDNSZone(context.Context, CreateDNSZoneRequest) (DNSZone, error)
	DeleteDNSZone(context.Context, DeleteDNSZoneRequest) (string, error)
	ListDNSRecords(context.Context, ListDNSRecordsRequest) ([]DNSRecord, error)
	CreateDNSRecord(context.Context, CreateDNSRecordRequest) (DNSRecord, error)
	UpdateDNSRecord(context.Context, UpdateDNSRecordRequest) (DNSRecord, error)
	DeleteDNSRecord(context.Context, DeleteDNSRecordRequest) (string, error)
	ListDHCPLeases(context.Context, ListDHCPLeasesRequest) ([]DHCPLease, error)
	CreateDHCPLease(context.Context, CreateDHCPLeaseRequest) (DHCPLease, error)
	DeleteDHCPLease(context.Context, DeleteDHCPLeaseRequest) (string, error)
}
