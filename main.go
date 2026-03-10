package main

import (
	"fmt"
	"log"
	"os"

	"sgi-tickets-back/handlers"
	"sgi-tickets-back/migrations"
	"sgi-tickets-back/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	html "github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error cargando archivo .env")
	}

	storage.DBConnection()

	// Flag --rollback: revierte la ultima migracion y sale
	for _, arg := range os.Args[1:] {
		if arg == "--rollback" {
			if err := migrations.RollbackMigration(storage.DB); err != nil {
				log.Fatal(err)
			}
			fmt.Println("Rollback completado. Saliendo...")
			os.Exit(0)
		}
	}

	templatesEngine := html.New("./templates", ".html")

	app := fiber.New(fiber.Config{
		Views:     templatesEngine,
		BodyLimit: 100 * 1024 * 1024, // 100MB
	})

	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: os.Getenv("SECRET_KEY"),
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowCredentials: true,
	}))

	app.Static("/", "./public", fiber.Static{Browse: false})

	// React SPA fallback
	app.Use(func(c *fiber.Ctx) error {
		if len(c.Path()) >= 4 && c.Path()[:4] == "/api" {
			return c.Next()
		}
		return c.SendFile("./public/index.html")
	})

	apiRoutes := app.Group("/api")
	v1 := apiRoutes.Group("/v1")

	// Version publica
	apiRoutes.Get("/version", handlers.GetVersion)

	// ==========================================
	// RUTAS DE AUTENTICACIÓN
	// ==========================================
	authRoutes := v1.Group("/auth")

	// Login y recuperación (sin autenticación)
	authRoutes.Post("/login", handlers.Login)
	authRoutes.Post("/recover", handlers.RecoverPassword)
	authRoutes.Post("/reset", handlers.ResetPassword)

	// Setup y verificación 2FA (requieren cookie sgi_user_email)
	auth2faRoutes := authRoutes.Group("/2fa", handlers.CookieMiddleware())
	auth2faRoutes.Get("/setup", handlers.Setup2FA)
	auth2faRoutes.Post("/setup", handlers.Setup2FA)
	auth2faRoutes.Post("/verify", handlers.Verify2FA)

	// Logout (requiere cookie sgi_user_email)
	authRoutes.Post("/logout", handlers.CookieMiddleware(), handlers.Logout)

	// ==========================================
	// RUTAS DE PERFIL
	// ==========================================
	perfilRoutes := v1.Group("/perfil", handlers.TwoFaMiddleware())
	perfilRoutes.Get("/", handlers.GetPerfil)
	perfilRoutes.Put("/", handlers.UpdatePerfil)
	perfilRoutes.Put("/password", handlers.ChangePassword)
	perfilRoutes.Get("/2fa/activar", handlers.Activar2FA)
	perfilRoutes.Post("/2fa/activar", handlers.Activar2FA)
	perfilRoutes.Post("/2fa/desactivar", handlers.Desactivar2FA)

	log.Fatal(app.Listen(":" + os.Getenv("SERVER_PORT")))
}
