# Due Draghi - Combattimenti

D&D 5e encounter calculator supporting both 2014 and 2024 (One D&D) rules.

## Quick Start

```bash
# Build
templ generate && go build -o bin/combattimenti cmd/encounters/main.go

# Run (default port 8080)
./bin/combattimenti
```

## Development Commands

```bash
# Install templ
go install github.com/a-h/templ/cmd/templ@latest

# Generate templates
templ generate

# Run tests
go test ./...

# Format code
go fmt ./...

# Vet code
go vet ./...
```

## Architecture

Clean Architecture with Domain-Driven Design:

```
cmd/encounters/           - Application entry point
internal/
  domain/                 - Core business logic, entities, value objects
  application/            - Use cases and application services
  infrastructure/
    persistence/          - In-memory repositories
    web/
      handlers/           - HTTP handlers
      templates/          - Templ templates (.templ files)
    static/               - CSS and JavaScript assets
```

## Tech Stack

- Go 1.25
- chi router (github.com/go-chi/chi/v5)
- templ templates (github.com/a-h/templ)
- HTMX for dynamic interactions

## API Endpoints

- `GET /` - Main calculator page
- `POST /calculate` - Calculate encounter XP budget
- `GET /party-input` - Get party input options
- `GET /api/difficulties` - Get difficulties for ruleset
- `GET /health` - Health check
- `GET /ready` - Readiness check
