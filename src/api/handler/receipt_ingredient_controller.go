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

type ListReceiptIngredientAPIResponse struct {
	APIResponse
	List []receipt.ReceiptIngredient `json:"list"`
}

type ReceiptIngredientAPIResponse struct {
	APIResponse
	Item receipt.ReceiptIngredient `json:"item"`
}

// GetReceiptIngredients godoc
// @Summary Get receipt ingredients
// @Description find receipt ingredients by params
// @Tags receipts
// @Produce  json
// @Param   id     path    int     true        "Receipt id"
// @Success 200 {object} handler.ListReceiptIngredientAPIResponse
// @Failure 401 {object} handler.APIResponse
// @Failure 400 {object} handler.APIResponse
// @Failure 500 {object} handler.APIResponse
// @Security ApiKeyAuth
// @Router /v1/receipts/{id}/ingredients [get]
func (*Controller) GetReceiptIngredients(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Message: "Given request to get receipt ingredient is invalid"})
		return
	}

	db, err := database.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{Message: "Error occurred when try to get receipts"})
		return
	}


	receiptService := services.GetReceiptService(db)
	receiptIngredients, err := receiptService.GetAllReceiptIngredientsById(uint(id))
	if err != nil {
		log.Printf("internal error: `%s`", err)
		c.JSON(http.StatusInternalServerError, APIResponse{Message: "Error occurred when get receipt ingredients"})
		return
	}


	c.JSON(http.StatusOK, ListReceiptIngredientAPIResponse{APIResponse: APIResponse{}, List: receiptIngredients})
}

// CreateReceiptIngredient godoc
// @Summary Create receipt ingredient
// @Tags receipts
// @Produce  json
// @Param   id     path    int     true        "Receipt id"
// @Param ingredient body services.CreateReceiptIngredientRequest true "params"
// @Success 200 {object} handler.ReceiptIngredientAPIResponse
// @Failure 401 {object} handler.APIResponse
// @Failure 400 {object} handler.APIResponse
// @Failure 500 {object} handler.APIResponse
// @Security ApiKeyAuth
// @Router /v1/receipts/{id}/ingredients [post]
func (*Controller) CreateReceiptIngredient(c *gin.Context) {
	var request services.CreateReceiptIngredientRequest
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
	i, err := svc.CreateReceiptIngredient(uint(id), userClaims.Id, request)
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


	c.JSON(http.StatusOK, ReceiptIngredientAPIResponse{APIResponse: APIResponse{}, Item: i})
}

// UpdateReceiptIngredient godoc
// @Summary Update receipt
// @Tags receipts
// @Produce  json
// @Param   id     path    int     true        "Receipt id"
// @Param   ingredient_id     path    int     true        "Receipt ingredient id"
// @Param ingredient body services.UpdateReceiptIngredientRequest true "params"
// @Success 200 {object} handler.ReceiptIngredientAPIResponse
// @Failure 401 {object} handler.APIResponse
// @Failure 400 {object} handler.APIResponse
// @Failure 403 {object} handler.APIResponse
// @Failure 500 {object} handler.APIResponse
// @Security ApiKeyAuth
// @Router /v1/receipts/{id}/ingredients/{ingredient_id} [put]
func (*Controller) UpdateReceiptIngredient(c *gin.Context) {
	var request services.UpdateReceiptIngredientRequest
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

	ingredientIdParam := c.Param("ingredient_id")
	ingredientId, err := strconv.Atoi(ingredientIdParam)
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
	i, err := svc.UpdateReceiptIngredient(uint(id), uint(ingredientId), userClaims.Id, request)
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


	c.JSON(http.StatusOK, ReceiptIngredientAPIResponse{APIResponse: APIResponse{}, Item: i})
}

// DeleteReceiptIngredient godoc
// @Summary Delete receipt ingredient
// @Tags receipts
// @Produce  json
// @Param   id     path    int     true        "Receipt id"
// @Param   ingredient_id     path    int     true        "Receipt ingredient id"
// @Success 204
// @Failure 401 {object} handler.APIResponse
// @Failure 400 {object} handler.APIResponse
// @Failure 403 {object} handler.APIResponse
// @Failure 500 {object} handler.APIResponse
// @Security ApiKeyAuth
// @Router /v1/receipts/{id}/ingredients/{ingredient_id} [delete]
func (*Controller) DeleteReceiptIngredient(c *gin.Context) {
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

	ingredientIdParam := c.Param("ingredient_id")
	ingredientId, err := strconv.Atoi(ingredientIdParam)
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
	err = svc.DeleteReceiptIngredient(uint(id), uint(ingredientId), userClaims.Id)
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
