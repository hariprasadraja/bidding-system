package model

import "time"

const (
	AdminUser = iota
	NormalUser
)

type User struct {
	ID           int64     `json:"id,omitempty"`
	Name         string    `json:"name,omitempty"`
	Role         int32     `json:"role,omitempty"`
	Email        string    `json:"email,omitempty"`
	Password     string    `json:"password,omitempty"`
	DateCreated  time.Time `json:"date_created,omitempty"`
	DateModified time.Time `json:"date_modified,omitempty"`
}
