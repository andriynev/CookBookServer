package handler

import (
	"food/src/api/config"
	"food/src/api/jwt_auth"
	"food/src/api/models/user"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"              // gin-swagger middleware
	"github.com/swaggo/gin-swagger/swaggerFiles" // swagger embed files

	_ "food/src/api/docs"
)

type APIResponse struct {
	Message string `json:"message,omitempty"` // need fill only if error occurred
}

type AuthAPIResponse struct {
	APIResponse
	Id          uint   `json:"id,omitempty"`
	AccessToken string `json:"access_token"`
}

type ProfileAPIResponse struct {
	APIResponse
}

// @title Food API
// @version 1.0
// @description This is Swagger docs for Food API

// @host api.food.test
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func SetupHandler() http.Handler {
	r := gin.Default()

	c := NewController()

	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		config.GetConfig().DefaultUsername: config.GetConfig().DefaultPassword,
	}))
	authorized.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	authorized.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	userCtrl := r.Group("/user")
	{
		userCtrl.POST("/signUp", c.SignUp)
		userCtrl.POST("/signIn", c.SignIn)
	}

	// Secure API
	v1Api := r.Group("/v1")
	{
		ctrlSecure := v1Api.Use(auth())
		ctrlSecureRegular := ctrlSecure.Use(regularAccess())
		ctrlSecureRegular.GET("/media/:folder/:filename", c.GetMedia)

		ctrlSecureRegular.GET("/receipts", c.GetReceipts)
		ctrlSecureRegular.POST("/receipts", c.CreateReceipt)
		ctrlSecureRegular.PUT("/receipts/:id", c.UpdateReceipt)
		ctrlSecureRegular.DELETE("/receipts/:id", c.DeleteReceipt)
		ctrlSecureRegular.POST("/receipts/:id/media", c.UploadReceiptMedia)

		ctrlSecureRegular.GET("/receipts/:id/ingredients", c.GetReceiptIngredients)
		ctrlSecureRegular.POST("/receipts/:id/ingredients/", c.CreateReceiptIngredient)
		ctrlSecureRegular.PUT("/receipts/:id/ingredients/:ingredient_id", c.UpdateReceiptIngredient)
		ctrlSecureRegular.DELETE("/receipts/:id/ingredients/:ingredient_id", c.DeleteReceiptIngredient)

		ctrlSecureRegular.GET("/receipts/:id/directions", c.GetReceiptDirections)
		ctrlSecureRegular.POST("/receipts/:id/directions/", c.CreateReceiptDirection)
		ctrlSecureRegular.PUT("/receipts/:id/directions/:direction_id", c.UpdateReceiptDirection)
		ctrlSecureRegular.DELETE("/receipts/:id/directions/:direction_id", c.DeleteReceiptDirection)

		ctrlSecureRegular.GET("/ingredients", c.GetIngredients)
		ctrlSecureRegular.POST("/ingredients", c.CreateIngredient)
		ctrlSecureRegular.PUT("/ingredients/:id", c.UpdateIngredient)
	}

	return r
}

func auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 {
			log.Printf("auth header has wrong format. auth header: `%s`", authHeader)
			c.AbortWithStatus(401)
			return
		}

		if authHeaderParts[0] != "Bearer" {
			log.Printf("auth header has wrong format. auth marker: `%s`", authHeaderParts[0])
			c.AbortWithStatus(401)
			return
		}

		accessToken := authHeaderParts[1]
		userClaims, err := jwt_auth.ParseToken(accessToken)
		if err != nil {
			log.Printf("token is not valid: `%s`", err.Error())
			c.AbortWithStatus(401)
			return
		}
		c.Set("claims", userClaims)
		c.Next()
	}
}

func regularAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := c.Get("claims")
		if !ok {
			c.AbortWithStatus(401)
			return
		}

		userClaims, ok := claims.(*jwt_auth.UserClaims)
		if !ok {
			c.AbortWithStatus(401)
			return
		}

		if !userClaims.Privileges.Has(user.RegularUserPrivilege) {
			c.AbortWithStatus(403)
			return
		}
		c.Next()
	}
}
