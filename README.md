# Hospital Backend

Sistema backend para gestión hospitalaria desarrollado en Go usando el framework Fiber y Supabase como base de datos PostgreSQL.

## Tecnologías
- Go (Fiber)
- Supabase (PostgreSQL)
- JWT (futuro)
- Docker (futuro)

## Estructura del proyecto
📁 db
  📄 db.go – conexión a Supabase
📁 handlers
  📄 usuario_handler.go – lógica de endpoints para usuarios
📁 models
  📄 usuario.go – definición del modelo de usuario
📄 main.go – arranque del servidor
📄 .env – variables de entorno
📄 go.mod / go.sum – gestión de dependencias
📄 CHANGELOG.md – registro de cambios

## Uso

```bash
go run main.go