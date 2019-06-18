package services

import (
	"food/src/api/models/tools"
	"food/src/api/models/user"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func GetUserService(db *gorm.DB) *User {
	return &User{repo: user.GetProfileRepository(db)}
}

type User struct {
	repo *user.ProfileRepository
}

func (s *User) Create(request user.SignUpRequest) (profile user.Profile, err error) {
	err = tools.Validator.Struct(request)
	if err != nil {
		err = errors.Wrap(tools.NewValidationErr(err), "Sign up request is not valid")
		return
	}
	err = request.Validate()
	if err != nil {
		err = errors.Wrap(tools.NewValidationErr(err), "Sign up request is not valid")
		return
	}

	_, err = s.repo.GetByUsername(request.Username)
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return
	}

	if err == nil {
		err = errors.Wrapf(tools.NewValidationErr(err), "profile with username `%s` is already registered", request.Username)
		return
	}

	profileParams := user.Profile{
		Username: request.Username,
		Password: request.Password,
	}
	profile, err = s.repo.Create(profileParams)
	if err != nil {
		return
	}

	return
}

func (s *User) FindUser(request user.SignInRequest) (profile user.Profile, err error) {

	err = tools.Validator.Struct(request)
	if err != nil {
		err = tools.NewValidationErr(err)
		return
	}
	profileParams := user.Profile{Password: request.Password, Username: request.Username}
	profile, err = s.repo.GetByUserParams(profileParams)
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return
	}

	if err != nil {
		err = tools.NewValidationErr(err)
	}
	return
}
