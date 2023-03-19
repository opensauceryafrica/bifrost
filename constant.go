package bifrost

/*
At a point, you might wonder why we have some structs and constants duplicated in the root package and in the subpackages.
This is because we want to keep the imports as simple as possible for the end user.
No need to import subpackages, just import the root package and you're good to go.

So, if you need to use the BridgeConfig struct, you can just import the root package and use it. And if you need to assert an error, you can just import the root package and use it.

It's just a design choice, others might oppose it, that's fine. But keeping the learning curve as low as possible is a priority for me.
*/

var (
	// providers is a map of the supported providers
	providers = map[string]string{
		"pinata": "Pinata Cloud Storage",
		"s3":     "Simple Storage Service",
		"gcs":    "Google Cloud Storage",
	}
)

// Misc constants
const (

	// PinataCloud is the identifier of the Pinata Cloud storage
	PinataCloud = "pinata"
	// SimpleStorageService is the identifier of the S3 provider
	SimpleStorageService = "s3"
	// GoogleCloudStorage is the identifier of the Google Cloud Storage provider
	GoogleCloudStorage = "gcs"

	// BridgeConfigType is the type of the bridge configuration
	bridgeConfigType = "BridgeConfig"
)

// Error constants.
const (
	// ErrBadRequest is returned when something fails due to client error.
	ErrBadRequest = "bad request"

	// ErrUnauthorized is returned when something fails due to client not being authorized.
	ErrUnauthorized = "unauthorized"

	// ErrInvalidConfig is returned when the config is invalid.
	ErrInvalidConfig = "invalid config"

	// ErrInvalidBucket is returned when the bucket is invalid.
	ErrInvalidBucket = "invalid bucket"

	// ErrInvalidProvider is returned when the provider is invalid.
	ErrInvalidProvider = "invalid provider"

	// ErrInvalidCredentials is returned when the authentication credentials are invalid.
	ErrInvalidCredentials = "invalid credentials"

	// ErrFileOperationFailed is returned when a file operation fails.
	ErrFileOperationFailed = "file operation failed"

	// ErrClientError is returned when the client returns an error.
	ErrClientError = "client error"
)

// Options constants.
const (
	// ACL is the option to set the ACL of the file.
	OptACL = "acl"
	// PublicRead is the option to set the ACL of the file to public read.
	ACLPublicRead = "public-read"
	// Private is the option to set the ACL of the file to private.
	ACLPrivate = "private"
	// ContentType is the option to set the content type of the file.
	OptContentType = "content-type"
	// Metadata is the option to set the metadata of the file.
	OptMetadata = "metadata"
	// OptPinata is the option to set the pinataOptions
	OptPinata = "pinataOptions"
)
