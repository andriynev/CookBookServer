package media

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type MediaRepository struct {
	db *gorm.DB
}

func GetMediaRepository(db *gorm.DB) *MediaRepository {
	return &MediaRepository{db: db}
}

func (r *MediaRepository) GetById(id uint) (media Media, err error) {
	if id == 0 {
		err = fmt.Errorf("media id cannot be empty")
		return
	}
	err = r.db.Where(&Media{Id: id}).First(&media).Error
	return
}

func (r *MediaRepository) Create(media *Media) (err error) {
	if media == nil {
		err = fmt.Errorf("media cannot be empty")
		return
	}
	if media.Id != 0 {
		err = fmt.Errorf("media id should be empty")
		return
	}
	err = r.db.Create(media).Error
	return
}

func (r *MediaRepository) Update(media *Media) (err error) {
	if media == nil {
		err = fmt.Errorf("media cannot be empty")
		return
	}
	if media.Id == 0 {
		err = fmt.Errorf("media id cannot be empty")
		return
	}

	err = r.db.Model(Media{}).Update(media).Error
	return
}

func (r *MediaRepository) Delete(mediaId uint) (err error) {
	if mediaId == 0 {
		err = fmt.Errorf("media id cannot be empty")
		return
	}

	err = r.db.Model(Media{}).Where(&Media{Id: mediaId}).Delete(Media{}).Error
	return
}
