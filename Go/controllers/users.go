package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"proyecto.com/m/conectors"
	"proyecto.com/m/models"
)

func GetUsersHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []models.User
		if err := db.Order("id asc").Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to load users", "error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

func GetUserInfoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")
		if username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "username is required"})
			return
		}

		user, err := getUserByUsername(db, username)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}

func LoginHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload conectors.LoginRequest
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid login payload"})
			return
		}

		user, err := getUserByUsername(db, payload.Username)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid username or password"})
			return
		}

		if user.Password != payload.Password {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid username or password"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "login successful", "user": user})
	}
}

func RegisterHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload conectors.RegisterRequest
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid register payload"})
			return
		}

		if payload.Username == "" || payload.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "username and password are required"})
			return
		}

		if _, err := getUserByUsername(db, payload.Username); err == nil {
			c.JSON(http.StatusConflict, gin.H{"message": "username already exists"})
			return
		}

		user := models.User{
			Name:     payload.Name,
			LastName: payload.LastName,
			Username: payload.Username,
			Email:    payload.Email,
			Password: payload.Password,
		}

		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not create user", "error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "user created successfully", "user": user})
	}
}

func UpdatePasswordHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload conectors.UpdatePasswordRequest
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid password payload"})
			return
		}

		if payload.Username == "" || payload.NewPassword == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "username and new password are required"})
			return
		}

		user, err := getUserByUsername(db, payload.Username)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
			return
		}

		if user.Password != payload.OldPassword {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "old password is incorrect"})
			return
		}

		user.Password = payload.NewPassword
		if err := db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not update password", "error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "password updated successfully"})
	}
}

func DeleteUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")
		if username == "" {
			username = c.Query("username")
		}
		if username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "username is required"})
			return
		}

		result := db.Where("username = ?", username).Delete(&models.User{})
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not delete user", "error": result.Error.Error()})
			return
		}
		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
	}
}

func getUserByUsername(db *gorm.DB, username string) (models.User, error) {
	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, err
		}
		return models.User{}, err
	}
	return user, nil
}
