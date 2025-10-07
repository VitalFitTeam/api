# VitalFit API

Backend API for VitalFit, a gym management system, built with **Go**.

---

## Table of Contents

- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
  - [1. Clone the Repository](#1-clone-the-repository)
  - [2. Configure Environment Variables](#2-configure-environment-variables)
  - [3. Start the Database](#3-start-the-database)
  - [4. Install Dependencies](#4-install-dependencies)
  - [5. Run the Application](#5-run-the-application)
- [Live Reloading with Air](#live-reloading-with-air)
- [API Documentation](#api-documentation)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [License](#license)

---

## Prerequisites

Before you begin, ensure you have the following installed:

- **Go**: 1.25.1 or higher
- **Docker & Docker Compose**: For running the database
- **Git**: For cloning the repository
- **[Air](https://github.com/cosmtrek/air)**: For live reloading during development
- **direnv**: For loading environment variables from `.envrc`

---

## Getting Started

### 1. Clone the Repository

```sh
git clone https://github.com/yourusername/vitalfit-api.git
cd vitalfit-api
```

### 2. Configure Environment Variables

This project uses `.envrc` for environment variable management. **godotenv is not used**, so you must export variables manually or use [direnv](https://direnv.net/).

- Copy the example environment file:

  ```sh
  cp env.example .envrc
  ```

- Edit `.envrc` and set your environment variables as needed. Example on .env.example


- Allow direnv to load the variables:

  ```sh
  direnv allow .
  ```

### 3. Start the Database

Use Docker Compose to start the database service:

```sh
docker-compose up -d
```
or 

```sh
docker-compose up --build
```

This will start all services defined in [`docker-compose.yml`](docker-compose.yml).

### Migrations
This will start the migrations
```sh
  make migrate-up
```

### 4. Install Dependencies

Make sure Go modules are downloaded:

```sh
go mod tidy
```

### 5. Run the Application

You can run the application directly:

```sh
go run cmd/api/main.go
```

Or use Air for live reloading (recommended for development):

```sh
air
```

---

## Live Reloading with Air

Air watches for file changes and automatically restarts your Go application. To use Air:

1. Install Air:

   ```sh
   go install github.com/cosmtrek/air@latest
   ```

2. Run Air in the project root:

   ```sh
   air
   ```

Configuration is handled via `.air.toml` (already present in the project).

---

## API Documentation

API documentation is available in [docs/swagger.yaml](docs/swagger.yaml) and [docs/swagger.json](docs/swagger.json).

You can view the documentation using [Swagger UI](https://swagger.io/tools/swagger-ui/) or import the files into your preferred API client.

---

## Tech Stack

- **Go**: Backend language
- **PostgreSQL**: Database (via Docker Compose)
- **Air**: Live reloading for development
- **direnv**: Environment variable management
- **Docker Compose**: Service orchestration

---

## Project Structure



## License

This project is licensed under the Apache 2.0 License.