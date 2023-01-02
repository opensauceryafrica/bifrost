// package bifrost provides a rainbow bridge for shipping files to cloud storage services.
package bifrost

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/opensaucerer/bifrost/errors"
	"github.com/opensaucerer/bifrost/gcs"
	"google.golang.org/api/option"
	"google.golang.org/appengine/log"
)

// NewRainbowBridge returns a new Rainbow Bridge for shipping files to your specified cloud storage service.
func NewRainbowBridge(bc *BridgeConfig) (rainbowBridge, error) {
	// vefify that the config is valid
	if bc == nil {
		return nil, &errors.BifrostError{
			Err:  fmt.Errorf(errors.ErrInvalidConfig),
			Code: errors.BadRequest,
		}
	}

	// verify that the config is a struct
	t := reflect.TypeOf(bc)
	if t.Elem().Kind() != reflect.Struct {
		return nil, &errors.BifrostError{
			Err:  fmt.Errorf(errors.ErrInvalidConfig),
			Code: errors.BadRequest,
		}
	}

	// verify that the config struct is of valid type
	if t.Elem().Name() != bridgeConfigType {
		return nil, &errors.BifrostError{
			Err:  fmt.Errorf(errors.ErrInvalidConfig),
			Code: errors.BadRequest,
		}
	}

	// verify that the config struct has a valid provider
	if bc.Provider == "" {
		return nil, &errors.BifrostError{
			Err:  fmt.Errorf(errors.ErrInvalidProvider),
			Code: errors.BadRequest,
		}
	}

	// verify that the provider is valid
	if _, ok := providers[strings.ToLower(bc.Provider)]; !ok {
		return nil, &errors.BifrostError{
			Err:  fmt.Errorf(errors.ErrInvalidProvider),
			Code: errors.BadRequest,
		}
	}

	// verify that the config struct has a valid bucket
	if bc.DefaultBucket == "" {
		// some provider might not require a bucket
		// Just log a warning
		if bc.EnableDebug {
			// TODO: create a logger
			log.Warningf(context.Background(), "No bucket specified for provider %s. This might cause errors or require you to specify a bucket for each operation.", bc.Provider)
		}
	}

	// Create a new bridge based on the provider
	switch bc.Provider {
	// case "aws":
	// 	return NewAmazonWebServices(bc), nil
	case "gcs":
		return newGoogleCloudStorage(bc)
	default:
		return nil, &errors.BifrostError{
			Err:  fmt.Errorf(errors.ErrInvalidProvider),
			Code: errors.BadRequest,
		}
	}

}

// newGoogleCloudStorage returns a new client for Google Cloud Storage.
func newGoogleCloudStorage(g *BridgeConfig) (rainbowBridge, error) {
	// first attempt to authenticate with credentials file
	client, err := storage.NewClient(context.Background(), option.WithCredentialsFile(g.CredentialsFile))
	if err != nil {
		// if authentication error occurs, reattempt without credentials file
		client, err = storage.NewClient(context.Background())
		if err != nil {
			return nil, &errors.BifrostError{
				Err:  fmt.Errorf(errors.ErrInvalidCredentials),
				Code: errors.Unauthorized,
			}
		}
	}

	return &gcs.GoogleCloudStorage{
		DefaultBucket:   g.DefaultBucket,
		CredentialsFile: g.CredentialsFile,
		Project:         g.Project,
		Region:          g.Region,
		Zone:            g.Zone,
		DefaultTimeout:  g.DefaultTimeout,
		Client:          client,
	}, nil
}
