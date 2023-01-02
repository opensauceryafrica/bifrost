package errors

func (e *BifrostError) Error() string {
	// return the error message
	return e.Err.Error()
}
