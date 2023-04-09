// Bifrost interface for Simple Storage Service (S3)
package s3

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/opensaucerer/bifrost/shared/config"
	"github.com/opensaucerer/bifrost/shared/errors"
	"github.com/opensaucerer/bifrost/shared/types"
)

/*
UploadFile uploads a file to S3 and returns an error if one occurs.

Note: UploadFile requires that a default bucket be set in bifrost.BridgeConfig.
*/
func (s *SimpleStorageService) UploadFile(fileFace interface{}) (*types.UploadedFile, error) {
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

	// validate struct
	if err := bFile.Validate(); err != nil {
		return nil, &errors.BifrostError{
			Err:       err,
			ErrorCode: errors.ErrInvalidParameters,
		}
	}

	if !s.IsConnected() {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("no active S3 client"),
			ErrorCode: errors.ErrClientError,
		}
	}
	var ctx context.Context
	var cancel context.CancelFunc
	ctx = context.Background()
	if s.DefaultTimeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(s.DefaultTimeout)*time.Second)
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

	var params *s3.PutObjectInput = &s3.PutObjectInput{
		Bucket: aws.String(s.DefaultBucket),
		Key:    aws.String(bFile.Filename),
		Body:   file,
	}
	// check the bridge config for default acl settings
	if s.PublicRead {
		// set public read permissions
		params.ACL = awsTypes.ObjectCannedACLPublicRead
	}
	// configure upload options
	for k, v := range bFile.Options {
		switch k {
		// check the options map for acl settings
		case config.OptACL:
			if v, ok := v.(string); ok {
				switch v {
				case config.ACLPublicRead:
					params.ACL = awsTypes.ObjectCannedACLPublicRead
				case config.ACLPrivate:
					params.ACL = awsTypes.ObjectCannedACLPrivate
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
				params.Metadata = v
			}
		}
	}
	// Upload the file to S3
	if _, err := s.Client.PutObject(ctx, params); err != nil {
		return nil, &errors.BifrostError{
			Err:       err,
			ErrorCode: errors.ErrFileOperationFailed,
		}
	}
	// head object details
	obj, err := s.Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.DefaultBucket),
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
		Bucket:         s.DefaultBucket,
		Path:           bFile.Path,
		Preview:        fmt.Sprintf(config.URLSimpleStorageService, s.DefaultBucket, s.Region, bFile.Filename),
		Size:           obj.ContentLength,
		ProviderObject: obj,
		URL:            fmt.Sprintf(config.URLSimpleStorageService, s.DefaultBucket, s.Region, bFile.Filename),
	}, nil
}

// UploadMultiFile
func (s *SimpleStorageService) UploadMultiFile(multiFace interface{}) ([]*types.UploadedFile, error) {
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

	// validate struct
	if err := multiFile.Validate(); err != nil {
		return nil, &errors.BifrostError{
			Err:       err,
			ErrorCode: errors.ErrInvalidParameters,
		}
	}

	if !s.IsConnected() {
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

		uploadedFile, err := s.UploadFile(file)
		if err != nil {
			if s.EnableDebug {
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

// Config returns the s3 configuration.
func (s *SimpleStorageService) Config() *types.BridgeConfig {
	return &types.BridgeConfig{
		DefaultBucket:  s.DefaultBucket,
		Region:         s.Region,
		AccessKey:      s.AccessKey,
		SecretKey:      s.SecretKey,
		DefaultTimeout: s.DefaultTimeout,
		EnableDebug:    s.EnableDebug,
		Provider:       s.Provider,
		UseAsync:       s.UseAsync,
	}
}

/*
Disconnect closes the S3 connection and returns an error if one occurs.

Disconnect should only be called when the connection is no longer needed.
*/
func (s *SimpleStorageService) Disconnect() error {
	if s.IsConnected() {
		s.Client = nil
	}
	return nil
}

// IsConnected returns true if the S3 connection is open.
func (s *SimpleStorageService) IsConnected() bool {
	return s.Client != nil
}

/*
UploadFolder uploads a folder to the provider storage and returns an error if one occurs.

Note: for some providers, UploadFolder requires that a default bucket be set in bifrost.BridgeConfig.
*/
func (s *SimpleStorageService) UploadFolder(foldFace interface{}) ([]*types.UploadedFile, error) {
	return nil, nil
}
