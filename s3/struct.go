package s3

import "cloud.google.com/go/storage"

type SimpleStorageService struct {
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
}
