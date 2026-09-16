package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"proyecto.com/m/controllers"
	"proyecto.com/m/models"
)

func HeartbeatRequest(lbHost string, instanceID string) {
	url := "http://" + lbHost + ":8080/heartbeat"
	type HeartbeatPayload struct {
		InstanceID string `json:"instanceID"`
	}
	go func() {
		for {
			resp, err := http.Get(url)
			if err != nil {
				log.Printf("Error checking heartbeat for %s: %v", lbHost, err)
				time.Sleep(20 * time.Second)
				continue
			}
			if resp.StatusCode != http.StatusOK {
				log.Printf("Backend %s is not alive", lbHost)
			}
			time.Sleep(20 * time.Second)
		}
	}()
}

func main() {

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	lbHost := os.Getenv("LB_HOST")
	instanceID := os.Getenv("INSTANCE_ID")

	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + name + " port=" + port + " sslmode=disable TimeZone=America/Mexico_City"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
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

	HeartbeatRequest(lbHost, instanceID)

	router.POST("/user", controllers.WithInstanceHeader(controllers.CreateUserHandler(db)))
	router.PUT("/user/:userID", controllers.WithInstanceHeader(controllers.UpdateUserHandler(db)))
	router.DELETE("/user/:userID", controllers.WithInstanceHeader(controllers.DeleteUserHandler(db)))
	router.GET("/user/:userID", controllers.WithInstanceHeader(controllers.GetUserHandler(db)))
	router.GET("/users", controllers.WithInstanceHeader(controllers.GetUsersHandler(db)))
	router.GET("/heartbeat", controllers.WithInstanceHeader(controllers.HeartbeatHandler()))

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
