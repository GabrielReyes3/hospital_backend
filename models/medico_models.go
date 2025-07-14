package models

import "time"

type CitaMedico struct {
	ID           int        `json:"id"`
	Tipo         string     `json:"tipo"`
	Horario      time.Time  `json:"horario"`
	Paciente     string     `json:"paciente"`
	PacienteID   int        `json:"paciente_id"`
	Consultorio  string     `json:"consultorio"`
	Diagnostico  *string    `json:"diagnostico,omitempty"`
}

type ExpedientePaciente struct {
	ID                     int       `json:"id"`
	PacienteID             int       `json:"paciente_id"`
	GrupoSanguineo         string    `json:"grupo_sanguineo"`
	Alergias               string    `json:"alergias"`
	EnfermedadesCronicas   string    `json:"enfermedades_cronicas"`
	AntecedentesFamiliares string    `json:"antecedentes_familiares"`
	AntecedentesPersonales string    `json:"antecedentes_personales"`
	MedicamentosHabituales string    `json:"medicamentos_habituales"`
	Vacunas                string    `json:"vacunas"`
	NotasGenerales         string    `json:"notas_generales"`
	FechaActualizacion     time.Time `json:"fecha_actualizacion"`
}

type CrearRecetaRequest struct {
	ConsultaID  int    `json:"consulta_id"`
	Medicamento string `json:"medicamento"`
	Dosis       string `json:"dosis"`
	Nota        string `json:"nota"`
}
