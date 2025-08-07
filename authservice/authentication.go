package main

import (
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

	logger.Info("Recieved Authenticate User Request", zap.String("user email", userEmail))

	// ? validation logic
	if userEmail == "" || userPassword == "" || !strings.Contains(userEmail, "@") || !strings.Contains(userEmail, ".") || len(userPassword) < 6 {
		logger.Warn("Invalid Request", zap.String("user email", userEmail))
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(iris.Map{"message": "The request contains missing or invalid fields"})
		return
	}

	// ? After fields are validated
	var existingUser *models.User

	userNotFoundError := dbConnector.Where("email = ?", user.Email).First(&existingUser).Error
	// ? If the user is already exist -> userNotFoundError = nil
	// ? If the user does not exist -> userNotFoundError = error

	if userNotFoundError != nil {
		logger.Warn("Authentication failed", zap.Error(userNotFoundError))
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(iris.Map{"message": "User not found"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(userPassword))

	if err != nil {
		logger.Warn("Authentication failed due to wrong password")
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(iris.Map{"message": "Authentication failed"})
		return
	}

	logger.Info("User Authenticated Successfully", zap.String("username", existingUser.Name),
		zap.String("useremail", userEmail))

	accessToken, err := jwtManager.GenerateAccessToken(existingUser)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(iris.Map{"message": "Couldn't generate access token"})
		return
	}

	refreshToken, err := jwtManager.GenerateRefreshToken(existingUser)
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
	// Use the new Verify method from your global jwtManager
	claims, err := jwtManager.Verify(requestPayload.AccessToken)
	if err != nil {
		logger.Warn("Token validation failed", zap.Error(err))
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(iris.Map{"message": "Invalid or expired token."})
		return
	}

	logger.Info("Token validated successfully", zap.String("user_email", claims.UserEmail))
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

	// Verify the refresh token
	claims, err := jwtManager.Verify(requestBody.RefreshToken)
	if err != nil {
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(iris.Map{"message": "Invalid or expired refresh token"})
		return
	}

	// Check if the token type is "refresh"
	if claims.Type != "refresh" {
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(iris.Map{"message": "Invalid token type provided"})
		return
	}

	// If valid, generate a new access token
	user := &models.User{Email: claims.UserEmail} // Create a user model from claims
	newAccessToken, err := jwtManager.GenerateAccessToken(user)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(iris.Map{"message": "Failed to generate new access token"})
		return
	}

	ctx.StatusCode(http.StatusOK)
	ctx.JSON(iris.Map{"access_token": newAccessToken})
}
