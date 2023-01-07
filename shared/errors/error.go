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

const (
   // Status200 refers to 200 status code
   Status200 = 200
   // Status400 refers to 400 status code
   Status400 = 400
)