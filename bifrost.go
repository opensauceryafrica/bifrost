/* package bifrost provides a rainbow bridge for shipping files to any cloud storage service.
it's like bifrost from marvel comics, but for files. */
package bifrost

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"

	"cloud.google.com/go/storage"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/opensaucerer/bifrost/gcs"
	"github.com/opensaucerer/biforst/gdrive"
	"github.com/opensaucerer/bifrost/pinata"
	bs3 "github.com/opensaucerer/bifrost/s3"
	bconfig "github.com/opensaucerer/bifrost/shared/config"
	"github.com/opensaucerer/bifrost/shared/errors"
	"github.com/opensaucerer/bifrost/shared/request"
	"google.golang.org/api/option"
)

// NewRainbowBridge returns a new Rainbow Bridge for shipping files to your specified cloud storage service.
func NewRainbowBridge(bc *BridgeConfig) (RainbowBridge, error) {
	// vefify that the config is valid
	if bc == nil {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("config is nil"),
			ErrorCode: errors.ErrBadRequest,
		}
	}

	// verify that the config is a struct
	t := reflect.TypeOf(bc)
	if t.Elem().Kind() != reflect.Struct {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("invalid config type: %s", t.Elem().Kind()),
			ErrorCode: errors.ErrBadRequest,
		}
	}

	// verify that the config struct is of valid type
	if t.Elem().Name() != bridgeConfigType {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("invalid config type: %s", t.Elem().Name()),
			ErrorCode: errors.ErrBadRequest,
		}
	}

	// verify that the config struct has a valid provider
	if bc.Provider == "" {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("no provider specified"),
			ErrorCode: errors.ErrBadRequest,
		}
	}

	// verify that the provider is valid
	if _, ok := providers[strings.ToLower(bc.Provider)]; !ok {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("invalid provider: %s", bc.Provider),
			ErrorCode: errors.ErrBadRequest,
		}
	}

	// verify that the config struct has a valid bucket
	if bc.DefaultBucket == "" {
		// some provider might not require a bucket
		// Just log a warning
		if bc.EnableDebug {
			// TODO: create a logger
			log.Printf(errors.WARN+"WARN: "+errors.NONE+"No bucket specified for provider %s. This might cause errors or require you to specify a bucket for each operation.", providers[strings.ToLower(bc.Provider)])
		}
	}

	// Create a new bridge based on the provider
	switch bc.Provider {
	case "s3":
		return newSimpleStorageService(bc)
	case "gdrive":
		return NewGoogleDriveStorage(bc), nil
	case "gcs":
		return newGoogleCloudStorage(bc)
	case "pinata":
		return newPinataCloud(bc)
	default:
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("invalid provider: %s", bc.Provider),
			ErrorCode: errors.ErrBadRequest,
		}
	}

}

// newPinataCloud returns a new client for Pinata Cloud.
func newPinataCloud(g *BridgeConfig) (RainbowBridge, error) {
	// TODO: add support for API key and API secret
	if g.PinataJWT == "" {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("pinata JWT is required"),
			ErrorCode: errors.ErrUnauthorized,
		}
	}

	var p = pinata.PinataCloud{
		PinataJWT:      g.PinataJWT,
		Provider:       providers[strings.ToLower(g.Provider)],
		DefaultTimeout: g.DefaultTimeout,
		PublicRead:     g.PublicRead,
		UseAsync:       g.UseAsync,
		EnableDebug:    g.EnableDebug,
		Client:         request.NewClient(bconfig.URLPinataAuth, g.PinataJWT, g.DefaultTimeout),
	}
	// authenticate with Pinata Cloud
	if err := p.Preflight(); err != nil {
		return nil, err
	}
	// return a new Pinata Cloud Storage Provider
	return &p, nil
}

// newGoogleCloudStorage returns a new client for Google Cloud Storage.
func newGoogleCloudStorage(g *BridgeConfig) (RainbowBridge, error) {
	var client *storage.Client
	var err error
	if g.CredentialsFile != "" {
		// first attempt to authenticate with credentials file
		client, err = storage.NewClient(context.Background(), option.WithCredentialsFile(g.CredentialsFile))
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
		Provider:        providers[strings.ToLower(g.Provider)],
		DefaultBucket:   g.DefaultBucket,
		CredentialsFile: g.CredentialsFile,
		Project:         g.Project,
		DefaultTimeout:  g.DefaultTimeout,
		Client:          client,
		EnableDebug:     g.EnableDebug,
		PublicRead:      g.PublicRead,
		UseAsync:        g.UseAsync,
	}, nil
}

// newSimpleStorageService returns a new client for AWS S3
func newSimpleStorageService(g *BridgeConfig) (RainbowBridge, error) {
	var client *awss3.Client
	if g.AccessKey != "" && g.SecretKey != "" {
		creds := credentials.NewStaticCredentialsProvider(g.AccessKey, g.SecretKey, "")
		cfg, err := awsconfig.LoadDefaultConfig(context.Background(), awsconfig.WithCredentialsProvider(creds), awsconfig.WithRegion(g.Region))
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
		Provider:      providers[strings.ToLower(g.Provider)],
		DefaultBucket: g.DefaultBucket,
		Region:        g.Region,
		PublicRead:    g.PublicRead,
		SecretKey:     g.SecretKey,
		AccessKey:     g.AccessKey,
		Client:        client,
	}, nil
}

func newGoogleDriveStoage(g *BridgeConfig) (rainbowBridge, error) {
	var client *http.Client
	var token oauth2.Token{}
	var err error
	if g.CredentialsFile != nil {
		f, err := os..Open(g.CredentialsFile)
		if err != nil {
			return nil, &BifrostError{
				Err: err,
				ErrorCode: errors.ErrUnauthorized,
			}
		}

		defer f.Close()
		err = json.NewDecoder(f).Decode(tok)
		if err != nil {
			return nil, &BifrostError{
				Err: err,
				ErrorCode: errors.ErrUnauthorized,
			}
		}
	}

	client = config.Client(ctx, tok)
	return &gdrive.GoogleDriveStorage{
		Client: client,
		PublicRead: g.PublicRead,
		EnableDebug: g.EnableDebug,
		DefaultTimeout: g.DefaultTimeout,
	}, nil
}
