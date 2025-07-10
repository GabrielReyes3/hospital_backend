package middleware

import (
    "github.com/gofiber/fiber/v2"
    "github.com/xeipuuv/gojsonschema"
    "os"
)

func ValidateSchema(schemaPath string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        var body map[string]interface{}
        if err := c.BodyParser(&body); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "error": "Cuerpo JSON inválido",
            })
        }

        // Ruta absoluta al esquema
        wd, _ := os.Getwd()
        fullPath := "file://" + wd + "/" + schemaPath

        schemaLoader := gojsonschema.NewReferenceLoader(fullPath)
        documentLoader := gojsonschema.NewGoLoader(body)

        result, err := gojsonschema.Validate(schemaLoader, documentLoader)
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": "Error al validar schema",
            })
        }

        if !result.Valid() {
            var errs []string
            for _, e := range result.Errors() {
                errs = append(errs, e.String())
            }
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "error":  "Validación fallida",
                "detail": errs,
            })
        }

        // Reescribimos el body para que el handler lo pueda usar
        c.Locals("validatedBody", body)
        return c.Next()
    }
}
