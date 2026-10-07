# Chirpy

Chirpy is a backend service built with Go. It provides an HTTP API for user accounts, authentication, chirps, and server-side data management.

## Features

- Create and manage users
- User authentication
- Create, retrieve, and delete chirps
- HTTP API built with Go
- PostgreSQL database integration
- JSON request and response handling
- Protected endpoints using authentication

## Tech Stack

- Go
- net/http
- PostgreSQL
- SQL
- JSON
- HTTP REST API

## Getting Started

### Prerequisites

- Go installed
- PostgreSQL installed and running
- A configured environment file

### Installation

```bash
git clone https://github.com/OrationKnight21/Chirpy.git
cd chirpy
go mod download
```

### **Running server**
```bash
go run .
```
The server will start on the configured port.

## API EXAMPLES
### Creating a user
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}'
```
### Get Chirps
```bash 
curl http://localhost:8080/api/chirps
```
## Project Structure
```text
.
├── main.go
├── internal/
├── migrations/
├── sql/
├── go.mod
└── README.md
```
## Testing 
```bash
go test ./...
```
## Future Improvements
Add docker support
Add rate limiting
Add frontend client (hopefully)