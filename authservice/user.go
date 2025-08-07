package main

import (
	"authservice/models"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"regexp"

	"github.com/kataras/iris/v12"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func validateUser(user models.User) error {
	if !isValidEmail(user.Email) {
		return fmt.Errorf("invalid email format")
	}

	if len(user.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	if !isValidPhone(user.Phone) {
		return fmt.Errorf("invalid phone number format")
	}
	return nil
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
	return re.MatchString(email)
}

func isValidPhone(phone string) bool {
	re := regexp.MustCompile(`^\d{10}$`)
	return re.MatchString(phone)
}

func AddUser(ctx iris.Context) {
	var user models.User
	ctx.ReadJSON(&user)

	logger.Info("Recieved User Request", zap.String("useremail", user.Email), zap.String("username", user.Name))

	if err := validateUser(user); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(iris.Map{"message": err.Error()})
		return
	}

	var existingUser models.User

	userNotFoundError := dbConnector.Where("email = ?", user.Email).First(&existingUser).Error

	if userNotFoundError == gorm.ErrRecordNotFound {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

		if err != nil {
			logger.Error("Fail to Hashing Password", zap.String("message", err.Error()))
		}
		userId := uuid.New()
		newUser := &models.User{
			ID:       userId,
			Name:     user.Name,
			Email:    user.Email,
			Password: string(hashedPassword),
			Phone:    user.Phone}

		primaryKey := dbConnector.Create(newUser)

		if primaryKey.Error != nil {
			logger.Error("Failed to Create user", zap.String("userPhone ", user.Phone), zap.Error(primaryKey.Error))
			ctx.StatusCode(http.StatusConflict)
			ctx.JSON(iris.Map{"message": "The Phone is already registered"})
			return
		}
		logger.Info(fmt.Sprintf("User %s created successfully", user.Name))
		ctx.JSON(iris.Map{"message": "User created successfully"})

	} else {
		logger.Warn("User Email Already Exist", zap.String("usermail", user.Email))
		ctx.StatusCode(http.StatusConflict)
		ctx.JSON(iris.Map{"message": "User Email Already Exist"})
	}

}

func GetAllUsers(ctx iris.Context) {
	var users []models.User

	err := dbConnector.Find(&users).Error

	if err != nil {
		logger.Error("Failed to fetch users from database", zap.Error(err))
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"message": "An error occurred while fetching users."})
		return
	}

	var userDetailsList []models.UserDetails

	for _, user := range users {
		userDetails := models.UserDetails{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Phone: user.Phone,
			Role:  user.Role,
		}
		userDetailsList = append(userDetailsList, userDetails)
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(userDetailsList)
}
