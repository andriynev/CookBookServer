package user

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type ProfileRepository struct {
	db *gorm.DB
}

const secret  = "secret"

func GetProfileRepository(db *gorm.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

func (r *ProfileRepository) GetByIds(ids []uint) (users []Profile, err error) {
	if len(ids) == 0 {
		err = fmt.Errorf("user ids cannot be empty")
		return
	}
	err = r.db.Debug().Preload("Media").Preload("Activities").Where("id in (?)", ids).Find(&users).Error
	return
}

func (r *ProfileRepository) GetById(id uint) (user Profile, err error) {
	if id == 0 {
		err = fmt.Errorf("user id cannot be empty")
		return
	}
	err = r.db.Debug().Where(&Profile{Id: id}).First(&user).Error
	return
}

func (r *ProfileRepository) GetByUsername(username string) (user Profile, err error) {
	if len(username) == 0 {
		err = fmt.Errorf("user username cannot be empty")
		return
	}

	err = r.db.Where(&Profile{Username: username}).First(&user).Error
	return
}


func (r *ProfileRepository) Create(user Profile) (createdUser Profile, err error) {
	if user.Id != 0 {
		err = fmt.Errorf("user id should be empty")
		return
	}
	err = r.db.Exec("INSERT INTO users (username, password) "+
		"VALUES(?, AES_ENCRYPT(?, ?));",
		user.Username, user.Password, secret).Error
	if err != nil {
		return
	}

	createdUser, err = r.GetByUserParams(user)
	return
}

func (r *ProfileRepository) Update(id uint, profileDetails map[string]interface{}) (err error) {
	if id == 0 {
		err = fmt.Errorf("user id cannot be empty")
		return
	}
	fmt.Println(profileDetails)


	err = r.db.Debug().Model(Profile{}).Where(&Profile{Id: id}).Update(profileDetails).Error
	return
}

func (r *ProfileRepository) GetByUserParams(userParams Profile) (user Profile, err error) {
	if len(userParams.Username) == 0 {
		err = fmt.Errorf("user username cannot be empty")
		return
	}

	err = r.db.Where(&Profile{Username:userParams.Username}).
		Where("AES_DECRYPT(password, ?) = ?", secret, userParams.Password).First(&user).Error
	return
}
