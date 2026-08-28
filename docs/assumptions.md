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

## Statistics Scope

Statistics are calculated using all numeric values contained in:

- Matrix Q
- Matrix R

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

---

## Security

JWT authentication is considered optional and is not included in the initial implementation.

---

## Testing

Basic unit tests may be implemented if time allows.

Testing is not considered a blocker for delivering the solution.
