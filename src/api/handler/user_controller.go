package handler

import (
	"fmt"
	"food/src/api/database"
	"food/src/api/jwt_auth"
	"food/src/api/models/tools"
	"food/src/api/models/user"
	"food/src/api/services"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"log"
	"net/http"
)

// SignUp godoc
// @Summary Add a new user to the DB
// @Description create new user by params
// @Tags auth
// @Accept  json
// @Produce  json
// @Param account body user.SignUpRequest true "User params"
// @Success 200 {object} handler.AuthAPIResponse
// @Failure 400 {object} handler.APIResponse
// @Failure 500 {object} handler.APIResponse
// @Router /user/signUp [post]
func (*Controller) SignUp(c *gin.Context) {

	var params user.SignUpRequest
	err := c.ShouldBindJSON(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Message: fmt.Sprintf("Given request to create user is invalid. Orig err: `%s`", err)})
		return
	}

	db, err := database.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{Message: "Error occurred when try to create new user"})
		return
	}
	params.TrimSpaces()

	userService := services.GetUserService(db)
	newUser, err := userService.Create(params)
	if err != nil {
		switch errors.Cause(err).(type) {
		case *tools.ValidationErr:
			log.Printf("create new user validation error: `%s`", err)
			c.JSON(http.StatusBadRequest, APIResponse{Message: fmt.Sprintf("Given request is invalid. Orig err: `%s`", err)})
			return
		}

		log.Printf("create new user error: `%s`", err)
		c.JSON(http.StatusInternalServerError, APIResponse{Message: "Error occurred when create new user"})
		return
	}

	token, err := jwt_auth.GetToken(jwt_auth.GetUserClaims(newUser.Id, user.DefaultUserPrivileges))
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{Message: "Error occurred when generate access token"})
		return
	}

	c.JSON(http.StatusOK, AuthAPIResponse{APIResponse: APIResponse{}, Id: newUser.Id, AccessToken: token})
}

// SignIn godoc
// @Summary Get a user from the DB
// @Description get a user by params
// @Tags auth
// @Accept  json
// @Produce  json
// @Param account body user.SignInRequest true "User params"
// @Success 200 {object} handler.AuthAPIResponse
// @Failure 400 {object} handler.APIResponse
// @Failure 401 {object} handler.APIResponse
// @Failure 500 {object} handler.APIResponse
// @Router /user/signIn [post]
func (*Controller) SignIn(c *gin.Context) {
	var request user.SignInRequest
	err := c.Bind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Message: "Given request to get user is invalid"})
		return
	}
	log.Println(request)

	db, err := database.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{Message: "Error occurred when try to get user"})
		return
	}

	request.TrimSpaces()


	userService := services.GetUserService(db)
	newUser, err := userService.FindUser(request)
	if err != nil {
		switch errors.Cause(err).(type) {
		case *tools.ValidationErr:
			log.Printf("validate error %s", err)
			c.JSON(http.StatusUnauthorized, APIResponse{Message: fmt.Sprintf("Given credentials is invalid.")})
			return
		}

		log.Printf("internal error %s", err)
		c.JSON(http.StatusInternalServerError, APIResponse{Message: "Error occurred when get user"})
		return
	}

	token, err := jwt_auth.GetToken(jwt_auth.GetUserClaims(newUser.Id, user.DefaultUserPrivileges))
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{Message: "Error occurred when generate access token"})
		return
	}

	c.JSON(http.StatusOK, AuthAPIResponse{APIResponse: APIResponse{}, Id: newUser.Id, AccessToken: token})
}
