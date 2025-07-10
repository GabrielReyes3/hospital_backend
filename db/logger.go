package db

import (
    "context"
    "time"

    "github.com/jackc/pgx/v4/pgxpool"
)

func RegistrarLog(pool *pgxpool.Pool, usuarioID *int, accion string, exito bool, mensaje string, ip string, userAgent string) error {
    _, err := pool.Exec(
        context.Background(),
        `INSERT INTO logs (usuario_id, accion, exito, mensaje, ip, user_agent, timestamp)
         VALUES ($1, $2, $3, $4, $5, $6, $7)`,
        usuarioID, accion, exito, mensaje, ip, userAgent, time.Now(),
    )
    return err
}
