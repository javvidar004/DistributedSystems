package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"proyecto.com/m/controllers"
	"proyecto.com/m/grpcserver"
	"proyecto.com/m/models"
	"proyecto.com/m/pb"
	"proyecto.com/m/service"
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
			resp.Body.Close()
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
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

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

	authSvc := service.NewAuthService(db)

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

	router.POST("/login", controllers.WithInstanceHeader(controllers.LoginHandler(authSvc)))
	router.POST("/register", controllers.WithInstanceHeader(controllers.RegisterHandler(authSvc)))
	router.PUT("/update", controllers.WithInstanceHeader(controllers.UpdatePasswordHandler(authSvc)))
	router.DELETE("/delete/:userID", controllers.WithInstanceHeader(controllers.DeleteUserHandler(authSvc)))
	router.GET("/heartbeat", controllers.WithInstanceHeader(controllers.HeartbeatHandler()))

	srv := &http.Server{
		Addr:    ":" + httpPort,
		Handler: router,
	}

	grpcLis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatal("failed to listen for gRPC:", err)
	}
	grpcSrv := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcSrv, grpcserver.NewAuthServer(authSvc))

	errCh := make(chan error, 2)
	var once sync.Once
	reportErr := func(err error) {
		if err == nil {
			return
		}
		once.Do(func() {
			errCh <- err
		})
	}

	go func() {
		log.Printf("HTTP server listening on :%s", httpPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			reportErr(err)
		}
	}()

	go func() {
		log.Printf("gRPC server listening on :%s", grpcPort)
		if err := grpcSrv.Serve(grpcLis); err != nil {
			reportErr(err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Printf("shutting down %s", instanceID)
		controllers.DeregisterFromConsul(consulAddr, instanceID)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("error shutting down HTTP server: %v", err)
		}
		grpcSrv.GracefulStop()
		reportErr(http.ErrServerClosed)
	}()

	if err := <-errCh; err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("server stopped with error:", err)
	}
}
