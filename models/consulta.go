package models

type CrearConsultaRequest struct {
    IdPaciente   int    `json:"id_paciente" validate:"required"`
    IdMedico     int    `json:"id_medico" validate:"required"`
    IdConsultorio int   `json:"id_consultorio" validate:"required"`
    Tipo         string `json:"tipo" validate:"required"`
    Horario      string `json:"horario" validate:"required,datetime=2006-01-02T15:04:05Z07:00"` // RFC3339
}
