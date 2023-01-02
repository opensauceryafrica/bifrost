package bifrost

import (
	"cloud.google.com/go/storage"
)

type Bifrost struct {
}

type GoogleCloudStorage struct {
	// DefaultBucket is the Google Cloud Storage bucket to use for storage
	DefaultBucket string
	// CredentialsFile is the path to the Google Cloud Storage credentials file
	CredentialsFile string
	// Project is the Google Cloud Storage project to use for storage
	Project string
	// Region is the Google Cloud Storage region to use for storage
	Region string
	// Zone is the Google Cloud Storage zone to use for storage
	Zone string
	// DefaultTimeout is the default timeout for Google Cloud Storage operations
	DefaultTimeout string
	// Client is the Google Cloud Storage client
	Client *storage.Client
}

type SimpleStorageService struct {
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
	DefaultTimeout string
}
