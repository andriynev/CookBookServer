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

type ListReceiptDirectionAPIResponse struct {
	APIResponse
	List []receipt.ReceiptDirection `json:"list"`
}

type ReceiptDirectionAPIResponse struct {
	APIResponse
	Item receipt.ReceiptDirection `json:"item"`
}

// GetReceiptDirections godoc
// @Summary Get receipt directions
// @Description find receipt directions by params
// @Tags receipts
// @Produce  json
// @Param   id     path    int     true        "Receipt id"
// @Success 200 {object} handler.ListReceiptDirectionAPIResponse
// @Failure 401 {object} handler.APIResponse
// @Failure 400 {object} handler.APIResponse
// @Failure 500 {object} handler.APIResponse
// @Security ApiKeyAuth
// @Router /v1/receipts/{id}/directions [get]
func (*Controller) GetReceiptDirections(c *gin.Context) {

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Message: "Given request to get receipt directions is invalid"})
		return
	}

	db, err := database.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{Message: "Error occurred when try to get receipts"})
		return
	}


	receiptService := services.GetReceiptService(db)
	directions, err := receiptService.GetAllReceiptDirectionsById(uint(id))
	if err != nil {
		log.Printf("internal error: `%s`", err)
		c.JSON(http.StatusInternalServerError, APIResponse{Message: "Error occurred when get receipt directions"})
		return
	}


	c.JSON(http.StatusOK, ListReceiptDirectionAPIResponse{APIResponse: APIResponse{}, List: directions})
}

// CreateReceiptDirection godoc
// @Summary Create receipt direction
// @Tags receipts
// @Produce  json
// @Param   id     path    int     true        "Receipt id"
// @Param ingredient body services.CreateReceiptDirectionRequest true "params"
// @Success 200 {object} handler.ReceiptDirectionAPIResponse
// @Failure 401 {object} handler.APIResponse
// @Failure 400 {object} handler.APIResponse
// @Failure 500 {object} handler.APIResponse
// @Security ApiKeyAuth
// @Router /v1/receipts/{id}/directions [post]
func (*Controller) CreateReceiptDirection(c *gin.Context) {
	var request services.CreateReceiptDirectionRequest
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

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Message: "Given request is invalid"})
		return
	}

	db, err := database.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{Message: "Error occurred when try to process request"})
		return
	}


	svc := services.GetReceiptService(db)
	i, err := svc.CreateReceiptDirection(uint(id), userClaims.Id, request)
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


	c.JSON(http.StatusOK, ReceiptDirectionAPIResponse{APIResponse: APIResponse{}, Item: i})
}

// UpdateReceiptDirection godoc
// @Summary Update receipt direction
// @Tags receipts
// @Produce  json
// @Param   id     path    int     true        "Receipt id"
// @Param   direction_id     path    int     true        "Receipt direction id"
// @Param direction body services.UpdateReceiptDirectionRequest true "params"
// @Success 200 {object} handler.ReceiptDirectionAPIResponse
// @Failure 401 {object} handler.APIResponse
// @Failure 400 {object} handler.APIResponse
// @Failure 403 {object} handler.APIResponse
// @Failure 500 {object} handler.APIResponse
// @Security ApiKeyAuth
// @Router /v1/receipts/{id}/directions/{direction_id} [put]
func (*Controller) UpdateReceiptDirection(c *gin.Context) {
	var request services.UpdateReceiptDirectionRequest
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

	directionIdParam := c.Param("direction_id")
	directionId, err := strconv.Atoi(directionIdParam)
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
	i, err := svc.UpdateReceiptDirection(uint(id), uint(directionId), userClaims.Id, request)
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


	c.JSON(http.StatusOK, ReceiptDirectionAPIResponse{APIResponse: APIResponse{}, Item: i})
}

// DeleteReceiptDirection godoc
// @Summary Delete receipt direction
// @Tags receipts
// @Produce  json
// @Param   id     path    int     true        "Receipt id"
// @Param   direction_id     path    int     true        "Receipt direction id"
// @Success 204
// @Failure 401 {object} handler.APIResponse
// @Failure 400 {object} handler.APIResponse
// @Failure 403 {object} handler.APIResponse
// @Failure 500 {object} handler.APIResponse
// @Security ApiKeyAuth
// @Router /v1/receipts/{id}/directions/{direction_id} [delete]
func (*Controller) DeleteReceiptDirection(c *gin.Context) {
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

	directionIdParam := c.Param("direction_id")
	directionId, err := strconv.Atoi(directionIdParam)
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
	err = svc.DeleteReceiptDirection(uint(id), uint(directionId), userClaims.Id)
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


	c.Status(http.StatusNoContent)
}
