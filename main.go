package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/limiter"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/joho/godotenv"
    "github.com/GabrielReyes3/hospital_backend/db"
    "github.com/GabrielReyes3/hospital_backend/handlers"
    "github.com/GabrielReyes3/hospital_backend/middleware"
    "log"
    "os"
    "time"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Println("⚠️ No se pudo cargar .env")
    }

    db.Connect()

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

    // 📌 Rutas públicas con validación por JSON Schema
    app.Post("/login", handlers.Login)
    app.Post("/register", handlers.CrearUsuario)
    app.Post("/activar-mfa", handlers.ActivarMFA)
    app.Post("/mfa/setup", middleware.RequireAuth(), handlers.MFASetup)


    app.Get("/consultorios", handlers.ObtenerConsultorios)
    app.Post("/consultas", handlers.CrearConsulta)


    app.Post("/refresh", handlers.RefreshToken) // no necesita esquema si solo usa el token JWT

    // Rutas protegidas por token JWT
    api := app.Group("/api", middleware.RequireAuth())

    api.Get("/usuarios", handlers.ObtenerUsuarios)
    // Aquí puedes agregar otras rutas protegidas y aplicar el validador si son POST/PUT

    port := os.Getenv("PORT")
    log.Fatal(app.Listen(":" + port))
}
