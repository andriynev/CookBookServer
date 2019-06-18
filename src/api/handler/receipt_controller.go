package handler

import (
	"fmt"
	"food/src/api/database"
	"food/src/api/jwt_auth"
	"food/src/api/models/receipt"
	"food/src/api/models/tools"
	"food/src/api/services"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strconv"
)

type ListAPIResponse struct {
	APIResponse
	List []receipt.Receipt `json:"list"`
}

type ReceiptAPIResponse struct {
	APIResponse
	Item receipt.Receipt `json:"item"`
}

// GetReceipts godoc
// @Summary Get receipts
// @Description find receipts by params
// @Tags receipts
// @Produce  json
// @Param category query string false "category"
// @Success 200 {object} handler.ListAPIResponse
// @Failure 401 {object} handler.APIResponse
// @Failure 400 {object} handler.APIResponse
// @Failure 500 {object} handler.APIResponse
// @Security ApiKeyAuth
// @Router /v1/receipts/ [get]
func (*Controller) GetReceipts(c *gin.Context) {
	category := c.Request.URL.Query().Get("category")
	db, err := database.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{Message: "Error occurred when try to get receipts"})
		return
	}


	receiptService := services.GetReceiptService(db)
	categories, err := receiptService.GetAllReceiptsByCategory(category)
	if err != nil {
		log.Printf("get receipts error: `%s`", err)
		c.JSON(http.StatusInternalServerError, APIResponse{Message: "Error occurred when get receipts"})
		return
	}


	c.JSON(http.StatusOK, ListAPIResponse{APIResponse: APIResponse{}, List: categories})
}

// CreateReceipt godoc
// @Summary Create receipt
// @Tags receipts
// @Produce  json
// @Param ingredient body services.CreateReceiptRequest true "params"
// @Success 200 {object} handler.ReceiptAPIResponse
// @Failure 401 {object} handler.APIResponse
// @Failure 400 {object} handler.APIResponse
// @Failure 500 {object} handler.APIResponse
// @Security ApiKeyAuth
// @Router /v1/receipts/ [post]
func (*Controller) CreateReceipt(c *gin.Context) {
	var request services.CreateReceiptRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Message: "Given request to create receipt is invalid"})
		return
	}

	claims, _ := c.Get("claims")
	userClaims, ok := claims.(*jwt_auth.UserClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, APIResponse{Message: "Unauthorized access"})
		return
	}

	db, err := database.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{Message: "Error occurred when try to create receipt"})
		return
	}


	svc := services.GetReceiptService(db)
	i, err := svc.CreateReceipt(userClaims.Id, request)
	if err != nil {
		switch errors.Cause(err).(type) {
		case *tools.NotPermittedErr:
			c.JSON(http.StatusForbidden, APIResponse{Message: "Not permitted"})
			return
		case *tools.ValidationErr:
			log.Printf("validate error %s", err)
			c.JSON(http.StatusBadRequest, APIResponse{Message: fmt.Sprintf("Given request is invalid.")})
			return
		}
		log.Printf("internal error: `%s`", err)
		c.JSON(http.StatusInternalServerError, APIResponse{Message: "Error occurred when create receipt"})
		return
	}


	c.JSON(http.StatusOK, ReceiptAPIResponse{APIResponse: APIResponse{}, Item: i})
}

// UpdateReceipt godoc
// @Summary Update receipt
// @Tags receipts
// @Produce  json
// @Param   id     path    int     true        "Receipt id"
// @Param ingredient body services.UpdateReceiptRequest true "params"
// @Success 200 {object} handler.ReceiptAPIResponse
// @Failure 401 {object} handler.APIResponse
// @Failure 400 {object} handler.APIResponse
// @Failure 403 {object} handler.APIResponse
// @Failure 500 {object} handler.APIResponse
// @Security ApiKeyAuth
// @Router /v1/receipts/{id} [put]
func (*Controller) UpdateReceipt(c *gin.Context) {
	var request services.UpdateReceiptRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Message: "Given request to update receipt is invalid"})
		return
	}

	claims, _ := c.Get("claims")
	userClaims, ok := claims.(*jwt_auth.UserClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, APIResponse{Message: "Unauthorized access"})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Message: "Given request to update receipt is invalid"})
		return
	}

	db, err := database.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{Message: "Error occurred when try to update receipt"})
		return
	}

	svc := services.GetReceiptService(db)
	i, err := svc.UpdateReceipt(uint(id), userClaims.Id, request)
	if err != nil {
		switch errors.Cause(err).(type) {
		case *tools.NotPermittedErr:
			log.Printf("validate error %s", err)
			c.JSON(http.StatusForbidden, APIResponse{Message: fmt.Sprintf("Not permitted")})
			return
		case *tools.ValidationErr:
			log.Printf("validate error %s", err)
			c.JSON(http.StatusBadRequest, APIResponse{Message: fmt.Sprintf("Given request is invalid.")})
			return
		}
		log.Printf("internal error: `%s`", err)
		c.JSON(http.StatusInternalServerError, APIResponse{Message: "Error occurred when update receipt"})
		return
	}


	c.JSON(http.StatusOK, ReceiptAPIResponse{APIResponse: APIResponse{}, Item: i})
}

// DeleteReceipt godoc
// @Summary Delete receipt
// @Tags receipts
// @Produce  json
// @Param   id     path    int     true        "Receipt id"
// @Success 204
// @Failure 401 {object} handler.APIResponse
// @Failure 400 {object} handler.APIResponse
// @Failure 500 {object} handler.APIResponse
// @Security ApiKeyAuth
// @Router /v1/receipts/{id} [delete]
func (*Controller) DeleteReceipt(c *gin.Context) {
	claims, _ := c.Get("claims")
	userClaims, ok := claims.(*jwt_auth.UserClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, APIResponse{Message: "Unauthorized access"})
		return
	}
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Message: "Given request to update receipt is invalid"})
		return
	}
	if id == 0 {
		c.JSON(http.StatusBadRequest, APIResponse{Message: "Given request to update receipt is invalid"})
		return
	}

	db, err := database.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{Message: "Error occurred when try to update receipt"})
		return
	}

	svc := services.GetReceiptService(db)
	err = svc.DeleteReceipt(uint(id), userClaims.Id)
	if err != nil {
		switch errors.Cause(err).(type) {
		case *tools.ValidationErr:
			log.Printf("validate error %s", err)
			c.JSON(http.StatusUnauthorized, APIResponse{Message: fmt.Sprintf("Given request is invalid.")})
			return
		}
		log.Printf("get ingredients error: `%s`", err)
		c.JSON(http.StatusInternalServerError, APIResponse{Message: "Error occurred when update receipt"})
		return
	}


	c.Status(http.StatusNoContent)
}

// UploadReceiptMedia godoc
// @Summary Update a receipt media
// @Description upload a receipt media
// @Tags receipts
// @Accept  multipart/form-data
// @Produce json
// @Param   id     path    int     true        "Receipt id"
// @Param media_file formData file true "Media file"
// @Success 200 {object} handler.APIResponse
// @Failure 400 {object} handler.APIResponse
// @Failure 401 {object} handler.APIResponse
// @Failure 403 {object} handler.APIResponse
// @Failure 500 {object} handler.APIResponse
// @Security ApiKeyAuth
// @Router /v1/receipts/{id}/media [post]
func (*Controller) UploadReceiptMedia(c *gin.Context) {
	claims, _ := c.Get("claims")
	userClaims, ok := claims.(*jwt_auth.UserClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, APIResponse{Message: "Unauthorized access"})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Message: "Given request to update receipt is invalid"})
		return
	}


	if id == 0 {
		c.JSON(http.StatusBadRequest, APIResponse{Message: "Given request to update receipt is invalid"})
		return
	}

	if c.ContentType() == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Message: http.ErrNotMultipart.Error(),
		})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Message: fmt.Sprintf("Given request isn't valid multipart/form-data. Orig err: `%s`", err.Error()),
		})
		return
	}

	if len(form.File) == 0 {
		c.JSON(http.StatusBadRequest, APIResponse{
			Message: "Given multipart form is empty.",
		})
		return
	}

	if len(form.File) > 1 {
		c.JSON(http.StatusBadRequest, APIResponse{
			Message: "Given multipart form contains more than 1 file.",
		})
		return
	}

	if len(form.File["media_file"]) != 1 {
		c.JSON(http.StatusBadRequest, APIResponse{
			Message: "Given multipart isn't valid.",
		})
		return
	}

	db, err := database.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{Message: "Error occurred when try to update post media"})
		return
	}

	fileHeader := form.File["media_file"][0]

	svc := services.GetReceiptService(db)
	err = svc.UpdateReceiptMedia(uint(id), userClaims.Id, fileHeader, services.Options{Filename:"dish"})
	if err != nil {
		switch errors.Cause(err).(type) {
		case *tools.NotPermittedErr:
			c.JSON(http.StatusForbidden, APIResponse{
				Message: "Not permitted",
			})
			return
		case *tools.ValidationErr:
			c.JSON(http.StatusBadRequest, APIResponse{
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, APIResponse{Message: "Error occurred when try to update receipt media"})
		return
	}


	c.JSON(http.StatusOK, APIResponse{})
}
