package handler

import (
	"food/src/api/models/media"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
)

// GetMedia godoc
// @Summary Get media from DB
// @Description get media
// @Tags media
// @Param   folder     path    string     true      "1ac4f9a135d204e71ed41fa8accfbe42"
// @Param   filename   path    string     true      "avatar.png"
// @Produce png
// @Produce jpeg
// @Produce gif
// @Produce json
// @Success 200 {array} integer
// @Failure 400 {object} handler.APIResponse
// @Failure 401 {object} handler.APIResponse
// @Failure 500 {object} handler.APIResponse
// @Security ApiKeyAuth
// @Router /v1/media/{folder}/{filename} [get]
func (*Controller) GetMedia(c *gin.Context) {
	folder := c.Param("folder")
	if len(folder) == 0 {
		c.JSON(http.StatusBadRequest, APIResponse{Message: "Request is invalid. Folder cannot be empty"})
		return
	}
	filename := c.Param("filename")
	if len(filename) == 0 {
		c.JSON(http.StatusBadRequest, APIResponse{Message: "Request is invalid. Filename cannot be empty"})
		return
	}

	fi, err := os.Stat(path.Join(media.MediaFolderRoot, folder, filename))
	if err != nil {
		log.Printf("Given file by path is not file. Stat err %s", err)
		c.JSON(http.StatusNotFound, APIResponse{Message: "Request is invalid. Filename or folder is invalid"})
		return
	}

	if fi.IsDir() {
		log.Printf("Given file path `%s` is dir", fi.Name())
		c.JSON(http.StatusNotFound,APIResponse{Message: "Request is invalid. Filename or folder is invalid"})
		return
	}
	c.File(path.Join(media.MediaFolderRoot, folder, filename))
}
