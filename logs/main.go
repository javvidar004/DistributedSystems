package main

import (
	//"bytes"
	//"encoding/json"
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
	lbhost := os.Getenv("LB_HOST")
	instanceID := os.Getenv("INSTANCE_ID")

	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + name + " port=" + port + " sslmode=disable TimeZone=America/Mexico_City"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := db.AutoMigrate(&models.Log{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	HeartbeatRequest(lbhost, instanceID)

	router.GET("/log", controllers.WithInstanceHeader(controllers.GetLogHandler(db)))
	router.POST("/log", controllers.WithInstanceHeader(controllers.NewLogHandler(db)))
	router.GET("/heartbeat", controllers.WithInstanceHeader(controllers.HeartbeatHandler()))

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
