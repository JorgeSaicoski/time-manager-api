package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/JorgeSaicoski/time-manager-api/internal/database"
	"github.com/JorgeSaicoski/time-manager-api/internal/models"
	"github.com/JorgeSaicoski/time-manager-api/internal/routes"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}
}

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.TotalTime{},
		&models.WorkTime{},
		&models.Project{},
		&models.Company{},
		&models.UserCompanyRole{},
		&models.Task{},
		&models.Cost{},
		&models.BreakTime{},
		&models.Brb{},
		&models.ProjectMember{},
		&models.ResolutionTracker{},
		&models.ResolutionUnit{},
	)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	cfg := &database.Config{
		DB: db,
	}

	router := routes.SetupRouter(cfg)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	serverAddr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on port %s", port)

	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
