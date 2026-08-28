package clients

type UnavailableError struct {
	Message string
}

func (e *UnavailableError) Error() string {
	return e.Message
}

func newUnavailableError(message string) *UnavailableError {
	return &UnavailableError{Message: message}
}
