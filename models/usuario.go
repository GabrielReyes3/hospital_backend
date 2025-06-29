package models

import "database/sql"

type Usuario struct {
    ID              int            `json:"id"`
    Nombre          string         `json:"nombre"`
    Apellidos       string         `json:"apellidos"`
    Tipo            string         `json:"tipo"`
    FechaNacimiento sql.NullString `json:"fecha_nacimiento"`
    Genero          sql.NullString `json:"genero"`
    Correo          string         `json:"correo"`
    Contrasena      string         `json:"contrasena"`
}
