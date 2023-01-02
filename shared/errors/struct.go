package errors

// BifrostError is the error struct.
type BifrostError struct {
	// Error is the error message
	Err error
	// Code is the error code
	ErrorCode int
}
