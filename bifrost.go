/*
	package bifrost

provides a rainbow bridge for shipping files to any cloud storage service.

It's like bifrost from marvel comics, but for files.
*/
package bifrost

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"cloud.google.com/go/storage"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/opensaucerer/bifrost/gcs"
	"github.com/opensaucerer/bifrost/pinata"
	bs3 "github.com/opensaucerer/bifrost/s3"
	bconfig "github.com/opensaucerer/bifrost/shared/config"
	"github.com/opensaucerer/bifrost/shared/errors"
	"github.com/opensaucerer/bifrost/shared/request"
	"github.com/opensaucerer/bifrost/wasabi"
	"google.golang.org/api/option"

	awsv1 "github.com/aws/aws-sdk-go/aws"
	credv1 "github.com/aws/aws-sdk-go/aws/credentials"
	sessv1 "github.com/aws/aws-sdk-go/aws/session"
	s3v1 "github.com/aws/aws-sdk-go/service/s3"
)

// NewRainbowBridge returns a new Rainbow Bridge for shipping files to your specified cloud storage service.
func NewRainbowBridge(bc *BridgeConfig) (RainbowBridge, error) {
	// vefify that the config is valid
	if bc == nil {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("config is nil"),
			ErrorCode: errors.ErrInvalidConfig,
		}
	}

	// verify that the config is a struct
	t := reflect.TypeOf(bc)
	if t.Elem().Kind() != reflect.Struct {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("invalid config type: %s", t.Elem().Kind()),
			ErrorCode: errors.ErrInvalidConfig,
		}
	}

	// verify that the config struct is of valid type
	if t.Elem().Name() != bridgeConfigType {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("invalid config type: %s", t.Elem().Name()),
			ErrorCode: errors.ErrInvalidConfig,
		}
	}

	// verify that the config struct has a valid provider
	if bc.Provider == "" {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("no provider specified"),
			ErrorCode: errors.ErrInvalidProvider,
		}
	}

	// verify that the provider is valid
	if _, ok := providers[bc.Provider]; !ok {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("invalid provider: %s", bc.Provider),
			ErrorCode: errors.ErrInvalidProvider,
		}
	}

	// verify that the config struct has a valid bucket
	if bc.DefaultBucket == "" {
		// some provider might not require a bucket
		// Just log a warning
		if bc.EnableDebug {
			// @TODO: create a logger
			log.Printf(errors.WARN+"WARN: "+errors.NONE+"No bucket specified for provider %s. This might cause errors or require you to specify a bucket for each operation.", providers[bc.Provider])
		}
	}

	// Create a new bridge based on the provider
	switch bc.Provider {
	case SimpleStorageService:
		return newSimpleStorageService(bc)
	case GoogleCloudStorage:
		return newGoogleCloudStorage(bc)
	case PinataCloud:
		return newPinataCloud(bc)
	case WasabiCloudStorage:
		return newWasabiCloudStorage(bc)
	default:
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("invalid provider: %s", bc.Provider),
			ErrorCode: errors.ErrBadRequest,
		}
	}

}

// newPinataCloud returns a new client for Pinata Cloud.
func newPinataCloud(bc *BridgeConfig) (RainbowBridge, error) {
	// @TODO: add support for API key and API secret
	if bc.PinataJWT == "" {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("pinata JWT is required"),
			ErrorCode: errors.ErrUnauthorized,
		}
	}

	var p = pinata.PinataCloud{
		PinataJWT:      bc.PinataJWT,
		Provider:       providers[bc.Provider],
		DefaultTimeout: bc.DefaultTimeout,
		PublicRead:     bc.PublicRead,
		UseAsync:       bc.UseAsync,
		EnableDebug:    bc.EnableDebug,
		Client:         request.NewClient(bconfig.URLPinataAuth, bc.PinataJWT, bc.DefaultTimeout),
	}
	// authenticate with Pinata Cloud
	if err := p.Preflight(); err != nil {
		return nil, err
	}
	// return a new Pinata Cloud Storage Provider
	return &p, nil
}

// newGoogleCloudStorage returns a new client for Google Cloud Storage.
func newGoogleCloudStorage(bc *BridgeConfig) (RainbowBridge, error) {
	var client *storage.Client
	var err error
	if bc.CredentialsFile != "" {
		// first attempt to authenticate with credentials file
		client, err = storage.NewClient(context.Background(), option.WithCredentialsFile(bc.CredentialsFile))
		if err != nil {
			return nil, &errors.BifrostError{
				Err:       err,
				ErrorCode: errors.ErrUnauthorized,
			}
		}
	} else {
		// if no credentials file is specified, attempt to authenticate without credentials file
		client, err = storage.NewClient(context.Background())
		if err != nil {
			return nil, &errors.BifrostError{
				Err:       err,
				ErrorCode: errors.ErrUnauthorized,
			}
		}
	}

	// return a new Google Cloud Storage client
	return &gcs.GoogleCloudStorage{
		Provider:        providers[bc.Provider],
		DefaultBucket:   bc.DefaultBucket,
		CredentialsFile: bc.CredentialsFile,
		Project:         bc.Project,
		DefaultTimeout:  bc.DefaultTimeout,
		Client:          client,
		EnableDebug:     bc.EnableDebug,
		PublicRead:      bc.PublicRead,
		UseAsync:        bc.UseAsync,
	}, nil
}

// newSimpleStorageService returns a new client for AWS S3
func newSimpleStorageService(bc *BridgeConfig) (RainbowBridge, error) {
	var client *awss3.Client
	if bc.AccessKey != "" && bc.SecretKey != "" {
		creds := credentials.NewStaticCredentialsProvider(bc.AccessKey, bc.SecretKey, "")
		cfg, err := awsconfig.LoadDefaultConfig(context.Background(), awsconfig.WithCredentialsProvider(creds), awsconfig.WithRegion(bc.Region))
		if err != nil {
			return nil, &errors.BifrostError{
				Err:       err,
				ErrorCode: errors.ErrUnauthorized,
			}
		}
		client = awss3.NewFromConfig(cfg)
	} else {
		// Load AWS Shared Configuration
		cfg, err := awsconfig.LoadDefaultConfig(context.TODO())
		if err != nil {
			return nil, &errors.BifrostError{
				Err:       err,
				ErrorCode: errors.ErrUnauthorized,
			}
		}
		client = awss3.NewFromConfig(cfg)
	}
	return &bs3.SimpleStorageService{
		Provider:       providers[bc.Provider],
		DefaultBucket:  bc.DefaultBucket,
		Region:         bc.Region,
		DefaultTimeout: bc.DefaultTimeout,
		PublicRead:     bc.PublicRead,
		SecretKey:      bc.SecretKey,
		AccessKey:      bc.AccessKey,
		Client:         client,
		EnableDebug:    bc.EnableDebug,
		UseAsync:       bc.UseAsync,
	}, nil
}

// newWasabiCloudStorage returns a new client for Wasabi Cloud Storage
func newWasabiCloudStorage(bc *BridgeConfig) (RainbowBridge, error) {

	// Wasabi is an S3-compatible service, so we can use the same client
	var client *s3v1.S3
	if bc.AccessKey != "" && bc.SecretKey != "" {

		s3Config := awsv1.Config{
			Endpoint:         awsv1.String(fmt.Sprintf(bconfig.URLWasabiEndpoint, bc.Region)),
			Region:           awsv1.String(bc.Region),
			S3ForcePathStyle: awsv1.Bool(true),
			Credentials:      credv1.NewStaticCredentials(bc.AccessKey, bc.SecretKey, ""),
		}

		s3session, err := sessv1.NewSessionWithOptions(sessv1.Options{
			Config: s3Config,
		})

		if err != nil {
			return nil, &errors.BifrostError{
				Err:       err,
				ErrorCode: errors.ErrUnauthorized,
			}
		}

		client = s3v1.New(s3session)

	} else {

		// load default config
		s3Config := awsv1.Config{
			Endpoint:         awsv1.String(fmt.Sprintf(bconfig.URLWasabiEndpoint, bc.Region)),
			Region:           awsv1.String(bc.Region),
			S3ForcePathStyle: awsv1.Bool(true),
		}

		s3session, err := sessv1.NewSessionWithOptions(sessv1.Options{
			Config:  s3Config,
			Profile: bconfig.WasabiCloudStorage,
		})
		if err != nil {
			return nil, &errors.BifrostError{
				Err:       err,
				ErrorCode: errors.ErrUnauthorized,
			}
		}

		client = s3v1.New(s3session)

	}
	return &wasabi.WasabiCloudStorage{
		Provider:       providers[bc.Provider],
		DefaultBucket:  bc.DefaultBucket,
		Region:         bc.Region,
		DefaultTimeout: bc.DefaultTimeout,
		PublicRead:     bc.PublicRead,
		SecretKey:      bc.SecretKey,
		AccessKey:      bc.AccessKey,
		Client:         client,
		EnableDebug:    bc.EnableDebug,
		UseAsync:       bc.UseAsync,
	}, nil
}
