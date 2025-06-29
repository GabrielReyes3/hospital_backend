package db

import (
    "context"
    "log"
    "os"
    "github.com/jackc/pgx/v4/pgxpool"
)

var Pool *pgxpool.Pool

func Connect() {
    connStr := os.Getenv("SUPABASE_CONN_STRING")
    var err error
    Pool, err = pgxpool.Connect(context.Background(), connStr)
    if err != nil {
        log.Fatalf("❌ Error conectando a Supabase: %v", err)
    }
    log.Println("✅ Conectado exitosamente a Supabase.")
}
