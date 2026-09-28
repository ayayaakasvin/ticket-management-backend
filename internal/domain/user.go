package domain

import "time"

type User struct {
	ID           int64  `json:"user_id" example:"123"`
	Username     string `json:"username" example:"alice"`
	Email        string `json:"email" example:"example@example.com"`
	PasswordHash string `json:"hashed_password,omitempty" example:"$2a1ASDf"`
	Role         Role
	CreatedAt    time.Time `json:"created_at" example:"2025-11-14 10:23:45"`
}

type Role string

// Roles
const (
	Admin  Role = "admin"
	Client Role = "client"
)
