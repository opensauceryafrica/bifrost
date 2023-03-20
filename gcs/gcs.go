package gcs

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
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
func (g *GoogleCloudStorage) UploadFile(fileFace interface{}) (*types.UploadedFile, error) {
	// marshal interface to bytes
	fileBytes, err := json.Marshal(fileFace)
	if err != nil {
		return nil, &errors.BifrostError{
			Err:       err,
			ErrorCode: errors.ErrBadRequest,
		}
	}

	// unmarshal bytes to struct
	var bFile types.File
	if err := json.Unmarshal(fileBytes, &bFile); err != nil {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("argument must be of type bifrost.File"),
			ErrorCode: errors.ErrBadRequest,
		}
	}

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
	if _, err := os.Stat(bFile.Path); os.IsNotExist(err) {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("file does not exist: %s", bFile.Path),
			ErrorCode: errors.ErrBadRequest,
		}
	}
	// open file
	file, err := os.Open(bFile.Path)
	if err != nil {
		return nil, &errors.BifrostError{
			Err:       err,
			ErrorCode: errors.ErrFileOperationFailed,
		}
	}
	// close file
	defer file.Close()

	// ensure filename
	if bFile.Filename == "" {
		bFile.Filename = filepath.Base(bFile.Path)
	}

	// Upload file to Google Cloud Storage
	obj := g.Client.Bucket(g.DefaultBucket).Object(bFile.Filename)
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
	for k, v := range bFile.Options {
		switch k {
		// acl
		case config.OptACL:
			if v, ok := v.(string); ok {
				switch v {
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
			}
		// content type
		case config.OptContentType:
			if v, ok := v.(string); ok {
				if _, err := obj.Update(ctx, storage.ObjectAttrsToUpdate{ContentType: v}); err != nil {
					return nil, &errors.BifrostError{
						Err:       err,
						ErrorCode: errors.ErrFileOperationFailed,
					}
				}
			}
		// metadata
		case config.OptMetadata:
			if v, ok := v.(map[string]string); ok {
				if _, err := obj.Update(ctx, storage.ObjectAttrsToUpdate{Metadata: v}); err != nil {
					return nil, &errors.BifrostError{
						Err:       err,
						ErrorCode: errors.ErrFileOperationFailed,
					}
				}
			}
		}
	}

	// get object attributes
	objAttrs, _ := obj.Attrs(ctx)
	return &types.UploadedFile{
		Name:           objAttrs.Name,
		Bucket:         objAttrs.Bucket,
		Path:           bFile.Path,
		Size:           objAttrs.Size,
		URL:            objAttrs.MediaLink,
		Preview:        fmt.Sprintf(config.URLGoogleCloudStorage, objAttrs.Bucket, objAttrs.Name),
		ProviderObject: obj,
	}, nil
}

/*
	UploadMultiFile uploads mutliple files to the provider storage and returns an error if one occurs. If any of the uploads fail, the error is appended
	to the []UploadedFile.Error and also logged when debug is enabled while the rest of the uploads continue.

	Note: for some providers, UploadMultiFile requires that a default bucket be set in bifrost.BridgeConfig.
*/
func (g *GoogleCloudStorage) UploadMultiFile(multiFace interface{}) ([]*types.UploadedFile, error) {
	// marshal interface to bytes
	multiBytes, err := json.Marshal(multiFace)
	if err != nil {
		return nil, &errors.BifrostError{
			Err:       err,
			ErrorCode: errors.ErrBadRequest,
		}
	}

	// unmarshal bytes to struct
	var multiFile types.MultiFile
	if err := json.Unmarshal(multiBytes, &multiFile); err != nil {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("argument must be of type bifrost.MultiFile"),
			ErrorCode: errors.ErrBadRequest,
		}
	}

	if !g.IsConnected() {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("no active Google Cloud Storage client"),
			ErrorCode: errors.ErrClientError,
		}
	}

	if len(multiFile.Files) == 0 {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("no files to upload"),
			ErrorCode: errors.ErrBadRequest,
		}
	}

	uploadedFiles := make([]*types.UploadedFile, 0, len(multiFile.Files))

	// TODO: add concurrency when UseAsync is true
	for _, file := range multiFile.Files {
		if multiFile.GlobalOptions != nil {
			// merge global options with file options
			for k, v := range multiFile.GlobalOptions {
				// don't override if file has option already
				if _, ok := file.Options[k]; !ok {
					file.Options[k] = v
				}
			}
		}

		uploadedFile, err := g.UploadFile(file)
		if err != nil {
			if g.EnableDebug {
				// log failed file and continue
				log.Printf("Upload for file at path %s failed with err: %s\n", file.Path, err.Error())
			}
			uploadedFiles = append(uploadedFiles, &types.UploadedFile{Error: err})
			continue
		}
		uploadedFiles = append(uploadedFiles, uploadedFile)
	}

	return uploadedFiles, nil
}

/*
Disconnect closes the Google Cloud Storage connection and returns an error if one occurs.

Disconnect should only be called when the connection is no longer needed.
*/
func (g *GoogleCloudStorage) Disconnect() error {
	if g.Client != nil {
		g.Client.Close()
		g.Client = nil
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
func (g *GoogleCloudStorage) UploadFolder(foldFace interface{}) ([]*types.UploadedFile, error) {
	return nil, nil
}
