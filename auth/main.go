package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
				time.Sleep(60 * time.Second)
				continue
			}
			if resp.StatusCode != http.StatusOK {
				log.Printf("Backend %s is not alive", lbHost)
			}
			time.Sleep(60 * time.Second)
		}
	}()
}

func main() {

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	middlewareHost := os.Getenv("MIDDLEWARE_HOST")
	instanceID := os.Getenv("INSTANCE_ID")

	consulAddr := os.Getenv("CONSUL_HTTP_ADDR")

	if err := controllers.RegisterWithConsul(consulAddr, instanceID); err != nil {
		// Non-fatal: the service still runs and serves requests, it just
		// won't be discoverable via Consul until this succeeds. Logged
		// loudly because silently running "invisible" is a confusing
		// failure mode.
		log.Printf("WARNING: %v", err)
	}

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

	HeartbeatRequest(middlewareHost, instanceID)

	router.POST("/login", controllers.WithInstanceHeader(controllers.LoginHandler(db)))
	router.POST("/register", controllers.WithInstanceHeader(controllers.RegisterHandler(db)))
	router.PUT("/update", controllers.WithInstanceHeader(controllers.UpdatePasswordHandler(db)))
	router.DELETE("/delete/:userID", controllers.WithInstanceHeader(controllers.DeleteUserHandler(db)))
	router.GET("/heartbeat", controllers.WithInstanceHeader(controllers.HeartbeatHandler()))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Printf("shutting down %s", instanceID)
		controllers.DeregisterFromConsul(consulAddr, instanceID)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	}()

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
