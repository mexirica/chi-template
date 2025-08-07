<div align="center">
  <img src="assets/template.png" alt="Logo" width="120">
</div>
A modern, production-ready Go REST API template with full observability, Docker Compose orchestration, and best practices for configuration, documentation, and testing.

## Features

- **Go 1.24** with idiomatic project structure
- **Chi** router for fast, composable HTTP routing
- **PostgreSQL** with SQLC for type-safe queries
- **Redis** for caching
- **OpenTelemetry** for distributed tracing and metrics
- **Prometheus, Grafana, Loki, Tempo, cAdvisor, Node Exporter** for full observability (metrics, logs, traces)
- **Swagger (Swaggo)** for API documentation
- **Viper** for configuration management
- **GoMock** for service and repository testing
- **Docker Compose** for local development and orchestration
- **Provisioned Grafana dashboards and datasources**
- **GoDoc** and Swagger comments for all public APIs and endpoints

## Structure

- `cmd/` — Application entrypoint
- `internal/` — Main application code (handlers, services, repositories, models, middleware, etc)
- `deploy/` — Docker Compose, Grafana, Prometheus, Loki, Tempo, and other infra configs
- `docs/` — Swagger/OpenAPI docs

## Getting Started

1. **Clone the repository:**
   ```sh
   git clone <repo-url>
   cd chi-template
   ```
2. **Copy and edit the `.env` file:**
   ```sh
   cp .env.example .env
   # Edit as needed
   ```
3. **Install the dependencies:**
   ```sh
   make setup
   ```
4. **Build and run with Docker Compose:**
   ```sh
   make compose_up
   ```

## Development

- **Run tests:**
  ```sh
  make test
  ```
- **Generate Swagger docs:**
  ```sh
  make docgen
  ```
- **Hot reload (optional):**
  Use [Air](https://github.com/cosmtrek/air) for live reload during development.

## What’s Included

- CRUD endpoints for movies (with SQLC and repository/service pattern)
- Full observability stack (metrics, logs, traces) wired to Grafana
- Example Grafana dashboards for Docker, host, and application metrics
- API documentation with Swaggo/Swagger
- GoDoc comments for all packages and public functions
- Example GoMock usage for service and repository tests
- Docker Compose for all services and infrastructure

## How to Extend

- Add new endpoints in `internal/handler/` and `internal/service/`
- Add new SQLC queries in `internal/db/queries/` and regenerate with `sqlc generate`
- Add new Grafana dashboards in `deploy/grafana/provisioning/dashboards/`
- Add new Prometheus scrape configs in `deploy/prometheus/prometheus.yml`
