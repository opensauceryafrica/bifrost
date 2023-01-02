package errors

const (
	// ErrBadRequest is returned when something fails due to client error.
	BadRequest = 400

	// ErrUnauthorized is returned when something fails due to client not being authorized.
	Unauthorized = 401

	// ErrInvalidConfig is returned when the config is invalid.
	ErrInvalidConfig = "invalid BridgeConfig"

	// ErrInvalidBucket is returned when the bucket is invalid.
	ErrInvalidBucket = "invalid bucket"

	// ErrInvalidProvider is returned when the provider is invalid.
	ErrInvalidProvider = "invalid provider"
)
