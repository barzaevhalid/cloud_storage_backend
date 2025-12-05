package main

import (
	"context"
	"fmt"
	"log"

	"github.com/barzaevhalid/cloud_storage_backend/config"
	"github.com/barzaevhalid/cloud_storage_backend/db"
	"github.com/barzaevhalid/cloud_storage_backend/handlers"
	"github.com/barzaevhalid/cloud_storage_backend/repositories"
	"github.com/barzaevhalid/cloud_storage_backend/routes"
	"github.com/barzaevhalid/cloud_storage_backend/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"

	_ "github.com/barzaevhalid/cloud_storage_backend/docs"
)

// @host localhost:8080
// @BasePath /
func main() {
	//config
	cfg := config.LoadConfig()

	//context for bd
	ctx := context.Background()

	//conection to bd
	pool, err := db.NewPool(ctx, *cfg)

	if err != nil {
		log.Fatalf("Unable to conne ct to database: %v\n", err)
	}

	defer pool.Close()
	fmt.Println("Database connected")

	//repositores
	fileRepo := repositories.NewFileRepository(pool)
	userRepo := repositories.NewUserRepository(pool)
	//services
	fileService := services.NewFileService(fileRepo)
	userService := services.NewUserService(userRepo, cfg.Auth.Secret, cfg.Auth.AccessTokenMin, cfg.Auth.RefreshTokenDays)
	//handlers
	userHandler := handlers.NewUserHandler(userService)
	fileHandler := handlers.NewFileHandler(fileService)

	app := fiber.New()
	app.Use(cors.New())
	app.Static("/uploads", "./uploads")
	app.Get("/swagger/*", swagger.HandlerDefault)

	//routes
	routes.SetupRoutes(app, fileHandler, userHandler)

	//start server
	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}
