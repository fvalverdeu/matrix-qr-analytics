# Design

## Architecture Overview

The solution consists of two independent services.

```text
Client
  |
  v

QR Service (Go + Fiber)

  |
  | HTTP POST
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

## Service Responsibilities

### QR Service

Responsibilities:

- Request validation
- QR decomposition
- Communication with Statistics Service
- Response aggregation

Suggested structure:

```text
qr-service-go/

cmd/server

internal/
  handlers/
  services/
  clients/
  models/
```

---

### Statistics Service

Responsibilities:

- Statistics calculation
- Diagonal matrix verification

Suggested structure:

```text
statistics-service-node/

src/
  routes/
  controllers/
  services/
  models/
```

---

## API Contracts

### Go Service

POST /api/v1/qr

Request:

```json
{
  "matrix": [
    [1, 2],
    [3, 4]
  ]
}
```

Response:

```json
{
  "q": [],
  "r": [],
  "statistics": {}
}
```

---

### Node Service

POST /api/v1/statistics

Request:

```json
{
  "q": [],
  "r": []
}
```

Response:

```json
{
  "max": 0,
  "min": 0,
  "average": 0,
  "sum": 0,
  "hasDiagonalMatrix": false
}
```

---

## Technical Decisions

### QR Calculation

Use Gonum.

Reason:

- Mature numerical library.
- Reliable implementation.
- Reduces implementation risk.

---

### Inter-Service Communication

Use synchronous HTTP.

Reason:

- Simple.
- Matches challenge requirements.
- Easy to test.

---

### Deployment

Use Docker Compose for local execution.

Containers:

- qr-service-go
- statistics-service-node
