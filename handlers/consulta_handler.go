package handlers

import (
	"context"
	"github.com/GabrielReyes3/hospital_backend/db"
	"github.com/GabrielReyes3/hospital_backend/models"
	"github.com/GabrielReyes3/hospital_backend/validators"
	"github.com/gofiber/fiber/v2"
)

func CrearConsulta(c *fiber.Ctx) error {
	var consulta models.CrearConsultaRequest

	if err := c.BodyParser(&consulta); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Datos inválidos",
			"error":   err.Error(),
		})
	}

	if err := validators.Validate.Struct(&consulta); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Validación fallida",
			"error":   err.Error(),
		})
	}

	query := `
        INSERT INTO consultas (id_paciente, id_medico, id_consultorio, tipo, horario)
        VALUES ($1, $2, $3, $4, $5)
    `
	_, err := db.Pool.Exec(context.Background(), query,
		consulta.IdPaciente,
		consulta.IdMedico,
		consulta.IdConsultorio,
		consulta.Tipo,
		consulta.Horario,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Error al crear consulta",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "ok",
		"message": "Consulta registrada con éxito",
	})
}
