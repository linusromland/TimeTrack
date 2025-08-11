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

type TimeEntryStatPerDate struct {
	TimeFrame string `bson:"timeframe" json:"timeframe"`  // ISO 8601 format for the correct time format. (e.g. for day 2025-08-11, month 2025-08, week 2025-W32)
	TotalTime float64 `bson:"total_time" json:"total_time"` // total time in seconds
}

type TimeEntryPerProject struct {
	ProjectID string  `bson:"project_id" json:"project_id"`
	TotalTime float64 `bson:"total_time" json:"total_time"` // total time in seconds
}

type TimeEntryStatistics struct {
	TotalEntries      int64                  `json:"total_entries"`       // total number of time entries used for statistics
	TotalTime         int64                  `json:"total_time"`          // total time in seconds
	Format            string                 `json:"format"`              // e.g. "d" for days, "w" for weeks, "m" for months
	EntriesPerDate    []TimeEntryStatPerDate `json:"entries_per_date"`    // list of time entries per date in the specified format
	EntriesPerProject []TimeEntryPerProject  `json:"entries_per_project"` // list of time entries per project
}
