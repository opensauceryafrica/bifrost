package s3

import "github.com/aws/aws-sdk-go-v2/service/s3"

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
	// UseAsync enables asynchronous operations with go routines.
	UseAsync bool
	// s3 client
	Client *s3.Client
	// PublicRead enables public read access to uploaded files.
	PublicRead bool
	// SecretKey is the secret key for IAM authentication.
	SecretKey string
	// AccessKey is the access key for IAM authentication.
	AccessKey string
	// EnableDebug enables debug logging.
	EnableDebug bool
}
