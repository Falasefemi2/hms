# Hospital Management System

A Hospital Management System built with Go and PostgreSQL.

## Setup

1. Clone the repository
2. Copy `.env.example` to `.env`
3. Fill in your database credentials
4. Run `go mod download`
5. Run `go run cmd/api/main.go`

## Project Structure

- `cmd/api/` - Application entry point
- `internal/` - Internal packages
  - `config/` - Configuration
  - `database/` - Database connection
  - `models/` - Data structures
  - `handlers/` - HTTP handlers
  - `middleware/` - HTTP middleware

## Environment Variables

See `.env.example` for required variables.
