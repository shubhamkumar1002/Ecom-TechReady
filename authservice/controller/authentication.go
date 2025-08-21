package controller

import (
	"authservice/common"
	"authservice/models"
	"net/http"
	"strings"

	"github.com/kataras/iris/v12"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func AuthenticateUser(ctx iris.Context) {
	var user models.UserLogin
	ctx.ReadJSON(&user)

	userEmail := user.Email
	userPassword := user.Password

	common.App.Logger.Info("Recieved Authenticate User Request", zap.String("user email", userEmail))

	if userEmail == "" || userPassword == "" || !strings.Contains(userEmail, "@") || !strings.Contains(userEmail, ".") || len(userPassword) < 6 {
		common.App.Logger.Warn("Invalid Request", zap.String("user email", userEmail))
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(iris.Map{"message": "The request contains missing or invalid fields"})
		return
	}

	var existingUser *models.User

	userNotFoundError := common.App.DbConnector.Where("email = ?", user.Email).First(&existingUser).Error

	if userNotFoundError != nil {
		common.App.Logger.Warn("Authentication failed", zap.Error(userNotFoundError))
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(iris.Map{"message": "User not found"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(userPassword))

	if err != nil {
		common.App.Logger.Warn("Authentication failed due to wrong password")
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(iris.Map{"message": "Authentication failed"})
		return
	}

	common.App.Logger.Info("User Authenticated Successfully", zap.String("username", existingUser.Name),
		zap.String("useremail", userEmail))

	accessToken, err := common.App.JWTManager.GenerateAccessToken(existingUser)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(iris.Map{"message": "Couldn't generate access token"})
		return
	}

	refreshToken, err := common.App.JWTManager.GenerateRefreshToken(existingUser)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(iris.Map{"message": "Couldn't generate refresh token"})
		return
	}

	ctx.StatusCode(http.StatusOK)
	ctx.JSON(iris.Map{
		"message":       "User login successfully",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

type TokenValidationRequest struct {
	AccessToken string `json:"access_token"`
}

func ValidateToken(ctx iris.Context) {
	var requestPayload TokenValidationRequest

	if err := ctx.ReadJSON(&requestPayload); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid input"})
		return
	}

	claims, err := common.App.JWTManager.Verify(requestPayload.AccessToken)
	if err != nil {
		common.App.Logger.Warn("Token validation failed", zap.Error(err))
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(iris.Map{"message": "Invalid or expired token."})
		return
	}

	common.App.Logger.Info("Token validated successfully", zap.String("user_email", claims.UserEmail))
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(iris.Map{
		"message": "Token is valid",
		"email":   claims.UserEmail,
	})
}

func RefreshTokenHandler(ctx iris.Context) {
	var requestBody struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := ctx.ReadJSON(&requestBody); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(iris.Map{"message": "Invalid request body"})
		return
	}

	claims, err := common.App.JWTManager.Verify(requestBody.RefreshToken)
	if err != nil {
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(iris.Map{"message": "Invalid or expired refresh token"})
		return
	}

	if claims.Type != "refresh" {
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(iris.Map{"message": "Invalid token type provided"})
		return
	}

	user := &models.User{Email: claims.UserEmail}
	newAccessToken, err := common.App.JWTManager.GenerateAccessToken(user)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(iris.Map{"message": "Failed to generate new access token"})
		return
	}

	ctx.StatusCode(http.StatusOK)
	ctx.JSON(iris.Map{"access_token": newAccessToken})
}
