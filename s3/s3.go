package s3

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/opensaucerer/bifrost/shared/errors"
	"github.com/opensaucerer/bifrost/shared/types"
)

func (s *SimpleStorageService) UploadFile(path, filename string, options map[string]interface{}) (*types.UploadedFile, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx = context.Background()
	if s.DefaultTimeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(s.DefaultTimeout)*time.Second)
		defer cancel()
	}
	// verify that file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &types.UploadedFile{}, &errors.BifrostError{
			Err:       fmt.Errorf("file does not exist: %s", path),
			ErrorCode: errors.ErrBadRequest,
		}
	}
	// open file
	file, err := os.Open(path)
	if err != nil {
		return &types.UploadedFile{}, &errors.BifrostError{
			Err:       err,
			ErrorCode: errors.ErrFileOperationFailed,
		}
	}
	// close file
	defer file.Close()

	// Upload the file to S3
	if s.PublicRead {
		_, err = s.Client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(s.DefaultBucket),
			Key:    aws.String(filename),
			Body:   file,
			ACL:    awsTypes.ObjectCannedACLPublicRead,
		})
	} else {
		_, err = s.Client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(s.DefaultBucket),
			Key:    aws.String(filename),
			Body:   file,
		})
	}
	if err != nil {
		return &types.UploadedFile{}, &errors.BifrostError{
			Err:       err,
			ErrorCode: errors.ErrFileOperationFailed,
		}
	}
	return &types.UploadedFile{
		Name:   filename,
		Bucket: s.DefaultBucket,
		Path:   path,
	}, nil
}

// Config returns the s3 configuration.
func (s *SimpleStorageService) Config() *types.BridgeConfig {
	return &types.BridgeConfig{
		DefaultBucket:  s.DefaultBucket,
		Region:         s.Region,
		AccessKey:      s.AccessKey,
		SecretKey:      s.SecretKey,
		DefaultTimeout: s.DefaultTimeout,
	}
}

func (s *SimpleStorageService) Disconnect() error {
	if s.Client != nil {
		// return s.Client.
	}
	return nil
}
