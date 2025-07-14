package handlers

import (
	"context"

	"github.com/GabrielReyes3/hospital_backend/db"
	"github.com/GabrielReyes3/hospital_backend/models"
	"github.com/gofiber/fiber/v2"
)



func GetCitasMedico(c *fiber.Ctx) error {
	rows, err := db.Pool.Query(context.Background(), `
		SELECT 
			con.id,  -- agregar id de consulta
			u.nombre || ' ' || u.apellidos AS paciente,
			co.tipo, 
			con.horario::TEXT, 
			co.nombre AS consultorio,
			COALESCE(con.diagnostico, '') AS diagnostico
		FROM consultas con
		JOIN usuarios u ON u.id = con.id_paciente
		JOIN consultorios co ON co.id = con.id_consultorio
		WHERE con.id_medico = $1
		ORDER BY con.horario DESC
	`, c.Locals("user_id").(int)) // asumiendo que usas auth y user_id en locals

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al obtener citas: " + err.Error()})
	}
	defer rows.Close()

	var citas []map[string]interface{}
	for rows.Next() {
		var id int
		var paciente, tipo, horario, consultorio, diagnostico string
		if err := rows.Scan(&id, &paciente, &tipo, &horario, &consultorio, &diagnostico); err != nil {
			continue
		}
		citas = append(citas, fiber.Map{
			"id":          id,
			"paciente":    paciente,
			"tipo":        tipo,
			"horario":     horario,
			"consultorio": consultorio,
			"diagnostico": diagnostico,
		})
	}

	return c.JSON(citas)
}




func GetExpedientePaciente(c *fiber.Ctx) error {
	pacienteID := c.Params("id")

	row := db.Pool.QueryRow(context.Background(), `
		SELECT id, paciente_id, grupo_sanguineo, alergias, enfermedades_cronicas, 
		       antecedentes_familiares, antecedentes_personales, medicamentos_habituales, 
		       vacunas, notas_generales, fecha_actualizacion
		FROM expedientes
		WHERE paciente_id = $1
	`, pacienteID)

	var exp models.ExpedientePaciente
	if err := row.Scan(
		&exp.ID, &exp.PacienteID, &exp.GrupoSanguineo, &exp.Alergias, &exp.EnfermedadesCronicas,
		&exp.AntecedentesFamiliares, &exp.AntecedentesPersonales, &exp.MedicamentosHabituales,
		&exp.Vacunas, &exp.NotasGenerales, &exp.FechaActualizacion); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Expediente no encontrado"})
	}

	return c.JSON(exp)
}



func CrearRecetaMedico(c *fiber.Ctx) error {
	conn := db.Pool
	medicoID := c.Locals("user_id").(int)

	var data struct {
		ConsultaID int    `json:"consultaID"`
		Medicamento string `json:"medicamento"`
		Dosis string       `json:"dosis"`
		Nota string       `json:"nota"`
	}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	_, err := conn.Exec(c.Context(), `
		INSERT INTO recetas (id_consulta, fecha, id_medico, medicamento, dosis, nota)
		VALUES ($1, CURRENT_DATE, $2, $3, $4, $5)
	`, data.ConsultaID, medicoID, data.Medicamento, data.Dosis, data.Nota)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "No se pudo crear la receta"})
	}

	// Obtener nombre del paciente asociado a la consulta
	var paciente string
	err = conn.QueryRow(c.Context(), `
		SELECT u.nombre || ' ' || u.apellidos
		FROM consultas con
		JOIN usuarios u ON con.id_paciente = u.id
		WHERE con.id = $1
	`, data.ConsultaID).Scan(&paciente)

	if err != nil {
		paciente = "Desconocido"
	}

	return c.JSON(fiber.Map{
		"status":   "Receta creada",
		"paciente": paciente,
	})
}
