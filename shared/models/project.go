package models

import (
	"time"
)

type IntegrationInfo struct {
	Type       string `bson:"type" json:"type"`               // e.g. "jira"
	Key        string `bson:"key" json:"key"`                 // e.g. "MNT-123"
	ExternalID string `bson:"external_id" json:"external_id"` // ID in the integration system
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
