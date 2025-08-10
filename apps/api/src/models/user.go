package models

import (
	"time"
)

type User struct {
	ID        string    `bson:"_id" json:"id"`
	Email     string    `bson:"email" json:"email"`
	Password  string    `bson:"password" json:"password,omitempty"`
	DeletedAt time.Time `bson:"deleted_at,omitempty" json:"-"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
	Integration UserIntegration `bson:"integration" json:"integration"`
}

type UserIntegration struct {
	Atlassian AtlassianIntegration `bson:"atlassian" json:"atlassian"`
}

type AtlassianIntegration struct {
	Enabled bool   `bson:"enabled" json:"enabled"`
	AccessToken string `bson:"access_token" json:"access_token"`
}