# Architecture

## Overview

The solution consists of two independent services communicating through HTTP.

The architecture follows a simple service-oriented approach with clear separation of responsibilities.

### Constraints

- Go + Fiber
- Node.js + Express
- Docker
- HTTP communication between services
- No persistence layer

---

## Services

### QR Service

Technology:

- Go
- Fiber

Responsibilities:

- Receive matrix input.
- Validate matrix structure.
- Compute QR decomposition.
- Send Q and R matrices to the Statistics Service.
- Aggregate and return the final response.

---

### Statistics Service

Technology:

- Node.js
- Express

Responsibilities:

- Receive Q and R matrices.
- Calculate statistical information.
- Return calculated results.

---

## Communication Flow

```text
Client
  |
  v

QR Service (Go + Fiber)

  |
  | POST /api/v1/statistics
  |
  v

Statistics Service (Node.js + Express)

  |
  v

Statistics Result

  |
  v

QR Service

  |
  v

Client
```

---

## Request Flow

1. Client sends matrix to QR Service.
2. QR Service validates the request.
3. QR Service calculates QR decomposition.
4. QR Service sends Q and R matrices to Statistics Service.
5. Statistics Service calculates statistics.
6. Statistics Service returns statistics.
7. QR Service returns a combined response.

---

## Project Structure

### Go Service

```text
qr-service-go/

cmd/server

internal/
  handlers/
  services/
  clients/
  models/

configs/
```

### Node Service

```text
statistics-service-node/

src/
  routes/
  controllers/
  services/
  models/

config/
```

---

## Deployment

Each service runs in its own Docker container.

Docker Compose is used for local orchestration.

No cloud deployment is included in the implementation.
