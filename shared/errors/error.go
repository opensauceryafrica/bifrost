package errors

// Error returns the error message.
func (e *BifrostError) Error() string {
	// return the error message
	return e.Err.Error()
}

// Code returns the error code.
func (e *BifrostError) Code() string {
	// return the error code
	return e.ErrorCode
}
