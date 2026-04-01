package models

import (
	"time"
)

// Node represents a registered HIL bench
type Node struct {
	ID              string    `gorm:"primaryKey" json:"id"` // Defaults to Hostname
	Hostname        string    `json:"hostname"`
	Status          string    `json:"status"` // "online", "offline", "in-use"
	AssignedSSHPort int       `gorm:"unique" json:"assigned_ssh_port"`
	LastSeenAt      time.Time `json:"last_seen_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// RegisterRequest is the payload sent by the Python agent
type RegisterRequest struct {
	Hostname string `json:"hostname"`
}
