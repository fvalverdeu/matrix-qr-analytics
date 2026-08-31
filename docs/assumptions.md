# Assumptions

## Challenge Interpretation

The challenge description contains references to both:

- QR decomposition
- Matrix rotation

These requirements appear to be inconsistent.

The implementation assumes that QR decomposition is the primary operation because it is explicitly stated in the functional requirements.

Therefore:

- Matrix rotation is not implemented.
- QR decomposition is considered the core business requirement.

---

## Rectangular Matrix Interpretation

The term rectangular matrix is interpreted to include all three shapes:

- square (`m = n`)
- tall (`m > n`)
- wide (`m < n`)

For input `A` with shape `m x n`, the service contract is:

- `Q` shape is `m x m`
- `R` shape is `m x n`

This interpretation is implemented and validated; no `m >= n` restriction is part of current behavior.

---

## Statistics Scope

Statistics are calculated using all numeric values contained in:

- Matrix Q
- Matrix R

`hasDiagonalMatrix` means at least one matrix (`Q` OR `R`) is diagonal.

Diagonal interpretation assumptions used by the implementation:

- the matrix must be square
- off-diagonal absolute values `<= 1e-10` are treated as zero

---

## Persistence

No database is required.

All processing is performed in memory.

---

## Communication

Services communicate through synchronous HTTP requests.

No message broker is required.

---

## Error Handling

If the Statistics Service is unavailable:

- The QR Service returns an error response.
- No retry mechanism is implemented.

The Go service also treats invalid downstream success contracts as unavailability at this boundary.

---

## Security

JWT authentication is considered optional and is not included in the initial implementation.

---

## Scope Boundaries

The current implementation scope excludes:

- frontend application
- cloud deployment
- CI/CD pipeline

These items are not required for the current backend runtime architecture and may be addressed independently.

---

## Testing

Testing is implemented and part of the delivered solution.

Current automated verification covers key boundaries, including:

- QR behavior and numerical invariants across representative matrix shapes
- matrix validation rules
- Go handler orchestration and error mapping
- Statistics client transport and downstream contract validation behavior
- Node validation, statistics, controller, and app behaviors
- end-to-end Go -> Node integration flow
