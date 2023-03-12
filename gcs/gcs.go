package gcs

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/opensaucerer/bifrost/shared/config"
	"github.com/opensaucerer/bifrost/shared/errors"
	"github.com/opensaucerer/bifrost/shared/types"
)

/*
UploadFile uploads a file to Google Cloud Storage and returns an error if one occurs.

Note: UploadFile requires that a default bucket be set in bifrost.BridgeConfig.
*/
func (g *GoogleCloudStorage) UploadFile(path, filename string, options map[string]interface{}) (*types.UploadedFile, error) {
	if !g.IsConnected() {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("no active Google Cloud Storage client"),
			ErrorCode: errors.ErrClientError,
		}
	}
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
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("file does not exist: %s", path),
			ErrorCode: errors.ErrBadRequest,
		}
	}
	// open file
	file, err := os.Open(path)
	if err != nil {
		return nil, &errors.BifrostError{
			Err:       err,
			ErrorCode: errors.ErrFileOperationFailed,
		}
	}
	// close file
	defer file.Close()
	// Upload file to Google Cloud Storage
	obj := g.Client.Bucket(g.DefaultBucket).Object(filename)
	wc := obj.NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return nil, &errors.BifrostError{
			Err:       err,
			ErrorCode: errors.ErrFileOperationFailed,
		}
	}
	// close writer
	if err := wc.Close(); err != nil {
		return nil, &errors.BifrostError{
			Err:       err,
			ErrorCode: errors.ErrFileOperationFailed,
		}
	}
	// configure upload options
	if options != nil {
		// check the options map for acl settings
		if acl, ok := options[config.OptACL]; ok {
			switch acl.(string) {
			case config.ACLPublicRead:
				// set public read permissions
				if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
					return nil, &errors.BifrostError{
						Err:       err,
						ErrorCode: errors.ErrFileOperationFailed,
					}
				}
			case config.ACLPrivate:
				// set private permissions
				if err := obj.ACL().Set(ctx, storage.AllAuthenticatedUsers, storage.RoleReader); err != nil {
					return nil, &errors.BifrostError{
						Err:       err,
						ErrorCode: errors.ErrFileOperationFailed,
					}
				}
			}
		} else if g.PublicRead { // check the bridge config for default acl settings
			// set public read permissions
			if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
				return nil, &errors.BifrostError{
					Err:       err,
					ErrorCode: errors.ErrFileOperationFailed,
				}
			}
		}
	}
	// configure upload options
	if options != nil {
		// check the options map for metadata settings
		if metadata, ok := options[config.OptMetadata]; ok {
			// set metadata
			if _, err := obj.Update(ctx, storage.ObjectAttrsToUpdate{Metadata: metadata.(map[string]string)}); err != nil {
				return nil, &errors.BifrostError{
					Err:       err,
					ErrorCode: errors.ErrFileOperationFailed,
				}
			}
		}
	}
	// get object attributes
	objAttrs, _ := obj.Attrs(ctx)
	return &types.UploadedFile{
		Name:           objAttrs.Name,
		Bucket:         objAttrs.Bucket,
		Path:           path,
		Size:           objAttrs.Size,
		URL:            objAttrs.MediaLink,
		Preview:        fmt.Sprintf(config.URLGoogleCloudStorage, objAttrs.Bucket, objAttrs.Name),
		ProviderObject: obj,
	}, nil
}

/*
Disconnect closes the Google Cloud Storage connection and returns an error if one occurs.

Disconnect should only be called when the connection is no longer needed.
*/
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
		UseAsync:        g.UseAsync,
	}
}

// IsConnected returns true if the Google Cloud Storage client is connected.
func (g *GoogleCloudStorage) IsConnected() bool {
	return g.Client != nil
}

/*
	UploadFolder uploads a folder to the provider storage and returns an error if one occurs.

	Note: for some providers, UploadFolder requires that a default bucket be set in bifrost.BridgeConfig.
*/
func (g *GoogleCloudStorage) UploadFolder(path string, options map[string]interface{}) ([]*types.UploadedFile, error) {
	return nil, nil
}
