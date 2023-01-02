package gcs

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/opensaucerer/bifrost/shared/errors"
	"github.com/opensaucerer/bifrost/shared/types"
)

/*
UploadFile uploads a file to Google Cloud Storage and returns an error if one occurs.

Note: UploadFile requires that a default bucket be set in bifrost.BridgeConfig.
*/
func (g *GoogleCloudStorage) UploadFile(path, filename string) error {
	// create context and add timeout if default timeout is set
	var ctx context.Context
	var cancel context.CancelFunc
	ctx = context.Background()
	if g.DefaultTimeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(g.DefaultTimeout)*time.Second)
		defer cancel()
	}

	// verify that file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &errors.BifrostError{
			Err:       fmt.Errorf("file does not exist: %s", path),
			ErrorCode: errors.ErrBadRequest,
		}
	}

	// open file
	file, err := os.Open(path)
	if err != nil {
		return &errors.BifrostError{
			Err:       err,
			ErrorCode: errors.ErrBadRequest,
		}
	}
	defer file.Close()

	// Upload file to Google Cloud Storage
	wc := g.Client.Bucket(g.DefaultBucket).Object(filename).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	return nil
}

/*
Disconnect closes the Google Cloud Storage connection and returns an error if one occurs.

Disconnect should only be called when the connection is no longer needed. */
func (g *GoogleCloudStorage) Disconnect() error {
	if g.Client != nil {
		return g.Client.Close()
	}
	return nil
}

// Config returns the Google Cloud Storage configuration.
func (g *GoogleCloudStorage) Config() *types.BridgeConfig {
	return &types.BridgeConfig{
		Provider:        g.Provider,
		DefaultBucket:   g.DefaultBucket,
		CredentialsFile: g.CredentialsFile,
		Project:         g.Project,
		DefaultTimeout:  g.DefaultTimeout,
		EnableDebug:     g.EnableDebug,
	}
}
