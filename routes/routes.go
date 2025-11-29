package routes

import (
	"os"

	"github.com/barzaevhalid/cloud_storage_backend/handlers"
	"github.com/barzaevhalid/cloud_storage_backend/middleware"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, fileHandler *handlers.FileHandler, userHandler *handlers.UserHandler) {
	secret := os.Getenv("JWT_SECRET")

	api := app.Group("/api")

	//auth
	api.Post("/register", userHandler.Register)
	api.Post("/login", userHandler.Login)

	//refresh token
	app.Post("/refresh", userHandler.Refresh)

	protected := api.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: []byte(secret),
		},
		ContextKey: "user",
	}))
	protected.Use(middleware.JWTProtected())
	protected.Get("/me", userHandler.GetMe)
	//files
	files := protected.Group("/files")
	files.Post("/upload", fileHandler.Upload)
	files.Get("/", fileHandler.FindAllFiles)
	files.Delete("/:ids", fileHandler.DeleteFiles)
}
