package handlers

import (

	"database/sql"
    "github.com/gofiber/fiber/v2"
    "github.com/GabrielReyes3/hospital_backend/db"
    "github.com/GabrielReyes3/hospital_backend/models"
)

func ObtenerUsuarios(c *fiber.Ctx) error {
    rows, err := db.Pool.Query(c.Context(), "SELECT id, nombre, apellidos, tipo, fecha_nacimiento, genero, correo, contrasena FROM usuarios")
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    defer rows.Close()

    var usuarios []models.Usuario
    for rows.Next() {
        var u models.Usuario
        err := rows.Scan(&u.ID, &u.Nombre, &u.Apellidos, &u.Tipo, &u.FechaNacimiento, &u.Genero, &u.Correo, &u.Contrasena)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": err.Error()})
        }
        usuarios = append(usuarios, u)
    }

    return c.JSON(usuarios)

}

type UsuarioInput struct {
    Nombre          string `json:"nombre"`
    Apellidos       string `json:"apellidos"`
    Tipo            string `json:"tipo"`
    FechaNacimiento string `json:"fecha_nacimiento"`
    Genero          string `json:"genero"`
    Correo          string `json:"correo"`
    Contrasena      string `json:"contrasena"`
}

func toNullString(s string) sql.NullString {
    if s == "" {
        return sql.NullString{Valid: false}
    }
    return sql.NullString{String: s, Valid: true}
}

func CrearUsuario(c *fiber.Ctx) error {
    var input UsuarioInput

    if err := c.BodyParser(&input); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "JSON inválido"})
    }

    usuario := models.Usuario{
        Nombre:          input.Nombre,
        Apellidos:       input.Apellidos,
        Tipo:            input.Tipo,
        FechaNacimiento: toNullString(input.FechaNacimiento),
        Genero:          toNullString(input.Genero),
        Correo:          input.Correo,
        Contrasena:      input.Contrasena,
    }

    query := `
        INSERT INTO usuarios (nombre, apellidos, tipo, fecha_nacimiento, genero, correo, contrasena)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id
    `

    err := db.Pool.QueryRow(
        c.Context(),
        query,
        usuario.Nombre,
        usuario.Apellidos,
        usuario.Tipo,
        usuario.FechaNacimiento.String,
        usuario.Genero.String,
        usuario.Correo,
        usuario.Contrasena,
    ).Scan(&usuario.ID)

    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.Status(201).JSON(usuario)
}