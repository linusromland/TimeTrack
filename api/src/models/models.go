package models

import (
	"time"
)

type User struct {
	ID        string    `bson:"_id" json:"id"`
	Email     string    `bson:"email" json:"email"`
	Password  string    `bson:"password" json:"password,omitempty"`
	DeletedAt time.Time `bson:"deleted_at,omitempty" json:"-"`
}

type IntegrationType string

const (
	IntegrationJira IntegrationType = "jira"
)

type IntegrationInfo struct {
	Type       IntegrationType `bson:"type" json:"type"`               // e.g. "jira"
	Key        string          `bson:"key" json:"key"`                 // e.g. "MNT-123"
	ExternalID string          `bson:"external_id" json:"external_id"` // ID in the integration system
}

type Project struct {
	ID          string          `bson:"_id" json:"id"`
	Name        string          `bson:"name" json:"name"`
	Integration IntegrationInfo `bson:"integration" json:"integration"`
	OwnerID     string          `bson:"owner_id" json:"owner_id"`
	CreatedAt   time.Time       `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time       `bson:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time      `bson:"deleted_at,omitempty" json:"-"`
}
type TimePeriod struct {
	Started  time.Time `bson:"started" json:"started"`
	Ended    time.Time `bson:"ended" json:"ended"`
	Duration int       `bson:"duration" json:"duration"` // duration in seconds
}

type ReportStatus struct {
	Done        bool            `bson:"done" json:"done"`
	Integration IntegrationType `bson:"integration" json:"integration"` // e.g. "jira"
	ExternalID  string          `bson:"external_id" json:"external_id"` // e.g. "12345" from Jira
	ReportedAt  *time.Time      `bson:"reported_at,omitempty" json:"reported_at,omitempty"`
}

type TimeEntry struct {
	ID        string        `bson:"_id" json:"id"`
	ProjectID string        `bson:"project_id" json:"project_id"`
	OwnerID   string        `bson:"owner_id" json:"owner_id"`
	Period    TimePeriod    `bson:"period" json:"period"`
	Note      string        `bson:"note" json:"note"`
	Reported  *ReportStatus `bson:"reported,omitempty" json:"reported,omitempty"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time    `bson:"deleted_at,omitempty" json:"-"`
}
