# Repository Guidelines

## Project Structure & Module Organization
- `main.tsp`, `service.tsp`, `models/`, and `operations/` define the API contract in TypeSpec.
- `tspconfig.yaml` configures OpenAPI generation into `openapi/`.
- `cmd/api/` contains the Go entrypoint.
- `internal/config/`, `internal/domain/`, `internal/service/`, `internal/httpapi/`, and `internal/repository/postgres/` hold configuration, domain logic, HTTP handlers, and persistence.
- `migrations/` contains PostgreSQL schema changes.
- Tests live next to the code they cover, for example `internal/service/service_test.go`.

## Build, Test, and Development Commands
- `go test ./...` runs the full Go test suite.
- `go run ./cmd/api` starts the HTTP API locally.
- `gofmt -w <files>` formats Go sources; run it before committing.
- `tsp compile main.tsp` generates OpenAPI from TypeSpec when the TypeSpec CLI is installed.
- `go mod tidy` updates module metadata after dependency changes.
- `make up` starts PostgreSQL, applies migrations, and launches the API through Compose.
- `make migrate-up` and `make migrate-down` manage the database schema through Compose.
- `make openapi` generates `openapi/openapi.yaml` through Compose and TypeSpec.

## Coding Style & Naming Conventions
- Use standard Go formatting and idioms; keep imports grouped and let `gofmt` manage layout.
- Use short, descriptive package names and export only what must be used across packages.
- Name HTTP handlers by action, such as `listBookings` or `createBooking`.
- Keep TypeSpec names aligned with domain language: `EventType`, `Booking`, `Slot`, `PublicEventType`.

## Testing Guidelines
- Prefer table-driven tests for service and handler behavior.
- Name tests `TestXxx` and keep them close to the code under test.
- Cover validation, booking conflicts, and the 14-day availability window.
- Run `go test ./...` before opening a pull request.

## Commit & Pull Request Guidelines
- Keep commit messages short and imperative, matching the existing history style, for example `add booking conflict check`.
- In pull requests, summarize the API or schema impact, note any new migrations, and mention generated OpenAPI changes.
- Link related issues when available and include request/response examples if HTTP behavior changes.

## Security & Configuration Tips
- Do not commit local secrets; use environment variables such as `DATABASE_URL`, `LISTEN_ADDR`, and `TIMEZONE`.
- Booking conflicts are enforced both in service logic and in the PostgreSQL schema; keep those rules in sync.
