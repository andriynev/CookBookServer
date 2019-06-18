package services

import (
	"crypto/md5"
	"fmt"
	"food/src/api/models/media"
	"food/src/api/models/tools"
	"github.com/pkg/errors"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/jinzhu/gorm"
)

type Media struct {
	mediaRepo *media.MediaRepository
}

func GetMediaService(db *gorm.DB) *Media {
	return &Media{mediaRepo: media.GetMediaRepository(db)}
}

type Options struct {
	Filename string
}

func (s *Media) ProcessMedia(formFile *multipart.FileHeader, opts Options) (newMedia media.Media, err error) {
	dateBytes, _ := time.Now().MarshalBinary()
	hashBytes := md5.Sum(dateBytes)

	currentUploadFolder := fmt.Sprintf("%x", hashBytes)
	folderPath := path.Join(media.MediaFolderRoot, currentUploadFolder)
	err = os.MkdirAll(folderPath, os.ModeDir)
	if err != nil {
		err =errors.Wrap(err, "Error occurred when create new folder.")
		return
	}
	filename := opts.Filename + filepath.Ext(formFile.Filename)
	mediaPath := path.Join(currentUploadFolder, filename)
	filePath := path.Join(folderPath, filename)
	newFile, err := os.Create(filePath)
	if err != nil {
		err = errors.Wrap(err, "Error occurred when create new file.")
		return
	}
	fd, err := formFile.Open()
	if err != nil {
		err = errors.Wrap(err, "Error occurred when open form file")
		return
	}
	_, err = io.Copy(newFile, fd)
	if err != nil {
		err = errors.Wrap(err, "Error occurred when fill file")
		return
	}
	fd.Close()

	newMedia = media.Media{Link: mediaPath, Format: formFile.Header.Get("Content-Type")}
	err = s.Save(&newMedia)

	if err != nil {
		err = errors.Wrap(err, "Error occurred when save new media")
		return
	}
	return

}

func (s *Media) Save(media *media.Media) (err error) {

	if media.Id != 0 {

		err = tools.NewValidationErr(fmt.Errorf("new media id is should be empty"))
		return
	}

	err = s.mediaRepo.Create(media)
	if err != nil {
		return
	}

	return nil
}

func (s *Media) Delete(mediaId uint) (err error) {
	if mediaId == 0 {
		err = tools.NewValidationErr(fmt.Errorf("media id cannot be empty"))
		return
	}
	existingMedia, err := s.mediaRepo.GetById(mediaId)
	if gorm.IsRecordNotFoundError(err) {
		err = tools.NewValidationErr(fmt.Errorf("media with id `%d` not found in DB", mediaId))
		return
	}
	if err != nil {
		return
	}

	if len(existingMedia.Link) == 0 {
		err = fmt.Errorf("media link is not defined")
		return
	}

	err = os.Remove(path.Join(media.MediaFolderRoot, existingMedia.Link))
	if err != nil {
		err = errors.Wrap(err, "Error occurred when remove file")
		return
	}

	err = s.mediaRepo.Delete(existingMedia.Id)
	if gorm.IsRecordNotFoundError(err) {
		err = nil
		return
	}

	if err != nil {
		err = errors.Wrap(err, "Error occurred when delete media from DB")
		return
	}

	return nil
}
