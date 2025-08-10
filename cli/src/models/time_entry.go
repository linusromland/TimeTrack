package models

import (
	"time"
)

type TimePeriod struct {
	Started  time.Time `bson:"started" json:"started"`
	Ended    time.Time `bson:"ended" json:"ended"`
	Duration int       `bson:"duration" json:"duration"` // duration in seconds
}

type ReportStatus struct {
	Done        bool       `bson:"done" json:"done"`
	Integration string     `bson:"integration" json:"integration"` // e.g. "jira"
	ExternalID  string     `bson:"external_id" json:"external_id"` // e.g. "12345" from Jira
	ReportedAt  *time.Time `bson:"reported_at,omitempty" json:"reported_at,omitempty"`
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
