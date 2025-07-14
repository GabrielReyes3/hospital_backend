package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/GabrielReyes3/hospital_backend/handlers"
)

func MedicoRoutes(router fiber.Router) {
	medico := router.Group("/medico")

	medico.Get("/citas", handlers.GetCitasMedico)
	medico.Get("/expedientes/:id", handlers.GetExpedientePaciente)
	medico.Post("/recetas", handlers.CrearRecetaMedico)
}
