package user

import (
	"time"
)

type Profile struct {
	// user identifier (required)
	Id                            uint         `json:"id" gorm:"primary_key"`
	Username                      string       `json:"username"`
	Password                      string       `json:"password"`
	CreatedAt                     time.Time    `json:"-"`
	UpdatedAt                     time.Time    `json:"-"`
	DeletedAt                     *time.Time   `json:"-"`
}

func (Profile) TableName() string {
	return "users"
}


