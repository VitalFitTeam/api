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

    This command will build and start the application containers along with the database. You will be able to see the live-reloading logs.

    ```bash
    make docker-up
    ```

    If you want to run it in the background:

    ```bash
    docker compose up -d
    ```

    The API will be available at the address specified in your `.env` file (default is `localhost:8080`). The service uses `Air` for hot-reloading when it detects code changes.

2.  **Stop the services:**

    To stop and remove the containers, run:

    ```bash
    make docker-down
    ```

## Makefile Commands

The `Makefile` provides several commands to facilitate common development tasks.

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

The server will start and listen on the port specified in your `.env` file (default: `http://localhost:8080`).

## Testing

The project is configured with a CI workflow that runs several quality checks. To run these checks locally:

```bash
# Run unit and integration tests
go test ./...

# Check for suspicious code
go vet ./...

# Run static analysis (requires prior installation)
# go install honnef.co/go/tools/cmd/staticcheck@latest
staticcheck ./...
```

## API Documentation

This project uses `swaggo/swag` to generate API documentation in Swagger/OpenAPI format.

*   **Access the documentation**: Once the server is running, you can view the interactive API docs in your browser at:
    `http://localhost:8080/swagger/index.html`

*   **Update the documentation**: If you make changes to the API comments (the `// @...` annotations), you should regenerate the documentation. Run the following command from the project root:

    ```bash
    swag init -g cmd/api/main.go
    ```

## Environment Variables

Below are the environment variables used by the application:

| Variable            | Description                                                               | Example                                                        |
| ------------------- | ------------------------------------------------------------------------- | -------------------------------------------------------------- |
| `ADDR`          | Port where the server will run.                                           | `8080`                                                         |
| `ENV`           | Application environment (`development`, `staging`, `production`).         | `development`                                                  |
| `DB_ADDR`            | Data Source Name for connecting to PostgreSQL.                            | `"host=localhost user=postgres ..."`                           |
| `JWT_SECRET`        | Secret key used to sign JWT tokens.                                       | `"a-very-long-hard-to-guess-secret"`                           |
| `JWT_EXPIRATION`    | Duration that JWT tokens are valid.                                       | `"72h"`                                                        |
| `RESEND_API_KEY`    | API key for the Resend email service.                                     | `"re_xxxxxxxxxxxx"`                                            |
| `RESEND_FROM_EMAIL` | Sender email address used by Resend.                                      | `"noreply@yourdomain.com"`                                     |

## License

This project is licensed under the Apache 2.0 License. See the `LICENSE` file for details.