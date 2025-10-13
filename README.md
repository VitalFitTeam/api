# VitalFit API

Quick start guide to set up and run the VitalFit backend using Docker and Make.

## Table of Contents

- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Configuration](#configuration)
- [Running the Project](#running-the-project)
- [Makefile Commands](#makefile-commands)

## Getting Started

Follow these instructions to get the project up and running on your local machine.

### Prerequisites

- **Docker and Docker Compose**: To run the application in containers.
- **Make**: To use the simplified commands from the `Makefile`.
- **Git**: To clone the repository.

### Configuration

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/vitalfit/api.git
    cd vitalfit-api
    ```

2.  **Create and configure your `.env` file:**

    Copy the example file `.env.example` to a new file named `.env`.

    ```bash
    cp .env.example .env
    ```

    Open the `.env` file and fill in the variables with your local credentials and settings (database, API keys, etc.).

## Running the Project

1.  **Start the services with Docker Compose:**

    This command will build and start the application containers along with the database, and you will be able to see the live-reloading logs.

    ```bash
    make docker-up
    ```

    If you want to run it in the background:

    ```bash
    docker compose up -d
    ```

    The API will be available at the address specified in your `.env` file (default is `localhost:8080`). The service uses `Air` for hot-reloading when it detects changes in the code.

2.  **Stop the services:**

    To stop and remove the containers, run:

    ```bash
    make docker-down
    ```

## Makefile Commands

The `makefile` provides several commands to facilitate common development tasks.

### Database Migrations

*   **Create a new migration:**

    ```bash
    make migration migration_name
    ```

*   **Apply all pending migrations:**

    ```bash
    make migrate-up
    ```

*   **Revert the last migration:**

    ```bash
    make migrate-down
    ```

### Database Seeding

*   **Seed the database with initial data (e.g., an admin user):**

    ```bash
    make seed
    ```

### API Documentation

*   **Generate or update Swagger documentation:**

    This command parses the code comments and updates the files in the `/docs` directory.

```bash
make gen-docs
```

El servidor se iniciará y estará escuchando en el puerto especificado en tu archivo `.env` (por defecto, `http://localhost:8080`).

## Testing

El proyecto está configurado con un flujo de trabajo de CI que ejecuta varias comprobaciones de calidad. Para ejecutar estas comprobaciones localmente:

```bash
# Ejecutar pruebas unitarias y de integración
go test ./...

# Verificar que no haya código "sospechoso"
go vet ./...

# Ejecutar análisis estático (requiere instalación previa)
# go install honnef.co/go/tools/cmd/staticcheck@latest
staticcheck ./...
```

## Documentación de la API

Este proyecto utiliza `swaggo/swag` para generar documentación de la API en formato Swagger/OpenAPI.

*   **Acceder a la documentación**: Una vez que el servidor esté en ejecución, puedes ver la documentación interactiva de la API en tu navegador visitando:
    `http://localhost:8080/swagger/index.html`

*   **Actualizar la documentación**: Si realizas cambios en los comentarios de la API (las anotaciones `// @...`), debes regenerar los archivos de documentación. Ejecuta el siguiente comando desde la raíz del proyecto:

    ```bash
    swag init -g cmd/api/main.go
    ```

## Variables de Entorno

A continuación se detallan las variables de entorno utilizadas por la aplicación:

| Variable            | Descripción                                                              | Ejemplo                                                        |
| ------------------- | ------------------------------------------------------------------------ | -------------------------------------------------------------- |
| `API_PORT`          | Puerto en el que se ejecutará el servidor.                               | `8080`                                                         |
| `API_ENV`           | Entorno de la aplicación (`development`, `staging`, `production`).       | `development`                                                  |
| `DB_DSN`            | Data Source Name para la conexión a PostgreSQL.                          | `"host=localhost user=postgres ..."`                           |
| `JWT_SECRET`        | Clave secreta para firmar los tokens JWT.                                | `"un-secreto-muy-largo-y-dificil-de-adivinar"`                 |
| `JWT_EXPIRATION`    | Duración de validez de los tokens JWT.                                   | `"72h"`                                                        |
| `RESEND_API_KEY`    | Clave API para el servicio de correo Resend.                             | `"re_xxxxxxxxxxxx"`                                            |
| `RESEND_FROM_EMAIL` | Dirección de correo electrónico remitente para Resend.                   | `"noreply@tudominio.com"`                                      |
                                     


## Licencia

Este proyecto está bajo la Licencia Apache 2.0. Consulta el archivo `LICENSE` para más detalles.