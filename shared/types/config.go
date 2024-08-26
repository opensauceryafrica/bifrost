package types

type BridgeConfig struct {
	// Provider is the name of the cloud storage service to use.
	Provider Provider
	// Zone is the service zone to use for storage.
	// This is only implemented by some providers (e.g. S3).
	Zone string
	// DefaultBucket is the default storage bucket to use for storing.
	// This is only implemented by some providers (e.g. Google Cloud Storage, S3).
	DefaultBucket string
	// CredentialsFile is the path to the credentials file.
	// This is only implemented by some providers (e.g. Google Cloud Storage).
	CredentialsFile string
	// SecretKey is the secret key for IAM authentication.
	SecretKey string
	// AccessKey is the access key for IAM authentication.
	AccessKey string
	// Region is the service region to use for storing.
	// This is only implemented by some providers (e.g. S3, Google Cloud Storage).
	Region string
	// DefaultTimeout is the time-to-live for time-dependent storage operations.
	DefaultTimeout int64
	// EnableDebug enables debug logging.
	EnableDebug bool
	// Project is the cloud project to use for storage.
	// This is only implemented by some providers (e.g. Google Cloud Storage).
	Project string
	// PublicRead enables public read access to uploaded files.
	PublicRead bool
	// UseAsync enables asynchronous operations with go routines.
	UseAsync bool
	// PinataJWT is the JWT generated for your Pinata cloud account
	PinataJWT string
	// Buckets specifics the list of bucket names to interact with
	Buckets []string
	// Object specifics an object name in a bucket to interact with
	Object string
}
