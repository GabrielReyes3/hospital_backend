package models

import "time"

type Cita struct {
    ID          int        `json:"id"`
    Consultorio string     `json:"consultorio"`
    Medico      string     `json:"medico"`
    Paciente    string     `json:"paciente"`
    Tipo        string     `json:"tipo"`
    Horario     time.Time  `json:"horario"`
    Diagnostico *string    `json:"diagnostico"`  // Puntero para NULL
}


type Expediente struct {
    ID                     int        `json:"id"`
    PacienteID             int        `json:"paciente_id"`
    PacienteNombre         string     `json:"paciente_nombre"`
    GrupoSanguineo         string     `json:"grupo_sanguineo"`
    Alergias               *string    `json:"alergias"`
    EnfermedadesCronicas   *string    `json:"enfermedades_cronicas"`
    AntecedentesFamiliares *string    `json:"antecedentes_familiares"`
    AntecedentesPersonales *string    `json:"antecedentes_personales"`
    MedicamentosHabituales *string    `json:"medicamentos_habituales"`
    Vacunas                *string    `json:"vacunas"`
    NotasGenerales         *string    `json:"notas_generales"`
    FechaActualizacion     string     `json:"fecha_actualizacion"`
}

