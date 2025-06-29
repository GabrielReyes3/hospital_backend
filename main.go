package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/joho/godotenv"
    "github.com/GabrielReyes3/hospital_backend/db"
    "github.com/GabrielReyes3/hospital_backend/handlers"
    "log"
    "os"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Println("⚠️ No se pudo cargar .env")
    }

    db.Connect()

    app := fiber.New()

    app.Get("/usuarios", handlers.ObtenerUsuarios)
	app.Post("/usuarios", handlers.CrearUsuario)


    port := os.Getenv("PORT")
    log.Fatal(app.Listen(":" + port))
}
