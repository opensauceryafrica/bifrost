// Bifrost interface for Wasabi Cloud Storage
package wasabi

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/opensaucerer/bifrost/shared/config"
	"github.com/opensaucerer/bifrost/shared/errors"
	"github.com/opensaucerer/bifrost/shared/types"
)

/*
UploadFile uploads a file to Wasabi and returns an error if one occurs.

Note: UploadFile requires that a default bucket be set in bifrost.BridgeConfig.
*/
func (w *WasabiCloudStorage) UploadFile(fileFace interface{}) (*types.UploadedFile, error) {

	// assert that the fileFace is of type bifrost.File
	bFile, ok := fileFace.(types.File)
	if !ok {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("argument must be of type bifrost.File"),
			ErrorCode: errors.ErrBadRequest,
		}
	}

	// validate struct
	if err := bFile.Validate(); err != nil {
		return nil, &errors.BifrostError{
			Err:       err,
			ErrorCode: errors.ErrInvalidParameters,
		}
	}

	if !w.IsConnected() {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("no active Wasabi client"),
			ErrorCode: errors.ErrClientError,
		}
	}
	// @TODO: add context with timeout
	// var ctx context.Context
	// var cancel context.CancelFunc
	// ctx = context.Background()
	// if w.DefaultTimeout > 0 {
	// 	ctx, cancel = context.WithTimeout(ctx, time.Duration(w.DefaultTimeout)*time.Second)
	// 	defer cancel()
	// }

	var f io.ReadSeeker

	if bFile.Path != "" {
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

		f = file

		// ensure filename
		if bFile.Filename == "" {
			bFile.Filename = filepath.Base(bFile.Path)
		}
	} else {
		// io.ReadSeeker requires that all the data be in memory and available to be seeked over, this is not ideal for large files but since it
		// is required by the v1 of the aws sdk which is Wasabi compatible, we have to do it this way

		// read the entire file into memory
		data, err := io.ReadAll(bFile.Handle)
		if err != nil {
			return nil, &errors.BifrostError{
				Err:       err,
				ErrorCode: errors.ErrFileOperationFailed,
			}
		}

		// create a new bytes reader
		f = bytes.NewReader(data)
	}

	var params *s3.PutObjectInput = &s3.PutObjectInput{
		Bucket: aws.String(w.DefaultBucket),
		Key:    aws.String(bFile.Filename),
		Body:   f,
	}
	// check the bridge config for default acl settings
	if w.PublicRead {
		// set public read permissions
		params.ACL = aws.String(s3.ObjectCannedACLPublicRead)
	}
	// configure upload options
	for k, v := range bFile.Options {
		switch k {
		// check the options map for acl settings
		case config.OptACL:
			if v, ok := v.(string); ok {
				switch v {
				case config.ACLPublicRead:
					params.ACL = aws.String(s3.ObjectCannedACLPublicRead)
				case config.ACLPrivate:
					params.ACL = aws.String(s3.ObjectCannedACLPrivate)
				}

			}
		// set the content type
		case config.OptContentType:
			if v, ok := v.(string); ok {
				params.ContentType = aws.String(v)
			}
		// set object metadata
		case config.OptMetadata:
			if v, ok := v.(map[string]string); ok {
				params.Metadata = aws.StringMap(v)
			}
		}
	}
	// Upload the file to Wasabi
	if _, err := w.Client.PutObject(params); err != nil {
		return nil, &errors.BifrostError{
			Err:       err,
			ErrorCode: errors.ErrFileOperationFailed,
		}
	}
	// head object details
	obj, err := w.Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(w.DefaultBucket),
		Key:    aws.String(bFile.Filename),
	})
	if err != nil {
		return nil, &errors.BifrostError{
			Err:       err,
			ErrorCode: errors.ErrFileOperationFailed,
		}
	}
	return &types.UploadedFile{
		Name:           bFile.Filename,
		Bucket:         w.DefaultBucket,
		Path:           bFile.Path,
		Preview:        fmt.Sprintf(config.URLWasabiCloudStorage, w.DefaultBucket, w.Region, bFile.Filename),
		Size:           *obj.ContentLength,
		ProviderObject: obj,
		URL:            fmt.Sprintf(config.URLWasabiCloudStorage, w.DefaultBucket, w.Region, bFile.Filename),
	}, nil
}

// UploadMultiFile
func (w *WasabiCloudStorage) UploadMultiFile(multiFace interface{}) ([]*types.UploadedFile, error) {

	// assert that the multiFace is of type bifrost.File
	multiFile, ok := multiFace.(types.MultiFile)
	if !ok {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("argument must be of type bifrost.MultiFile"),
			ErrorCode: errors.ErrBadRequest,
		}
	}

	// validate struct
	if err := multiFile.Validate(); err != nil {
		return nil, &errors.BifrostError{
			Err:       err,
			ErrorCode: errors.ErrInvalidParameters,
		}
	}

	if !w.IsConnected() {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("no active Google Cloud Storage client"),
			ErrorCode: errors.ErrClientError,
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

		uploadedFile, err := w.UploadFile(file)
		if err != nil {
			if w.EnableDebug {
				// log failed file and continue
				log.Printf("Upload for file at path %s failed with err: %s\n", file.Path, err.Error())
			}
			uploadedFiles = append(uploadedFiles, &types.UploadedFile{Error: err, Name: file.Filename, Path: file.Path})
			continue
		}
		uploadedFiles = append(uploadedFiles, uploadedFile)
	}

	return uploadedFiles, nil
}

// Config returns the wasabi configuration.
func (w *WasabiCloudStorage) Config() *types.BridgeConfig {
	return &types.BridgeConfig{
		DefaultBucket:  w.DefaultBucket,
		Region:         w.Region,
		AccessKey:      w.AccessKey,
		SecretKey:      w.SecretKey,
		DefaultTimeout: w.DefaultTimeout,
		EnableDebug:    w.EnableDebug,
		Provider:       w.Provider,
		UseAsync:       w.UseAsync,
	}
}

/*
Disconnect closes the Wasabi connection and returns an error if one occurs.

Disconnect should only be called when the connection is no longer needed.
*/
func (w *WasabiCloudStorage) Disconnect() error {
	if w.IsConnected() {
		w.Client = nil
	}
	return nil
}

// IsConnected returns true if the Wasabi connection is open.
func (w *WasabiCloudStorage) IsConnected() bool {
	return w.Client != nil
}

/*
UploadFolder uploads a folder to the provider storage and returns an error if one occurs.

Note: for some providers, UploadFolder requires that a default bucket be set in bifrost.BridgeConfig.
*/
func (w *WasabiCloudStorage) UploadFolder(foldFace interface{}) ([]*types.UploadedFile, error) {
	return nil, nil
}
