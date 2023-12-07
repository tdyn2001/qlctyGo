package controllers

import (
	"v2/middleware"
	"v2/services"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) AuthController {
	return AuthController{authService}
}

func (ac *AuthController) AuthController(rg *gin.RouterGroup) {
	router := rg.Group("/user")
	router.POST("/register", ac.authService.SignUpUser)
	router.POST("/login", ac.authService.SignInUser)
	router.GET("/logout", middleware.AuthorizeJWT(), ac.authService.LogoutUser)

}
