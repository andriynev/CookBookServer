package services

import (
	"food/src/api/models/user"
	"github.com/jinzhu/gorm"
)

func GetProfileService(db *gorm.DB) *Profile {
	return &Profile{
		profileRepo: user.GetProfileRepository(db),
	}
}

type Profile struct {
	profileRepo *user.ProfileRepository
}
