package api

import "time"

type Metadata map[string]string

type Resource struct {
	ID          string
	Kind        string
	Name        string
	Description string
	Driver      string
	ExternalID  string
	Metadata    Metadata
}

type Relation struct {
	ID       string
	Kind     string
	FromID   string
	ToID     string
	Metadata Metadata
}

type Annotation struct {
	ResourceID string
	Key        string
	Value      string
}

type Timestamped struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}
