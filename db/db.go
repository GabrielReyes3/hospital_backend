package db

import (
    "context"
    "log"
    "os"

    "github.com/jackc/pgx/v4/pgxpool"
)

var Pool *pgxpool.Pool

func Connect() error {
    connStr := os.Getenv("SUPABASE_CONN_STRING")

    config, err := pgxpool.ParseConfig(connStr)
    if err != nil {
        return err
    }

    // Deshabilita el uso de prepared statements automáticos para evitar el error
    config.ConnConfig.PreferSimpleProtocol = true

    Pool, err = pgxpool.ConnectConfig(context.Background(), config)
    if err != nil {
        return err
    }

    log.Println("✅ Conectado exitosamente a Supabase.")
    return nil
}
