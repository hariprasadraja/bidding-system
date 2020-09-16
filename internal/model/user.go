package model

import (
	"encoding/json"
	"time"

	log "github.com/micro/go-micro/v2/logger"
)

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

type ValidationError map[string]interface{}

func (ve ValidationError) Error() string {
	data, err := json.Marshal(ve)
	if err != nil {
		log.Error("validation.error.marshall", err.Error())
	}

	return string(data)
}

func (u User) CreateValidate() ValidationError {
	errors := make(ValidationError)

	if u.Name == "" {
		errors["name"] = "user name should not be empty."
	}

	if u.Email == "" {
		errors["email"] = "user email should not be empty."
	}

	if u.Password == "" {
		errors["password"] = "user password should not be empty."
	}

	return errors
}

func (u User) UpdateValidate() ValidationError {
	errors := make(ValidationError)

	if u.ID == 0 {
		errors["id"] = "user id is required."
	}

	if u.Name == "" {
		errors["name"] = "user name should not be empty."
	}

	if u.Email == "" {
		errors["email"] = "user email should not be empty."
	}

	return errors
}
