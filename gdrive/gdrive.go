package gdrive

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/opensaucerer/bifrost/shared/errors"
	"github.com/opensaucerer/bifrost/shared/types"
	"google.golang.org/api/drive/v3"
)

func (g *GoogleDriveStorage) UploadFile(fileFace interface{}) (*types.UploadedFile, error) {

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

	srv, err := drive.NewService(ctx)
	if err != nil {
		return nil, &errors.BifrostError{
			Err:       err,
			ErrorCode: errors.ErrFileOperationFailed,
		}
	}

	upload := &drive.File{Name: file.Name()}

	f, err := srv.Files.Create(upload).Do()
	if err != nil {
		return nil, &errors.BifrostError{
			Err:       err,
			ErrorCode: errors.ErrFileOperationFailed,
		}
	}

	return &types.UploadedFile{
		Name: f.Name,
		Path: bFile.Path,
		Size: f.Size,
		URL:  f.IconLink,
	}, nil
}

func (g *GoogleDriveStorage) UploadMultiFile(multiFace interface{}) ([]*types.UploadedFile, error) {
	return nil, nil
}

func (g *GoogleDriveStorage) IsConnected() bool {
	return false
}

func (g *GoogleDriveStorage) UploadFolder(foldFace interface{}) ([]*types.UploadedFile, error) {
	return nil, nil
}

func (g *GoogleDriveStorage) Disconnect() error {

	return nil
}

// Config returns the Google Cloud Storage configuration.
func (g *GoogleDriveStorage) Config() *types.BridgeConfig {
	return &types.BridgeConfig{
		Provider:        g.Provider,
		DefaultBucket:   g.DefaultBucket,
		CredentialsFile: g.CredentialsFile,
		Project:         g.Project,
		DefaultTimeout:  g.DefaultTimeout,
		EnableDebug:     g.EnableDebug,
	}
}
