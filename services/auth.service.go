package services

import (
	"net/http"
	"strings"
	"time"
	"v2/initializers"
	"v2/models"
	"v2/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthService struct {
	DB *gorm.DB
}

func NewAuthService(DB *gorm.DB) AuthService {
	return AuthService{DB}
}

func (as *AuthService) SignUpUser(ctx *gin.Context) {
	var payload *models.SignUpInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		response := utils.BuildErrorResponse("Failed to process request", err.Error(), utils.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if payload.Password != payload.PasswordConfirm {
		response := utils.BuildErrorResponse("Passwords do not match", "Invalid Credential", utils.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		response := utils.BuildErrorResponse("Error", err.Error(), utils.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	now := time.Now()
	newUser := models.User{
		Name:      payload.Name,
		Email:     strings.ToLower(payload.Email),
		Password:  hashedPassword,
		Role:      "user",
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := as.DB.Create(&newUser)

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		response := utils.BuildErrorResponse("User with that email already exists", result.Error.Error(), utils.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusConflict, response)
		return
	} else if result.Error != nil {
		response := utils.BuildErrorResponse("Something bad happened", result.Error.Error(), utils.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadGateway, response)
		return
	}

	config := initializers.GetConfig()
	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, newUser.ID, config.AccessTokenKey)
	if err != nil {
		response := utils.BuildErrorResponse("AccessToken failed", err.Error(), utils.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	userResponse := &models.UserResponse{
		ID:          newUser.ID,
		Name:        newUser.Name,
		Email:       newUser.Email,
		Role:        newUser.Role,
		CreatedAt:   newUser.CreatedAt,
		UpdatedAt:   newUser.UpdatedAt,
		AccessToken: access_token,
	}
	response := utils.BuildResponse(true, "OK!", gin.H{"user": userResponse})
	ctx.JSON(http.StatusCreated, response)
}

func (as *AuthService) SignInUser(ctx *gin.Context) {
	var payload *models.SignInInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var user models.User
	result := as.DB.First(&user, "email = ?", strings.ToLower(payload.Email))
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	if err := utils.VerifyPassword(user.Password, payload.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	config := initializers.GetConfig()
	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenKey)
	if err != nil {
		response := utils.BuildErrorResponse("AccessToken failed", result.Error.Error(), utils.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	userResponse := &models.UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Role:        user.Role,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		AccessToken: access_token,
	}
	response := utils.BuildResponse(true, "OK!", gin.H{"user": userResponse})

	ctx.JSON(http.StatusOK, response)
}

func (as *AuthService) LogoutUser(ctx *gin.Context) {
	print(ctx.Get("currentUser"))
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
