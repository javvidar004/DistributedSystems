package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"proyecto.com/m/controllers"
	"proyecto.com/m/models"
)

func main() {
	dsn := "host=db user=postgres password=mysecretpassword dbname=sampledb port=5432 sslmode=disable TimeZone=America/Mexico_City"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Office{}, &models.Workspace{}, &models.Booking{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/users", controllers.GetUsersHandler(db))
	router.GET("/users/:username", controllers.GetUserInfoHandler(db))
	router.POST("/login", controllers.LoginHandler(db))
	router.POST("/register", controllers.RegisterHandler(db))
	router.PUT("/update", controllers.UpdatePasswordHandler(db))
	router.DELETE("/delete/:username", controllers.DeleteUserHandler(db))
	router.GET("/offices", controllers.GetOfficesHandler(db))
	router.GET("/offices/:officeId/workspaces", controllers.GetOfficeWorkspacesHandler(db))
	router.GET("/bookings", controllers.GetBookingsHandler(db))
	router.POST("/bookings", controllers.CreateBookingHandler(db))

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
