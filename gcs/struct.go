package gcs

import "cloud.google.com/go/storage"

type GoogleCloudStorage struct {
	// DefaultBucket is the Google Cloud Storage bucket to use for storage
	DefaultBucket string
	// CredentialsFile is the path to the Google Cloud Storage credentials file
	CredentialsFile string
	// Project is the Google Cloud Storage project to use for storage
	Project string
	// DefaultTimeout is the time-to-live for time-dependent Google Cloud Storage operations
	DefaultTimeout int64
	// Client is the Google Cloud Storage client
	Client *storage.Client
}
