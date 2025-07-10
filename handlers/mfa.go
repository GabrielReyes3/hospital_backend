package handlers

import (
    "github.com/gofiber/fiber/v2"
    "github.com/pquerna/otp/totp"
)

func MFASetup(c *fiber.Ctx) error {
    userID := c.Locals("userID").(string) // extrae userID del middleware JWT

    key, err := totp.Generate(totp.GenerateOpts{
        Issuer:      "HospitalApp",
        AccountName: userID,
    })
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Error generando clave MFA",
        })
    }

    // Usar la función para guardar en BD que definiste en usuario_handler.go
    if err := UpdateUserMFASecret(userID, key.Secret()); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Error guardando clave MFA",
        })
    }

    return c.JSON(fiber.Map{
        "otp_url": key.URL(),
        "secret":  key.Secret(),
    })
}
