package handler

import (
	"fmt"
	"food/src/api/database"
	"food/src/api/models/ingredient"
	"food/src/api/models/tools"
	"food/src/api/services"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strconv"
)

type ListIngredientsAPIResponse struct {
	APIResponse
	List []ingredient.Ingredient `json:"list"`
}

type IngredientAPIResponse struct {
	APIResponse
	Item ingredient.Ingredient `json:"item"`
}

// GetIngredients godoc
// @Summary Get ingredients
// @Description find ingredients by params
// @Tags receipts
// @Produce  json
// @Success 200 {object} handler.ListIngredientsAPIResponse
// @Failure 401 {object} handler.APIResponse
// @Failure 400 {object} handler.APIResponse
// @Failure 500 {object} handler.APIResponse
// @Security ApiKeyAuth
// @Router /v1/ingredients/ [get]
func (*Controller) GetIngredients(c *gin.Context) {
	db, err := database.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{Message: "Error occurred when try to get ingredients"})
		return
	}


	svc := services.GetReceiptService(db)
	items, err := svc.GetAllIngredients()
	if err != nil {
		log.Printf("get ingredients error: `%s`", err)
		c.JSON(http.StatusInternalServerError, APIResponse{Message: "Error occurred when get ingredients"})
		return
	}


	c.JSON(http.StatusOK, ListIngredientsAPIResponse{APIResponse: APIResponse{}, List: items})
}

// CreateIngredient godoc
// @Summary Create ingredient
// @Tags receipts
// @Produce  json
// @Param ingredient body services.CreateIngredientRequest true "params"
// @Success 200 {object} handler.IngredientAPIResponse
// @Failure 401 {object} handler.APIResponse
// @Failure 400 {object} handler.APIResponse
// @Failure 500 {object} handler.APIResponse
// @Security ApiKeyAuth
// @Router /v1/ingredients/ [post]
func (*Controller) CreateIngredient(c *gin.Context) {
	var request services.CreateIngredientRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Message: "Given request to create ingredient is invalid"})
		return
	}

	db, err := database.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{Message: "Error occurred when try to create ingredient"})
		return
	}


	svc := services.GetReceiptService(db)
	i, err := svc.CreateIngredient(request)
	if err != nil {
		switch errors.Cause(err).(type) {
		case *tools.ValidationErr:
			log.Printf("validate error %s", err)
			c.JSON(http.StatusBadRequest, APIResponse{Message: fmt.Sprintf("Given request is invalid.")})
			return
		}
		log.Printf("get ingredients error: `%s`", err)
		c.JSON(http.StatusInternalServerError, APIResponse{Message: "Error occurred when create ingredient"})
		return
	}


	c.JSON(http.StatusOK, IngredientAPIResponse{APIResponse: APIResponse{}, Item: i})
}

// UpdateIngredient godoc
// @Summary Update ingredient
// @Tags receipts
// @Produce  json
// @Param   id     path    int     true        "Ingredient id"
// @Param ingredient body services.UpdateIngredientRequest true "params"
// @Success 200 {object} handler.IngredientAPIResponse
// @Failure 401 {object} handler.APIResponse
// @Failure 400 {object} handler.APIResponse
// @Failure 500 {object} handler.APIResponse
// @Security ApiKeyAuth
// @Router /v1/ingredients/{id} [put]
func (*Controller) UpdateIngredient(c *gin.Context) {
	var request services.UpdateIngredientRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Message: "Given request to update ingredient is invalid"})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Message: "Given request to update ingredient is invalid"})
		return
	}

	db, err := database.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{Message: "Error occurred when try to update ingredient"})
		return
	}

	svc := services.GetReceiptService(db)
	i, err := svc.UpdateIngredient(uint(id), request)
	if err != nil {
		switch errors.Cause(err).(type) {
		case *tools.ValidationErr:
			log.Printf("validate error %s", err)
			c.JSON(http.StatusBadRequest, APIResponse{Message: fmt.Sprintf("Given request is invalid.")})
			return
		}
		log.Printf("get ingredients error: `%s`", err)
		c.JSON(http.StatusInternalServerError, APIResponse{Message: "Error occurred when update ingredient"})
		return
	}


	c.JSON(http.StatusOK, IngredientAPIResponse{APIResponse: APIResponse{}, Item: i})
}
