package errors

type BifrostError struct {
	// Error is the error message
	Err error
	// Code is the error code
	Code int
}
