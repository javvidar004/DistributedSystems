package controllers

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"proyecto.com/m/conectors"
	"proyecto.com/m/models"
)

func LoginHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload conectors.LoginRequest
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid login payload"})
			return
		}

		user, err := getUserByUsername(db, payload.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid username or password"})
			return
		}

		//if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(payload.Password)); err != nil {
		//	c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid username or password"})
		//	return
		//}

		if user.PasswordHash != payload.Password {
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

		if payload.Email == "" || payload.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "email and password are required"})
			return
		}

		if _, err := getUserByUsername(db, payload.Email); err == nil {
			c.JSON(http.StatusConflict, gin.H{"message": "username already exists"})
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not hash password", "error": err.Error()})
			return
		}

		uuid, err := uuid.NewRandom()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not generate UUID", "error": err.Error()})
			return
		}

		user := models.User{
			ID:           uuid.String(),
			Email:        payload.Email,
			PasswordHash: strings.TrimSpace(string(hash)),
			IsActive:     true,
			Role:         "user",
			Group:        "default",
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

		if payload.Email == "" || payload.NewPassword == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "email and new password are required"})
			return
		}

		user, err := getUserByUsername(db, payload.Email)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
			return
		}

		if user.PasswordHash != payload.OldPassword {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "old password is incorrect"})
			return
		}

		user.PasswordHash = payload.NewPassword
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

func getUserByUsername(db *gorm.DB, email string) (models.User, error) {
	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, err
		}
		return models.User{}, err
	}
	return user, nil
}

//func withInstanceHeader(next http.HandlerFunc) http.HandlerFunc {
//
//	var instanceID string
//	if v := os.Getenv("INSTANCE_ID"); v != "" {
//		instanceID = v
//	}
//	return func(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("X-Instance-Id", instanceID)
//		next(w, r)
//	}
//}

func HeartbeatHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		instanceID := os.Getenv("INSTANCE_ID")
		c.JSON(http.StatusOK, gin.H{"message": "heartbeat received", "instanceID": instanceID})
	}
}

func WithInstanceHeader(next gin.HandlerFunc) gin.HandlerFunc {

	var instanceID string
	if v := os.Getenv("INSTANCE_ID"); v != "" {
		instanceID = v
	}
	return func(c *gin.Context) {
		c.Header("X-Instance-Id", instanceID)
		next(c)
	}
}
