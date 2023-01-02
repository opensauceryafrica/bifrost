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
)

// NewRainbowBridge returns a new Rainbow Bridge for shipping files to your specified cloud storage service.
func NewRainbowBridge(bc *BridgeConfig) (bridge, error) {
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
	if t.Elem().Name() != BridgeConfigType {
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
			fmt.Println(errors.ErrInvalidBucket)
		}
	}

	// Create a new bridge based on the provider
	switch bc.Provider {
	// case "aws":
	// 	return NewAmazonWebServices(bc), nil
	case "gcs":
		return NewGoogleCloudStorage(bc)
	default:
		return nil, &errors.BifrostError{
			Err:  fmt.Errorf(errors.ErrInvalidProvider),
			Code: errors.BadRequest,
		}
	}

}

// NewGoogleCloudStorage returns a new client for Google Cloud Storage.
func NewGoogleCloudStorage(g *BridgeConfig) (bridge, error) {
	// first attempt to authenticate with credentials file
	client, err := storage.NewClient(context.Background(), option.WithCredentialsFile(g.CredentialsFile))
	if err != nil {
		// if authentication error occurs, reattempt without credentials file
		client, err = storage.NewClient(context.Background())
		if err != nil {
			panic(err)
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
