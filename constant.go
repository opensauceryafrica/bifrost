package bifrost

var (
	// providers is a map of the supported providers
	providers = map[string]string{
		"s3":  "Simple Storage Service",
		"gcs": "Google Cloud Storage",
	}
)

// Misc constants
const (
	// SimpleStorageService is the identifier of the S3 provider
	SimpleStorageService = "s3"
	// GoogleCloudStorage is the identifier of the Google Cloud Storage provider
	GoogleCloudStorage = "gcs"

	// BridgeConfigType is the type of the bridge configuration
	bridgeConfigType = "BridgeConfig"
)

// Error constants
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

	// ErrInvalidCredentials is returned when the authentication credentials are invalid.
	ErrInvalidCredentials = "invalid credentials"
)
