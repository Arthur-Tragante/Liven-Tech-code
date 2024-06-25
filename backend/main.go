package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/arthur-tragante/liven-code-test/controllers"
	"github.com/arthur-tragante/liven-code-test/models"
	"github.com/arthur-tragante/liven-code-test/routes"
	"github.com/arthur-tragante/liven-code-test/services"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	jwtSecret := os.Getenv("JWT_SECRET")

	dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " port=" + dbPort + " sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Address{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	userService := &services.UserService{DB: db, JWTSecret: jwtSecret}
	addressService := &services.AddressService{DB: db}

	userController := &controllers.UserController{UserService: userService}
	addressController := &controllers.AddressController{AddressService: addressService}

	r := gin.Default()
	routes.SetupRoutes(r, userController, addressController)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Definindo uma porta padrão
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
