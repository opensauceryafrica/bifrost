package wasabi

import (
	s3v1 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/opensaucerer/bifrost/shared/types"
)

// WasabiCloudStorage is the S3 struct
type WasabiCloudStorage struct {
	// Provider is the name of the cloud storage service to use.
	Provider types.Provider
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
	// UseAsync enables asynchronous operations with go routines.
	UseAsync bool
	// s3 client
	Client *s3v1.S3
	// PublicRead enables public read access to uploaded files.
	PublicRead bool
	// SecretKey is the secret key for IAM authentication.
	SecretKey string
	// AccessKey is the access key for IAM authentication.
	AccessKey string
	// EnableDebug enables debug logging.
	EnableDebug bool
}
