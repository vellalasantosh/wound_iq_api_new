# Wound_IQ API (Gin + PostgreSQL)

This repository contains a production-ready REST API in Go (Gin) that connects to a local PostgreSQL database named `wound_iq`. It uses `database/sql` with the `pgx` driver.

## Features

- CRUD for Patients, Clinicians, Assessments
- Pagination and filtering for list endpoints
- Read endpoints that call existing DB functions:
  - `add_patient(full_name, date_of_birth, gender, medical_record_number)`
  - `add_full_assessment(...)`
  - `get_assessment_full(assessment_id)`
  - `get_patient_wound_history(patient_id)`
  - `get_all_patients()`
  - `get_all_assessments()`
- Graceful shutdown, structured logging (zap), request validation, error handling
- Unit test examples and Makefile
- OpenAPI spec (`openapi.yaml`)

---

## Prerequisites

- Go 1.22+
- PostgreSQL database `wound_iq` with schema/functions already created
- Git & GitHub account (your username: `vellalasantosh`)
- VS Code (already set up)

---

## Quick start (local)

1. Clone the repo:
```bash
git clone git@github.com:vellalasantosh/wound_iq_api_new.git
cd wound_iq_api_new
```

2. Copy environment file and update:
```bash
cp .env.example .env
# edit DB_DSN to point to your local DB credentials
```

3. Download dependencies and run:
```bash
go mod tidy
make run
```

4. Example create patient:
```bash
curl -X POST http://localhost:8080/v1/patients \
  -H "Content-Type: application/json" \
  -d '{"full_name":"John Doe","date_of_birth":"1970-01-01T00:00:00Z","gender":"male","medical_record_number":"MRN123"}'
```

---

## Files & Structure

```
cmd/api/main.go
internal/config/config.go
internal/db/db.go
internal/logger/logger.go
internal/models/
internal/handlers/
internal/router/router.go
openapi.yaml
Makefile
.env.example
README.md
```

---

## Tests

Run unit tests:
```bash
make test
```

Note: Unit tests in this repo use `sqlmock` placeholders in the examples; for full coverage add `github.com/DATA-DOG/go-sqlmock` and write table-driven tests.

---

## GitHub

Create a repository `wound_iq_api_new` under your GitHub account and push:
```bash
git init
git add .
git commit -m "initial: wound_iq_api_new"
git branch -M main
git remote add origin git@github.com:vellalasantosh/wound_iq_api_new.git
git push -u origin main
```

---

## Evolving with DB changes

1. Add migration with `golang-migrate`.
2. Update `internal/models/*.go` and handlers to reflect new columns.
3. Add tests and run CI.

For detailed developer notes see the comments in source files.
