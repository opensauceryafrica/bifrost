package s3_test

import (
	"os"
	"testing"

	"github.com/opensaucerer/bifrost"
)

var (
	bridge bifrost.RainbowBridge
	err    error

	AWS_ACCESS_KEY_ID     = os.Getenv("AWS_ACCESS_KEY_ID")
	AWS_SECRET_ACCESS_KEY = os.Getenv("AWS_SECRET_ACCESS_KEY")
	AWS_BUCKET_NAME       = os.Getenv("AWS_BUCKET_NAME")
)

func setup(t *testing.T) {

	bridge, err = bifrost.NewRainbowBridge(&bifrost.BridgeConfig{
		DefaultBucket:  AWS_BUCKET_NAME,
		DefaultTimeout: 10,
		Provider:       bifrost.SimpleStorageService,
		EnableDebug:    true,
		PublicRead:     true,
		AccessKey:      AWS_ACCESS_KEY_ID,
		SecretKey:      AWS_SECRET_ACCESS_KEY,
		Region:         "ap-northeast-1",
	})
	if err != nil {
		if err.(bifrost.Error).Code() == bifrost.ErrInvalidProvider {
			t.Error("Whoops, you didn't specify a valid provider!")
			return
		}
		t.Error(err.(bifrost.Error).Code(), err)
		return
	}
	defer bridge.Disconnect()

	t.Logf("Connected to %s\n", bridge.Config().Provider)
}

func TestUPload(t *testing.T) {
	setup(t)

	t.Run("Upload to S3", func(t *testing.T) {
		uploadedFile, err := bridge.UploadFile("../image/file.png", "file.png", map[string]interface{}{
			bifrost.OptACL: bifrost.ACLPublicRead,
			bifrost.OptMetadata: map[string]string{
				"originalname": "file.png",
			},
		})
		if err != nil {
			t.Errorf("Failed to upload file: %v", err)
		}
		t.Logf("Uploaded file: %s to %s\n", uploadedFile.Name, uploadedFile.Preview)
	})
}
