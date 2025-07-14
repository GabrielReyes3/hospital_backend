package models

import "database/sql"

type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
    TOTP     string `json:"totp"` // opcional si no tiene MFA activado
}

type TokenResponse struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
}

type RefreshTokenRequest struct {
    RefreshToken string `json:"refresh_token"`
}

type Usuario struct {
    ID              int            `json:"id"`
    Nombre          string         `json:"nombre"`
    Apellidos       string         `json:"apellidos"`
    Tipo            string         `json:"tipo"`
    FechaNacimiento sql.NullString `json:"fecha_nacimiento"`
    Genero          sql.NullString `json:"genero"`
    Correo          string         `json:"correo"`
    Contrasena      string         `json:"contrasena"`
    MfaSecret       sql.NullString `json:"mfa_secret"`
    RolID           int            `json:"rol_id"` // 👈 agregado
}
