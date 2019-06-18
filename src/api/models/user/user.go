package user

import (
	"fmt"
	"strings"
)

const (
	MinPasswordLen       = 6
	MaxPasswordLen       = 255
	MinLoginLen = 3
	MaxLoginLen = 50
	MinEmailLen = 3
	MaxEmailLen = 255
	RegularUserPrivilege = "regular"
)

var DefaultUserPrivileges = map[string]bool{
	RegularUserPrivilege: true,
}

type SignUpRequest struct {
	// (required)
	Username    string     `json:"username" minLength:"3" maxLength:"50" binding:"required" validate:"max=50,min=3"`
	// (required)
	Password    string     `json:"password" minLength:"6" maxLength:"255" binding:"required" validate:"max=255,min=6"`
}



func (u *SignUpRequest) TrimSpaces() {
	u.Username = strings.TrimSpace(u.Username)
}

func (u *SignUpRequest) Validate() (err error) {
	if len(u.Username) == 0 {
		err = fmt.Errorf("field username cannot be empty")
		return
	}

	if len(u.Username) < MinLoginLen {
		err = fmt.Errorf("field username is too short. "+
			"given length: `%d`, min length: `%d`", len(u.Username), MinLoginLen)
		return
	}

	if len(u.Username) > MaxLoginLen {
		err = fmt.Errorf("field username is too long. "+
			"given length: `%d`, max length: `%d`", len(u.Username), MaxLoginLen)
		return
	}

	if len(u.Password) == 0 {
		err = fmt.Errorf("field password cannot be empty")
		return
	}

	if len(u.Password) < MinPasswordLen {
		err = fmt.Errorf("field password is too short. "+
			"given length: `%d`, min length: `%d`", len(u.Password), MinPasswordLen)
		return
	}

	if len(u.Password) > MaxPasswordLen {
		err = fmt.Errorf("field password is too long. "+
			"given length: `%d`, max length: `%d`", len(u.Password), MaxPasswordLen)
		return
	}

	return
}

type SignInRequest struct {
	// (required)
	Username    string     `json:"username" minLength:"3" maxLength:"50" binding:"required" validate:"max=50,min=3"`
	// (required)
	Password    string     `json:"password" minLength:"6" maxLength:"255" binding:"required" validate:"max=255,min=6"`
}

func (u *SignInRequest) TrimSpaces() {
	u.Username = strings.TrimSpace(u.Username)
}
