package freeproxy

// Error represents an error response from the GetFreeProxy API
type Error struct {
	Message string
}

// Error implements the error interface for Error
func (e *Error) Error() string {
	return e.Message
}
