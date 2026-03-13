package main

import (
	"fmt"
	"log"
	"os"

	"sgi-tickets-back/handlers"
	"sgi-tickets-back/migrations"
	"sgi-tickets-back/storage"

	_ "sgi-tickets-back/docs" // Importar docs generados por swag

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	html "github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title SGI Tickets API
// @version 1.0
// @description API REST para Sistema de Gestión de Interventorías - Gestión de tickets de infraestructura
// @description
// @description **Autenticación**: Sistema basado en cookies con 2FA
// @description - Cookie `sgi_tickets_user_email`: Sesión de usuario
// @description - Cookie `sgi_tickets_identity`: Verificación 2FA
// @description
// @description **Roles disponibles**: admin, supervisor, agente, entidad, contratista
// @contact.name Soporte Técnico SGI
// @contact.email soporte@sgi.com
// @host localhost:8080
// @BasePath /api/v1
// @schemes http https
// @securityDefinitions.apikey CookieAuth
// @in cookie
// @name sgi_tickets_user_email
// @securityDefinitions.apikey TwoFactorAuth
// @in cookie
// @name sgi_tickets_identity

func main() {
	// Cargar .env si existe (opcional, para desarrollo local)
	// En Docker las variables vienen del docker-compose.yml
	_ = godotenv.Load()

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
		AllowOrigins:     "http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowCredentials: true,
	}))

	// Documentación Swagger (debe estar antes del SPA fallback)
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	app.Static("/", "./public", fiber.Static{Browse: false})

	// React SPA fallback
	app.Use(func(c *fiber.Ctx) error {
		path := c.Path()
		// Excluir /api y /swagger del SPA fallback
		if (len(path) >= 4 && path[:4] == "/api") || (len(path) >= 8 && path[:8] == "/swagger") {
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

	// Setup y verificación 2FA (requieren cookie sgi_tickets_user_email)
	auth2faRoutes := authRoutes.Group("/2fa", handlers.CookieMiddleware())
	auth2faRoutes.Get("/setup", handlers.Setup2FA)
	auth2faRoutes.Post("/setup", handlers.Setup2FA)
	auth2faRoutes.Post("/verify", handlers.Verify2FA)

	// Logout (requiere cookie sgi_tickets_user_email)
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

	// ==========================================
	// RUTAS DE USUARIOS
	// ==========================================
	usuariosRoutes := v1.Group("/usuarios", handlers.TwoFaMiddleware())
	usuariosRoutes.Get("/", handlers.ListarUsuarios)
	usuariosRoutes.Get("/:id", handlers.ObtenerUsuario)
	usuariosRoutes.Post("/", handlers.CrearUsuario)
	usuariosRoutes.Put("/:id", handlers.ActualizarUsuario)
	usuariosRoutes.Delete("/:id", handlers.EliminarUsuario)
	usuariosRoutes.Post("/:id/reset-password", handlers.ResetearPasswordUsuario)

	// ==========================================
	// RUTAS DE CATALOGOS
	// ==========================================
	catalogosRoutes := v1.Group("/catalogos", handlers.TwoFaMiddleware())
	catalogosRoutes.Get("/tipos-documento", handlers.GetTiposDocumentosIdentificacion)
	catalogosRoutes.Get("/regionales", handlers.GetRegionales)
	catalogosRoutes.Get("/departamentos", handlers.GetDepartamentos)
	catalogosRoutes.Get("/municipios", handlers.GetMunicipios)

	log.Fatal(app.Listen(":" + os.Getenv("SERVER_PORT")))
}
