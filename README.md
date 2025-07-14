# 🏥 Hospital Backend

Este es el backend del sistema hospitalario desarrollado en **Go** con el framework **Fiber**, usando **Supabase (PostgreSQL)** como base de datos. Implementa seguridad robusta con autenticación JWT, autenticación multifactor (MFA), control de roles y permisos, logging inmutable, validaciones, y estructuras modulares por tipo de usuario.

---

## 🚀 Tecnologías Principales

- **Go 1.21+**
- **Fiber** – Framework web rápido y minimalista
- **Supabase (PostgreSQL)** – Base de datos relacional con conexión por string
- **JWT (Access + Refresh)**
- **MFA con TOTP** (Google Authenticator / Authy)
- **Rate Limiting y Logging**
- **Validación de esquemas con JSON Schema + validator**
- **Roles y permisos por usuario**

---

## 📁 Estructura del Proyecto

```bash
HOSPITAL_BACKEND/
├── db/
│   ├── db.go              # Conexión a Supabase
│   └── logger.go          # Logging inmutable
├── handlers/
│   ├── auth_handler.go    # Login, MFA, registro
│   ├── consulta_handler.go
│   ├── consultorio_handler.go
│   ├── enfermera_handler.go
│   ├── medico_handler.go
│   ├── paciente_handler.go
│   └── usuario_handler.go
├── middleware/
│   ├── auth.go            # Middleware RequireAuth y permisos
│   └── validate_schema.go # Middleware de validación JSON
├── models/
│   ├── consulta.go
│   ├── enfermera_models.go
│   ├── medico_models.go
│   └── usuario.go
├── routes/
│   ├── enfermera_routes.go
│   └── medico_routes.go
├── schemas/               # Archivos JSON Schema para validación
├── validators/            # Validadores personalizados
├── .env                   # Variables de entorno (no subir a GitHub)
├── .gitignore
├── go.mod / go.sum
├── main.go
└── README.md
🔐 Seguridad Implementada
✅ Autenticación con JWT (Access y Refresh)

✅ MFA (TOTP con QR)

✅ Middleware RequireAuth y control de permisos

✅ Rate Limiting

✅ Validación con go-playground/validator + esquemas JSON

✅ Contraseñas hasheadas (SHA-256)

✅ Logging inmutable de login y actividad

✅ Roles (admin, paciente, medico, enfermera) con tabla rol_permisos

📦 Instalación
Clona el repositorio:

bash

git clone https://github.com/GabrielReyes3/hospital_backend
cd hospital_backend
Crea un archivo .env con tu cadena de conexión:

env

SUPABASE_CONN_STRING=postgresql://usuario:contraseña@host:puerto/basededatos
JWT_SECRET=supersecreto
JWT_REFRESH_SECRET=refresh_supersecreto

Instala dependencias y ejecuta:
go mod tidy
go run main.go

📡 Endpoints Principales
🔐 Autenticación
Método	Ruta	Descripción
POST	/login	Inicio de sesión con o sin MFA
POST	/register	Registro de usuario
POST	/refresh	Refrescar token JWT
GET	/mfa/setup	Genera QR para activar MFA
POST	/mfa/verify	Verifica código MFA

👤 Usuarios
GET /usuarios – Listar usuarios

POST /usuarios – Crear usuario

PUT /usuarios/:id – Actualizar usuario

DELETE /usuarios/:id – Eliminar usuario

👨‍⚕️ Médico
GET /medico/citas

POST /medico/recetas

GET /medico/expedientes

🧑‍⚕️ Enfermera
GET /enfermera/citas

GET /enfermera/expedientes

👨 Paciente
POST /paciente/solicitar-cita

GET /paciente/historial-citas

GET /paciente/recetas

🧾 Base de Datos
Las tablas principales utilizadas:

usuarios (con rol_id)

roles / permisos / rol_permisos

citas

consultorios

recetas

expedientes

logins (histórico de intentos)

🧪 Validaciones y Esquemas
Las validaciones se manejan con:

go-playground/validator en structs

Esquemas JSON en schemas/

Middleware de validación con mensajes personalizados

✅ Roadmap
 Autenticación JWT

 Refresh token

 MFA con TOTP

 Logging persistente

 Control de roles y permisos

 CRUD de usuarios y citas

 Documentación con Swagger

 Pruebas unitarias automatizadas

👨‍💻 Autor
Gabriel Reyes
GitHub - GabrielReyes3

