package handlers

import (
	"context"
	"log"

	"github.com/GabrielReyes3/hospital_backend/db"
	"github.com/GabrielReyes3/hospital_backend/models"
	"github.com/gofiber/fiber/v2"
)

func GetCitasEnfermera(c *fiber.Ctx) error {
    query := `
    SELECT 
      consultas.id,
      consultorios.nombre AS consultorio,
      CONCAT(medico.nombre, ' ', medico.apellidos) AS medico,
      CONCAT(paciente.nombre, ' ', paciente.apellidos) AS paciente,
      consultas.tipo,
      consultas.horario,
      consultas.diagnostico
    FROM consultas
    JOIN consultorios ON consultas.id_consultorio = consultorios.id
    JOIN usuarios AS medico ON consultas.id_medico = medico.id
    JOIN usuarios AS paciente ON consultas.id_paciente = paciente.id
    ORDER BY consultas.horario DESC
    LIMIT 50;
    `

log.Println("Ejecutando consulta GetCitasEnfermera")
rows, err := db.Pool.Query(context.Background(), query)
if err != nil {
    log.Println("Error en consulta GetCitasEnfermera:", err)
    return c.Status(500).JSON(fiber.Map{"error": "Error al obtener citas: " + err.Error()})
}


    defer rows.Close()

    var citas []models.Cita
    for rows.Next() {
        var cita models.Cita
        err := rows.Scan(&cita.ID, &cita.Consultorio, &cita.Medico, &cita.Paciente, &cita.Tipo, &cita.Horario, &cita.Diagnostico)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Error al leer datos de citas"})
        }
        citas = append(citas, cita)
    }

    return c.JSON(citas)
}

func GetExpedientesEnfermera(c *fiber.Ctx) error {
    query := `
    SELECT
      expedientes.id,
      expedientes.paciente_id,
      CONCAT(usuarios.nombre, ' ', usuarios.apellidos) AS paciente_nombre,
      expedientes.grupo_sanguineo,
      expedientes.alergias,
      expedientes.enfermedades_cronicas,
      expedientes.antecedentes_familiares,
      expedientes.antecedentes_personales,
      expedientes.medicamentos_habituales,
      expedientes.vacunas,
      expedientes.notas_generales,
      expedientes.fecha_actualizacion
    FROM expedientes
    JOIN usuarios ON expedientes.paciente_id = usuarios.id
    ORDER BY expedientes.fecha_actualizacion DESC
    LIMIT 50;
    `

    rows, err := db.Pool.Query(context.Background(), query)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Error al obtener expedientes"})
    }
    defer rows.Close()

    var expedientes []models.Expediente
    for rows.Next() {
        var exp models.Expediente
        err := rows.Scan(
            &exp.ID,
            &exp.PacienteID,
            &exp.PacienteNombre,
            &exp.GrupoSanguineo,
            &exp.Alergias,
            &exp.EnfermedadesCronicas,
            &exp.AntecedentesFamiliares,
            &exp.AntecedentesPersonales,
            &exp.MedicamentosHabituales,
            &exp.Vacunas,
            &exp.NotasGenerales,
            &exp.FechaActualizacion,
        )
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Error al leer datos de expedientes"})
        }
        expedientes = append(expedientes, exp)
    }

    return c.JSON(expedientes)
}
