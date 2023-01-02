package gcs

import "cloud.google.com/go/storage"

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
	DefaultTimeout int64
	// Client is the Google Cloud Storage client
	Client *storage.Client
}
