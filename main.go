package main

import (
	"log"
	"os"
	"time"

	"github.com/GabrielReyes3/hospital_backend/db"
	"github.com/GabrielReyes3/hospital_backend/handlers"
	"github.com/GabrielReyes3/hospital_backend/middleware"
	"github.com/GabrielReyes3/hospital_backend/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ No se pudo cargar .env")
	}

	// Conexión a la base de datos
	db.Connect()

	// Inicializar Fiber
	app := fiber.New()

	// Middleware CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:4200",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Middleware global de rate limiting
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 10 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Demasiadas peticiones, intenta más tarde.",
			})
		},
	}))

	// 📌 Rutas públicas
	app.Post("/login", handlers.Login)
	app.Post("/register", handlers.CrearUsuario)
	app.Post("/activar-mfa", handlers.ActivarMFA)
	app.Post("/mfa/setup", middleware.RequireAuth(), handlers.MFASetup)

	app.Get("/consultorios", handlers.ObtenerConsultorios)
	app.Post("/consultas", handlers.CrearConsulta)

	app.Post("/refresh", handlers.RefreshToken)

	// 🔐 Rutas protegidas por JWT
	api := app.Group("/api", middleware.RequireAuth())

	api.Get("/usuarios", handlers.ObtenerUsuarios)
	api.Get("/paciente/historial", handlers.GetHistorialCitasPaciente)

	// 🧑‍⚕️ Panel enfermera
	routes.EnfermeraRoutes(api)

	// 👨‍⚕️ Panel médico
	routes.MedicoRoutes(api) // ✅ Agregado

	// Puerto
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Puerto por defecto
	}
	log.Fatal(app.Listen(":" + port))
}
