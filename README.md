# Due Draghi - Combattimenti

Calcolatore di incontri per QuintaEdizione. Supporta sia le regole 2024 (One D&D) che le regole 2014 (5e).

## Caratteristiche

- **Supporto Multi-Ruleset**: Calcola il budget XP per entrambe le edizioni di 5e
- **Modalità Gruppo Flessibile**: Gestisce gruppi con tutti i personaggi allo stesso livello o livelli diversi
- **Regole 2024**: Sistema di difficoltà semplificato (Bassa, Moderata, Alta)
- **Regole 2014**: Sistema di difficoltà classico (Facile, Media, Difficile, Letale) con moltiplicatori per numero di mostri
- **Ricerca Mostri**: Integrazione con quintaedizione.online per trovare mostri appropriati
- **UI Moderna**: Interfaccia stile Notion con HTMX per interazioni dinamiche

## Requisiti

- Go 1.25 o superiore
- Docker e Docker Compose (per deployment containerizzato)

## Installazione

### Sviluppo Locale

```bash
# Clona il repository
git clone https://github.com/emiliopalmerini/due-draghi-combattimenti.git
cd due-draghi-combattimenti

# Installa dipendenze
go mod download

# Genera template Templ
go install github.com/a-h/templ/cmd/templ@latest
templ generate

# Avvia l'applicazione
go run cmd/encounters/main.go
```

L'applicazione sarà disponibile su `http://localhost:8080`

### Docker

```bash
# Build e avvio con Docker Compose
docker-compose up -d

# Verifica lo stato
docker-compose ps
```

### Build Manuale

```bash
# Build del binario
go build -o bin/combattimenti cmd/encounters/main.go

# Esegui
./bin/combattimenti
```

## Utilizzo

1. Seleziona il ruleset (2024 o 2014)
2. Scegli la modalità gruppo:
   - **Stesso livello**: Tutti i personaggi hanno lo stesso livello
   - **Livelli diversi**: Ogni personaggio ha il proprio livello
3. Seleziona la difficoltà desiderata
4. Per le regole 2014: Specifica il numero di mostri per calcolare il moltiplicatore
5. Ottieni il budget XP totale per l'incontro

## Architettura

Il progetto segue i principi di Clean Architecture e Domain-Driven Design:

```
cmd/encounters/          - Entry point dell'applicazione
internal/
  ├── domain/           - Logica di business core
  │   └── encounter/    - Entità e value objects degli incontri
  ├── application/      - Use cases e servizi applicativi
  │   └── encounter/    - Servizi di calcolo XP e query
  └── infrastructure/   - Dettagli implementativi
      ├── persistence/  - Repository in-memory per dati XP
      ├── web/          - Handlers HTTP e template
      └── static/       - Asset CSS e JavaScript
```

## Test

```bash
# Esegui tutti i test
go test ./...

# Test con coverage
go test -cover ./...

# Test verbose
go test -v ./...
```

## Tecnologie

- **Backend**: Go 1.25 con chi router
- **Template**: Templ per template type-safe
- **Frontend**: HTMX per interattività dinamica
- **Container**: Docker multi-stage per build ottimizzate
- **Testing**: Go standard testing con table-driven tests

## API Endpoints

- `GET /` - Pagina principale del calcolatore
- `POST /calculate` - Calcola il budget XP dell'incontro
- `GET /party-input` - Ottieni opzioni per input del gruppo
- `GET /api/difficulties` - Ottieni difficoltà per ruleset
- `GET /health` - Health check
- `GET /ready` - Readiness check

## Contribuire

Le pull request sono benvenute. Per modifiche importanti, apri prima un issue per discutere cosa vorresti cambiare.

Assicurati di aggiornare i test appropriatamente.

## Licenza

MIT

## Contatti

Emilio Palmerini - [@emiliopalmerini](https://github.com/emiliopalmerini)

Link Progetto: [https://github.com/emiliopalmerini/due-draghi-combattimenti](https://github.com/emiliopalmerini/due-draghi-combattimenti)
