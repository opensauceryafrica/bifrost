package bifrost

type BridgeConfig struct {
	// Provider is the name of the cloud storage service to use
	Provider string
	// Zone is the Google Cloud Storage zone to use for storage
	Zone string
	// The name of the S3 bucket to use for storing
	DefaultBucket string
	// Path to the S3 credentials file
	CredentialsFile string
	// The secret key for IAM authentication
	SecretKey string
	// The access key for IAM authentication
	AccessKey string
	// The name of the S3 region to use for storing
	Region string
	// The default timeout for S3 operations
	DefaultTimeout int64
	// EnableDebug enables debug logging
	EnableDebug bool
	// Project is the cloud project to use for storage
	// This is only implemented by some providers (e.g. Google Cloud Storage)
	Project string
}

type bridge interface {
	UploadFile(path, filename string) error
	Disconnect() error
}
