package main

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/thomasdarmawan9/datastream-backend/services/microB/internal/domain"
	"github.com/thomasdarmawan9/datastream-backend/services/microB/internal/infrastructure/auth"
	grpcInfra "github.com/thomasdarmawan9/datastream-backend/services/microB/internal/infrastructure/grpc"
	mysqlRepo "github.com/thomasdarmawan9/datastream-backend/services/microB/internal/infrastructure/mysql"
	"github.com/thomasdarmawan9/datastream-backend/services/microB/internal/interfaces/http"
	"github.com/thomasdarmawan9/datastream-backend/services/microB/internal/interfaces/middleware"
	"github.com/thomasdarmawan9/datastream-backend/services/microB/internal/usecase"

	echoSwagger "github.com/swaggo/echo-swagger"
	_ "github.com/thomasdarmawan9/datastream-backend/services/microB/docs"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	sensorpb "github.com/thomasdarmawan9/datastream-backend/proto/sensorpb"
)

// @title Microservice B API
// @version 1.0
// @description This is the backend service for receiving and managing sensor data.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@datastream.local

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api
func main() {
	_ = godotenv.Load()

	// --- Load ENV ---
	dsn := os.Getenv("DB_DSN") // contoh: root@tcp(localhost:3306)/datastream?parseTime=true
	if dsn == "" {
		log.Fatal("DB_DSN not set")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET not set")
	}
	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	// --- DB Init ---
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect db: ", err)
	}
	if err := db.AutoMigrate(&domain.User{}, &domain.SensorData{}); err != nil {
		log.Fatal("failed migrate: ", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("failed to get sql.DB from gorm.DB: ", err)
	}

	// --- Repository ---
	userRepo := mysqlRepo.NewUserRepository(sqlDB)
	sensorRepo := mysqlRepo.NewSensorRepository(sqlDB)

	// --- Usecase ---
	userUC := usecase.NewUserUsecase(userRepo)
	sensorUC := usecase.NewSensorUsecase(sensorRepo)
	jwtExpiry := 24 * time.Hour
	jwtManager := auth.NewJWTManager(jwtSecret, jwtExpiry)

	// --- Start gRPC Server ---
	go func() {
		lis, err := net.Listen("tcp", ":"+grpcPort)
		if err != nil {
			log.Fatalf("failed to listen on gRPC port %s: %v", grpcPort, err)
		}
		grpcServer := grpc.NewServer()
		sensorpb.RegisterSensorServiceServer(grpcServer, grpcInfra.NewSensorGRPCServer(sensorRepo))
		log.Println("Microservice B gRPC server running at :" + grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// --- Start HTTP Server (Echo) ---
	e := echo.New()
	e.Use(echomw.Logger())
	e.Use(echomw.Recover())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Public routes
	http.NewUserHandler(e, userUC, jwtManager)

	// Protected routes
	api := e.Group("/api")
	api.Use(middleware.JWTAuth(jwtManager, "admin", "user"))
	http.NewSensorHandler(api, sensorUC)

	log.Println("Microservice B HTTP server running at :" + httpPort)
	if err := e.Start(":" + httpPort); err != nil {
		log.Fatal(err)
	}
}
