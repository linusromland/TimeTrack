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

type Token struct {
	ID        string    `bson:"_id" json:"id"`
	Token     string    `bson:"token" json:"token"`
	UserID    string    `bson:"user_id" json:"user_id"`
	ExpiresAt time.Time `bson:"expires_at" json:"expires_at"`
	DeletedAt time.Time `bson:"deleted_at,omitempty" json:"-"`
}

type Integration struct {
	ID        string    `bson:"_id" json:"id"`
	Token     string    `bson:"token" json:"token"`
	UserID    string    `bson:"user_id" json:"user_id"`
	ExpiresAt time.Time `bson:"expires_at" json:"expires_at"`
}
