package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"proyecto.com/m/conectors"
	"proyecto.com/m/models"
)

func GetOfficesHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var offices []models.Office
		if err := db.Order("id asc").Find(&offices).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to load offices", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, offices)
	}
}

func GetOfficeWorkspacesHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		officeID, err := strconv.Atoi(c.Param("officeId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid office id"})
			return
		}

		var workspaces []models.Workspace
		if err := db.Where("office_id = ?", officeID).Order("id asc").Find(&workspaces).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to load workspaces", "error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, workspaces)
	}
}

func GetBookingsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var bookings []models.Booking
		query := db.Preload("User").Preload("Workspace").Where("booking_date >= ?", time.Now().UTC().Truncate(24*time.Hour))

		if username := c.Query("username"); username != "" {
			query = query.Joins("JOIN users ON users.id = bookings.user_id").Where("users.username = ?", username)
		}

		if err := query.Order("booking_date asc").Find(&bookings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to load bookings", "error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, bookings)
	}
}

func CreateBookingHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload conectors.BookingRequest
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid booking payload"})
			return
		}

		bookingDate, err := time.Parse("2006-01-02", payload.BookingDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "booking_date must be YYYY-MM-DD"})
			return
		}

		booking := models.Booking{
			UserID:      uint(payload.UserID),
			WorkspaceID: uint(payload.WorkspaceID),
			BookingDate: bookingDate,
		}

		if err := db.Create(&booking).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not create booking", "error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "booking created successfully", "booking": booking})
	}
}
