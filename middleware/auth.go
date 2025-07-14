package middleware

import (
    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v5"
    "os"
    "strconv"
)

var accessSecret = []byte(os.Getenv("ACCESS_SECRET"))

func RequireAuth() fiber.Handler {
    return func(c *fiber.Ctx) error {
        tokenString := c.Get("Authorization")
        if tokenString == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Token no proporcionado",
            })
        }

        // Soporta formato: Bearer <token>
        if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
            tokenString = tokenString[7:]
        }

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return accessSecret, nil
        })

        if err != nil || !token.Valid {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Token inválido o expirado",
            })
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok || claims["user_id"] == nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Token inválido",
            })
        }

        // Convertir user_id a int para mantener consistencia
        var userID int
        switch v := claims["user_id"].(type) {
        case float64:
            userID = int(v)
        case string:
            id, err := strconv.Atoi(v)
            if err != nil {
                return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                    "error": "user_id inválido en token",
                })
            }
            userID = id
        default:
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Tipo inválido para user_id",
            })
        }

        // Guardar user_id como int en el contexto
        c.Locals("user_id", userID)

        return c.Next()
    }
}
