# Matrix QR Analytics

Two independent REST APIs that communicate through HTTP to perform QR decomposition and calculate statistics from the resulting matrices.

## Architecture

- **QR Service** (`qr-service-go`): Go + Fiber — receives a matrix, validates input, computes QR decomposition, and aggregates the final response.
- **Statistics Service** (`statistics-service-node`): Node.js + Express — calculates statistics from matrices Q and R.

See [docs/architecture.md](docs/architecture.md) for the full architecture overview.

## Prerequisites

- Go 1.22+
- Node.js 18+
- Docker and Docker Compose (optional, for containerized execution)

## Project Structure

```text
matrix-qr-analytics/
├── qr-service-go/
│   ├── cmd/server/
│   ├── internal/{handlers,services,clients,models}/
│   └── configs/
├── statistics-service-node/
│   ├── src/{routes,controllers,services,models}/
│   └── config/
├── docker-compose.yml
└── .env.example
```

## Configuration

Copy the environment templates and adjust values as needed:

```bash
cp .env.example .env
cp qr-service-go/.env.example qr-service-go/.env
cp statistics-service-node/.env.example statistics-service-node/.env
```

| Variable | Service | Description | Default |
|----------|---------|-------------|---------|
| `PORT` | Both | HTTP port inside the container | `8080` / `3000` |
| `QR_SERVICE_PORT` | Docker Compose | Host port for QR service | `8080` |
| `STATISTICS_SERVICE_PORT` | Docker Compose | Host port for Statistics service | `3000` |
| `STATISTICS_SERVICE_URL` | QR Service | Base URL of the Statistics service | — |
| `STATISTICS_TIMEOUT_MS` | QR Service | HTTP timeout for Statistics calls (ms) | `5000` |

## Run Locally

### QR Service (Go)

```bash
cd qr-service-go
go run ./cmd/server
```

Health check: `GET http://localhost:8080/health`

### Statistics Service (Node.js)

```bash
cd statistics-service-node
npm install
npm start
```

Health check: `GET http://localhost:3000/health`

### Docker Compose

```bash
docker compose up --build
```

## API Endpoints

Business endpoints will be documented here as they are implemented:

| Service | Method | Path | Status |
|---------|--------|------|--------|
| QR Service | POST | `/api/v1/qr` | Planned |
| Statistics Service | POST | `/api/v1/statistics` | Planned |

## Documentation

- [Challenge specification](docs/challenge-spec.md)
- [Architecture](docs/architecture.md)
- [Assumptions](docs/assumptions.md)
