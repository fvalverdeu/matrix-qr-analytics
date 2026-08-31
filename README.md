# Matrix QR Analytics

Matrix QR Analytics is a two-service backend solution for QR decomposition and derived matrix statistics.

Flow:

1. A client sends a rectangular numeric matrix to the Go QR Service.
2. The Go service validates input and computes QR decomposition.
3. The Go service sends Q and R to the Node Statistics Service over HTTP.
4. The Node service computes aggregate statistics.
5. The Go service returns Q, R, and statistics to the client.

The normal public entry point is the Go service.

## Architecture

```text
Client
	|
	v
Go QR Service :8080
	|
	| HTTP
	v
Node Statistics Service :3000
```

Responsibilities:

- Go QR Service (`qr-service-go`): public facade/orchestration, matrix validation, QR decomposition, Node call, aggregated response.
- Node Statistics Service (`statistics-service-node`): validates Q/R payloads and computes statistics.

In local Docker Compose, both services are reachable on localhost for development/debugging.

## Technology Stack

- Go QR Service: Go 1.24+, Fiber, Gonum.
- Statistics Service: Node.js 22+, TypeScript, Express 4.x.
- Infrastructure: Docker, Docker Compose.
- API documentation: OpenAPI 3.1, Scalar static reference.
- API testing: Postman collection.
- Testing: Go built-in testing package; Node built-in `node:test`.
- Static analysis: `go vet`.

## Prerequisites

### Docker-first (recommended)

- Docker
- Docker Compose

### Local development (without Docker)

- Go 1.24+
- Node.js 22+
- npm

## Configuration

Configuration sources in this repository:

- Root Compose env template: `.env.example`
- Go local template: `qr-service-go/.env.example`
- Node local template: `statistics-service-node/.env.example`

Environment variables currently used:

| Variable | Used by | Purpose | Default | Required |
|---|---|---|---|---|
| `QR_SERVICE_PORT` | Docker Compose | Host port mapped to Go service container port 8080 | `8080` | No |
| `STATISTICS_SERVICE_PORT` | Docker Compose | Host port mapped to Node service container port 3000 | `3000` | No |
| `STATISTICS_TIMEOUT_MS` | Go service | Timeout for Go -> Node HTTP call (milliseconds) | `5000` | No |
| `STATISTICS_SERVICE_URL` | Go service | Base URL for Node statistics endpoint | none in code | Yes for direct local Go run |
| `PORT` | Go and Node processes | Service listen port | Go: `8080`, Node: `3000` | No |

Important env behavior:

- Docker Compose reads root `.env` values when present.
- Go does not auto-load `.env` files; it reads process environment variables.
- Node loads dotenv in `src/index.ts`, so running from `statistics-service-node` can use `statistics-service-node/.env`.

## Running with Docker Compose

From repository root:

```bash
docker compose up --build
```

Run detached if preferred:

```bash
docker compose up --build -d
```

Stop and remove containers/network:

```bash
docker compose down
```

Default local ports:

- Go QR Service: `http://localhost:8080`
- Node Statistics Service: `http://localhost:3000`

Quick checks:

```bash
curl http://localhost:8080/health
curl http://localhost:3000/health
```

## Running Locally

### Statistics Service

```bash
cd statistics-service-node
npm ci
npm run build
npm start
```

By default, the service listens on port `3000`.

### QR Service

Start the Statistics Service first, then run Go:

Direct local Go execution must know how to reach the Statistics Service.
Set `STATISTICS_SERVICE_URL` before starting the Go process.
Set `STATISTICS_TIMEOUT_MS` optionally to override the default timeout.

Unix-like shells:

```bash
cd qr-service-go
go mod download
export STATISTICS_SERVICE_URL=http://localhost:3000
# optional (default is 5000ms)
export STATISTICS_TIMEOUT_MS=5000
go run ./cmd/server
```

PowerShell:

```powershell
cd qr-service-go
go mod download
$env:STATISTICS_SERVICE_URL="http://localhost:3000"
# optional (default is 5000ms)
$env:STATISTICS_TIMEOUT_MS="5000"
go run ./cmd/server
```

## API Endpoints

| Service | Method | Endpoint | Purpose | Visibility |
|---|---|---|---|---|
| QR Service | GET | `/health` | Service health check | Public |
| QR Service | POST | `/api/v1/qr` | Validate matrix, compute QR, aggregate statistics | Public (primary entry point) |
| Statistics Service | GET | `/health` | Service health check | Local debug/internal |
| Statistics Service | POST | `/api/v1/statistics` | Compute statistics from Q and R | Internal service-to-service (locally reachable for debug/testing) |

### Primary API Example

Request to public operation (`m < n` wide matrix):

```bash
curl -X POST http://localhost:8080/api/v1/qr \
	-H "Content-Type: application/json" \
	-d '{"matrix":[[1,2,3],[4,5,6]]}'
```

The response contains `q`, `r`, and `statistics`.
See `docs/openapi.yaml` or `docs/api-reference.html` for the exact response contract.

## API Documentation

### OpenAPI

- Source of truth: `docs/openapi.yaml`
- Format: OpenAPI 3.1

### Scalar API Reference

- Static viewer: `docs/api-reference.html`
- Renders `docs/openapi.yaml` (no runtime integration in Go/Node services)
- Must be served over HTTP; opening with `file://` may fail because the browser cannot fetch `./openapi.yaml` from local file origin policies.
- Current Scalar viewer script is CDN-hosted; internet access is required to load the Scalar frontend asset.

One convenient local option (if Python is available):

```bash
cd docs
python -m http.server 8000
```

Then open: `http://localhost:8000/api-reference.html`

### Postman

- Collection: `postman/Matrix-QR-Analytics.postman_collection.json`
- Includes both services and collection variables (`qrBaseUrl`, `statisticsBaseUrl`)
- Includes health, success, invalid matrix, malformed JSON, and square/tall/wide scenarios with assertions

Import the collection into Postman and run requests using the default local variables.

## Testing

### Go Tests

```bash
cd qr-service-go
go test ./...
go vet ./...
```

### Node Tests

```bash
cd statistics-service-node
npm ci
npm test
```

`npm test` builds TypeScript and runs tests with built-in `node:test`.

### Integration Tests (Go -> Node)

Prerequisites:

1. Node.js is available in PATH.
2. Statistics TypeScript output is built.

Build Node artifacts first:

```bash
cd statistics-service-node
npm ci
npm run build
```

Then run integration tests:

```bash
cd ../qr-service-go
go test -tags=integration ./integration
```

The integration bootstrap starts a Node test server process from compiled output.

## Matrix and Statistics Semantics

Matrix/QR semantics:

- Input matrix `A` has shape `m x n` and must be rectangular with finite numeric values.
- Output `Q` shape is `m x m`.
- Output `R` shape is `m x n`.
- Square, tall, and wide matrices are supported.
- Numerical contract: `A ~= Q * R`, `Q^T * Q ~= I`, and `R` is upper trapezoidal.

Statistics semantics:

- Statistics are computed over all numeric cells across both `Q` and `R`.
- `average = sum / total number of cells in Q and R`.
- `hasDiagonalMatrix` is `true` when `Q` OR `R` is diagonal.
- A diagonal matrix must be square.
- Off-diagonal absolute values `<= 1e-10` are treated as zero.

Error behavior summary:

- Malformed JSON -> `400`.
- Invalid matrix payloads -> `400`.
- Go returns `502` when the statistics service is unavailable.
- Unexpected internal errors return generic `500` responses without leaking internals.

Refer to `docs/openapi.yaml` or `docs/api-reference.html` for exact error response contracts.

## Project Structure

```text
matrix-qr-analytics/
├── qr-service-go/
│   ├── cmd/server/
│   ├── configs/
│   ├── internal/
│   │   ├── clients/
│   │   ├── handlers/
│   │   ├── models/
│   │   └── services/
│   └── integration/
├── statistics-service-node/
│   ├── config/
│   ├── src/
│   │   ├── controllers/
│   │   ├── errors/
│   │   ├── routes/
│   │   ├── services/
│   │   ├── app.ts
│   │   ├── index.ts
│   │   └── types.ts
│   └── test-support/
├── docs/
│   ├── openapi.yaml
│   ├── api-reference.html
│   ├── architecture.md
│   ├── assumptions.md
│   └── challenge-spec.md
├── postman/
│   └── Matrix-QR-Analytics.postman_collection.json
├── docker-compose.yml
├── .env.example
└── README.md
```

## Design Decisions and Assumptions

- QR decomposition is treated as the primary matrix operation for this challenge scope.
- Services communicate synchronously over HTTP.
- No persistence/database layer is implemented.
- Go is the facade/orchestrator; Node is the specialized statistics service.
- Rectangular matrices are supported across square, tall, and wide shapes.
- Statistics aggregate values from both Q and R.
- JWT/authentication is not implemented in current scope.
- Frontend, cloud deployment, and CI/CD are not implemented in current scope.

Deeper context:

- `docs/architecture.md`
- `docs/assumptions.md`
- `docs/challenge-spec.md`
