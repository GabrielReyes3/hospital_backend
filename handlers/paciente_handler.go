package handlers

import (
	"fmt"
	"time"

	"github.com/GabrielReyes3/hospital_backend/db"
	"github.com/gofiber/fiber/v2"

	"strconv"
)

type ConsultaHistorial struct {
	ID              int       `json:"id"`
	Horario         time.Time `json:"horario"`
	Tipo            string    `json:"tipo"`
	Diagnostico     string    `json:"diagnostico"`
	MedicoNombre    string    `json:"medico_nombre"`
	MedicoApellidos string    `json:"medico_apellidos"`
}

func GetHistorialCitasPaciente(c *fiber.Ctx) error {
	userIDInterface := c.Locals("user_id")
	if userIDInterface == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "No autenticado",
		})
	}

	var userID int
	switch v := userIDInterface.(type) {
	case float64:
		userID = int(v)
	case int:
		userID = v
	case string:
		idInt, err := strconv.Atoi(v)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error interno con el ID de usuario",
			})
		}
		userID = idInt
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error interno con el ID de usuario",
		})
	}

	query := `
        SELECT 
            cons.id,
            cons.horario,
            cons.tipo,
            cons.diagnostico,
            u.nombre,
            u.apellidos
        FROM consultas cons
        JOIN usuarios u ON cons.id_medico = u.id
        WHERE cons.id_paciente = $1
        ORDER BY cons.horario DESC;
    `

	rows, err := db.Pool.Query(c.Context(), query, userID)
	if err != nil {
		fmt.Println("Error ejecutando consulta historial:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error en la consulta del historial",
		})
	}
	defer rows.Close()

	var historial []ConsultaHistorial
	for rows.Next() {
		var cHist ConsultaHistorial
		err := rows.Scan(&cHist.ID, &cHist.Horario, &cHist.Tipo, &cHist.Diagnostico, &cHist.MedicoNombre, &cHist.MedicoApellidos)
		if err != nil {
			fmt.Println("Error leyendo datos del historial:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error al leer datos del historial",
			})
		}
		historial = append(historial, cHist)
	}

	return c.JSON(fiber.Map{
		"status":    "success",
		"historial": historial,
	})
}
