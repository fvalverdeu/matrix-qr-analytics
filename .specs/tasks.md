# Tasks

## Task 1 - Project Scaffolding

Goal:

Create the initial project structure.

Deliverables:

- qr-service-go folder
- statistics-service-node folder
- Docker Compose file
- Environment variable templates
- Basic README

Acceptance Criteria:

- Both projects compile successfully.
- Folder structure matches architecture.md.
- No business logic implemented yet.

---

## Task 2 - QR Service

Goal:

Implement the QR decomposition service.

Deliverables:

- Matrix validation
- QR decomposition using Gonum
- POST /api/v1/qr endpoint

Acceptance Criteria:

- Valid matrices return Q and R.
- Invalid matrices return HTTP 400.

---

## Task 3 - Statistics Service

Goal:

Implement statistics calculations.

Deliverables:

- POST /api/v1/statistics endpoint
- Max calculation
- Min calculation
- Average calculation
- Sum calculation
- Diagonal matrix verification

Acceptance Criteria:

- Statistics are correctly calculated.
- Invalid requests return HTTP 400.

---

## Task 4 - Service Integration

Goal:

Integrate both services through HTTP.

Deliverables:

- Go HTTP client
- Communication with Statistics Service
- Aggregated response

Acceptance Criteria:

- QR endpoint returns decomposition and statistics.

---

## Task 5 - Containerization

Goal:

Containerize the solution.

Deliverables:

- Go Dockerfile
- Node Dockerfile
- docker-compose.yml

Acceptance Criteria:

- Entire solution starts with docker compose up.

---

## Task 6 - Documentation

Goal:

Document the solution.

Deliverables:

- README
- Setup instructions
- Architecture overview
- API examples

Acceptance Criteria:

- A developer can run the solution using the README only.
