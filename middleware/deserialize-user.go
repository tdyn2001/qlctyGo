package middleware

import (
	"fmt"
	"net/http"
	"v2/initializers"
	"v2/models"
	"v2/utils"

	"github.com/gin-gonic/gin"
)

func AuthorizeJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		access_token := ctx.GetHeader("x_access_token")
		if access_token == "" {
			response := utils.BuildErrorResponse("Failed to process request", "No token found", nil)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
		sub, err := utils.ValidateToken(access_token, initializers.GetConfig().AccessTokenKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
			return
		}

		var user models.User
		result := initializers.DB.First(&user, "id = ?", fmt.Sprint(sub))
		if result.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
			return
		}

		ctx.Set("currentUser", user)
		ctx.Next()
	}
}
