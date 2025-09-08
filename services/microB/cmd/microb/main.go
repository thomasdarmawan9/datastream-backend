package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/thomasdarmawan9/datastream-backend/services/microB/internal/domain"
	"github.com/thomasdarmawan9/datastream-backend/services/microB/internal/infrastructure/auth"
	mysqlRepo "github.com/thomasdarmawan9/datastream-backend/services/microB/internal/infrastructure/mysql"
	"github.com/thomasdarmawan9/datastream-backend/services/microB/internal/interfaces/http"
	"github.com/thomasdarmawan9/datastream-backend/services/microB/internal/interfaces/middleware"
	"github.com/thomasdarmawan9/datastream-backend/services/microB/internal/usecase"

	"time"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	echoSwagger "github.com/swaggo/echo-swagger"
	_ "github.com/thomasdarmawan9/datastream-backend/services/microB/docs"
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
	dsn := os.Getenv("DB_DSN") // contoh: root:password@tcp(mysql:3306)/datastream?parseTime=true
	if dsn == "" {
		log.Fatal("DB_DSN not set")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET not set")
	}

	// --- DB Init ---
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect db: ", err)
	}

	// Auto migrate
	if err := db.AutoMigrate(&domain.User{}, &domain.SensorData{}); err != nil {
		log.Fatal("failed migrate: ", err)
	}

	// --- Repository ---
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("failed to get sql.DB from gorm.DB: ", err)
	}
	userRepo := mysqlRepo.NewUserRepository(sqlDB)
	sensorRepo := mysqlRepo.NewSensorRepository(sqlDB)

	// --- Usecase ---
	userUC := usecase.NewUserUsecase(userRepo)
	sensorUC := usecase.NewSensorUsecase(sensorRepo)

	jwtExpiry := 24 * time.Hour // Set JWT expiry duration to 24 hours
	jwtManager := auth.NewJWTManager(jwtSecret, jwtExpiry)

	// --- Echo Server ---
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

	// --- Start ---
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Microservice B running at :" + port)
	if err := e.Start(":" + port); err != nil {
		log.Fatal(err)
	}
}
