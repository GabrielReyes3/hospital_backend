package handlers

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"regexp"
	"time"
	"github.com/GabrielReyes3/hospital_backend/db"
	"github.com/GabrielReyes3/hospital_backend/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"

)
var validate = validator.New()



// Hashea la contraseña usando SHA-256
func hashPassword(password string) string {
    hash := sha256.Sum256([]byte(password))
    return hex.EncodeToString(hash[:])
}

// Valida que tenga al menos 12 caracteres, un número y un símbolo
func isValidPassword(password string) bool {
    if len(password) < 12 {
        return false
    }
    hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
    hasSymbol := regexp.MustCompile(`[!@#~$%^&*()_+={}\[\]:;"'<>,.?\/\\|]`).MatchString(password)
    return hasNumber && hasSymbol
}
var (
    accessSecret  = []byte(os.Getenv("ACCESS_SECRET"))
    refreshSecret = []byte(os.Getenv("REFRESH_SECRET"))
)

// Genera un access token con duración corta (10 min)
func generateAccessToken(userID string) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(10 * time.Minute).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(accessSecret)
}

func RefreshToken(c *fiber.Ctx) error {
    var req models.RefreshTokenRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Solicitud inválida",
        })
    }

    token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
        // Validar que el método de firma sea el correcto
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fiber.NewError(fiber.StatusUnauthorized, "Método de firma inválido")
        }
        return refreshSecret, nil
    })

    if err != nil || !token.Valid {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Refresh token inválido",
        })
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || claims["user_id"] == nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Token inválido",
        })
    }

    userID := fmt.Sprintf("%v", claims["user_id"])
    newAccessToken, err := generateAccessToken(userID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Error generando nuevo access token",
        })
    }

    return c.JSON(fiber.Map{
        "access_token": newAccessToken,
    })
}


func ActivarMFA(c *fiber.Ctx) error {
    correo := c.Query("correo") // o por body
    if correo == "" {
        return c.Status(400).JSON(fiber.Map{"error": "Falta el correo"})
    }

    // Generar clave secreta TOTP
    key, err := totp.Generate(totp.GenerateOpts{
        Issuer:      "HospitalApp",
        AccountName: correo,
    })
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Error generando clave"})
    }

    // Guardar en la base de datos
    query := `UPDATE usuarios SET mfa_secret = $1 WHERE correo = $2`
    _, err = db.Pool.Exec(c.Context(), query, key.Secret(), correo)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Error guardando clave"})
    }

    // Generar código QR como PNG base64
    png, err := qrcode.Encode(key.URL(), qrcode.Medium, 256)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Error generando QR"})
    }

    qrBase64 := base64.StdEncoding.EncodeToString(png)

    return c.JSON(fiber.Map{
        "otpauth_url": key.URL(), // opcional
        "qr_base64":   qrBase64,
    })
}

// Genera un refresh token con duración más larga (7 días)
func generateRefreshToken(userID string) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(refreshSecret)
}



func obtenerRolID(tipo string) (int, error) {
    switch tipo {
    case "paciente":
        return 1, nil
    case "medico":
        return 2, nil
    case "enfermera":
        return 3, nil
    default:
        return 0, fmt.Errorf("tipo de usuario inválido: %s", tipo)
    }
}

func CrearUsuario(c *fiber.Ctx) error {
    var input UsuarioInput

    if err := c.BodyParser(&input); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "JSON inválido"})
    }

    if err := validate.Struct(input); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error":   "Validación fallida",
            "detalle": err.Error(),
        })
    }

    if !isValidPassword(input.Contrasena) {
        return c.Status(400).JSON(fiber.Map{
            "error": "La contraseña debe tener mínimo 12 caracteres, un número y un símbolo",
        })
    }

    rolID, err := obtenerRolID(input.Tipo)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    usuario := models.Usuario{
        Nombre:          input.Nombre,
        Apellidos:       input.Apellidos,
        Tipo:            input.Tipo,
        FechaNacimiento: toNullString(input.FechaNacimiento),
        Genero:          toNullString(input.Genero),
        Correo:          input.Correo,
        Contrasena:      hashPassword(input.Contrasena),
        RolID:           rolID,
    }

    query := `
        INSERT INTO usuarios (nombre, apellidos, tipo, fecha_nacimiento, genero, correo, contrasena, rol_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id
    `

    err = db.Pool.QueryRow(
        c.Context(),
        query,
        usuario.Nombre,
        usuario.Apellidos,
        usuario.Tipo,
        usuario.FechaNacimiento,
        usuario.Genero,
        usuario.Correo,
        usuario.Contrasena,
        usuario.RolID,
    ).Scan(&usuario.ID)

    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.Status(201).JSON(usuario)
}



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
    Nombre          string `json:"nombre" validate:"required"`
    Apellidos       string `json:"apellidos" validate:"required"`
    Tipo            string `json:"tipo" validate:"required"`
    FechaNacimiento string `json:"fecha_nacimiento"`
    Genero          string `json:"genero"`
    Correo          string `json:"correo" validate:"required,email"`
    Contrasena      string `json:"contrasena" validate:"required,min=12"`
}

func toNullString(s string) sql.NullString {
    if s == "" {
        return sql.NullString{Valid: false}
    }
    return sql.NullString{String: s, Valid: true}
}






func Login(c *fiber.Ctx) error {
    var req models.LoginRequest

    if err := c.BodyParser(&req); err != nil {
        _ = db.RegistrarLog(db.Pool, nil, "login", false, "Body inválido", c.IP(), c.Get("User-Agent"))
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Datos inválidos",
        })
    }

    if err := validate.Struct(req); err != nil {
        _ = db.RegistrarLog(db.Pool, nil, "login", false, "Validación fallida: "+err.Error(), c.IP(), c.Get("User-Agent"))
        return c.Status(400).JSON(fiber.Map{
            "error":   "Validación fallida",
            "detalle": err.Error(),
        })
    }

    var usuario struct {
        ID       int
        Correo   string
        Tipo     string
        RolID    int
        Contrasena string
        MfaSecret sql.NullString
    }

    query := "SELECT id, correo, tipo, rol_id, contrasena, mfa_secret FROM usuarios WHERE correo = $1"
    err := db.Pool.QueryRow(c.Context(), query, req.Email).Scan(
        &usuario.ID,
        &usuario.Correo,
        &usuario.Tipo,
        &usuario.RolID,
        &usuario.Contrasena,
        &usuario.MfaSecret,
    )

    if err == sql.ErrNoRows {
        _ = db.RegistrarLog(db.Pool, nil, "login", false, "Correo no encontrado: "+req.Email, c.IP(), c.Get("User-Agent"))
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Correo o contraseña incorrectos",
        })
    } else if err != nil {
        _ = db.RegistrarLog(db.Pool, nil, "login", false, "Error de servidor al buscar usuario", c.IP(), c.Get("User-Agent"))
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Error de servidor",
        })
    }

    hashedInput := hashPassword(req.Password)
    if hashedInput != usuario.Contrasena {
        _ = db.RegistrarLog(db.Pool, &usuario.ID, "login", false, "Contraseña incorrecta", c.IP(), c.Get("User-Agent"))
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Correo o contraseña incorrectos",
        })
    }

    if usuario.MfaSecret.Valid {
        valid := totp.Validate(req.TOTP, usuario.MfaSecret.String)
        if !valid {
            _ = db.RegistrarLog(db.Pool, &usuario.ID, "login", false, "Código MFA incorrecto", c.IP(), c.Get("User-Agent"))
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Código MFA incorrecto",
            })
        }
    }

    accessToken, err := generateAccessTokenWithRole(fmt.Sprint(usuario.ID), usuario.Tipo, usuario.RolID)
    if err != nil {
        _ = db.RegistrarLog(db.Pool, &usuario.ID, "login", false, "Error generando access token", c.IP(), c.Get("User-Agent"))
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Error generando access token",
        })
    }

    refreshToken, err := generateRefreshToken(fmt.Sprint(usuario.ID))
    if err != nil {
        _ = db.RegistrarLog(db.Pool, &usuario.ID, "login", false, "Error generando refresh token", c.IP(), c.Get("User-Agent"))
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Error generando refresh token",
        })
    }

    _ = db.RegistrarLog(db.Pool, &usuario.ID, "login", true, "Inicio de sesión exitoso", c.IP(), c.Get("User-Agent"))

    return c.JSON(fiber.Map{
        "access_token":  accessToken,
        "refresh_token": refreshToken,
        "usuario": fiber.Map{
            "id":      usuario.ID,
            "correo":  usuario.Correo,
            "tipo":    usuario.Tipo,
            "rol_id":  usuario.RolID,
        },
    })
}

// Ajusta la función para generar access token con ambos claims
func generateAccessTokenWithRole(userID, tipo string, rolID int) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "tipo":    tipo,
        "rol_id":  rolID,
        "exp":     time.Now().Add(10 * time.Minute).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(accessSecret)
}


// Actualiza el secret MFA en la tabla usuarios para el usuario dado
func UpdateUserMFASecret(userID, secret string) error {
    ctx := context.Background()
    sql := `UPDATE usuarios SET mfa_secret = $1 WHERE id = $2`

    cmdTag, err := db.Pool.Exec(ctx, sql, secret, userID) // asumiendo db.Pool es *pgxpool.Pool
    if err != nil {
        return err
    }
    if cmdTag.RowsAffected() != 1 {
        return fmt.Errorf("no se actualizó ningún usuario con id %s", userID)
    }
    return nil
}