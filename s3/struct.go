package s3

import "cloud.google.com/go/storage"

// SimpleStorageService is the S3 struct
type SimpleStorageService struct {
	// Provider is the name of the cloud storage service to use.
	Provider string
	// DefaultBucket is the S3 bucket to use for storage
	DefaultBucket string
	// CredentialsFile is the path to the S3 credentials file
	CredentialsFile string
	// Region is the S3 region to use for storage
	Region string
	// Zone is the S3 zone to use for storage
	Zone string
	// DefaultTimeout is the time-to-live for time-dependent S3 operations
	DefaultTimeout int64
	// Client is the S3 client
	Client *storage.Client
	// UseAsync enables asynchronous operations with go routines.
	UseAsync bool
}
