package handlers

import (
	"context"
	"github.com/GabrielReyes3/hospital_backend/db"
	"github.com/gofiber/fiber/v2"
)

func ObtenerConsultorios(c *fiber.Ctx) error {
	rows, err := db.Pool.Query(context.Background(), `
        SELECT id, nombre, ubicacion, horario, tipo
        FROM consultorios
        WHERE status = 'activo'
    `)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error al obtener consultorios", "error": err.Error()})
	}
	defer rows.Close()

	var consultorios []map[string]interface{}

	for rows.Next() {
		var id int
		var nombre, ubicacion, horario, tipo string
		err := rows.Scan(&id, &nombre, &ubicacion, &horario, &tipo)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error al escanear consultorios", "error": err.Error()})
		}
		consultorios = append(consultorios, fiber.Map{
			"id":        id,
			"nombre":    nombre,
			"ubicacion": ubicacion,
			"horario":   horario,
			"tipo":      tipo,
		})
	}

	return c.JSON(fiber.Map{"status": "ok", "data": consultorios})
}
