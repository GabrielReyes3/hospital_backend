package routes

import (
    "github.com/gofiber/fiber/v2"
    "github.com/GabrielReyes3/hospital_backend/handlers"
    "github.com/GabrielReyes3/hospital_backend/middleware"
)

func EnfermeraRoutes(app fiber.Router) {
    enfermera := app.Group("/enfermera", middleware.RequireAuth())

    enfermera.Get("/citas", handlers.GetCitasEnfermera)
    enfermera.Get("/expedientes", handlers.GetExpedientesEnfermera)
}
