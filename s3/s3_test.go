package s3_test

import (
	"fmt"
	"testing"

	"github.com/opensaucerer/bifrost"
)

var (
	bridge bifrost.RainbowBridge
	err    error
)

const (
	AWS_ACCESS_KEY_ID     = ""
	AWS_SECRET_ACCESS_KEY = ""
	AWS_BUCKET_NAME       = ""
)

func setup() {
	bridge, err = bifrost.NewRainbowBridge(&bifrost.BridgeConfig{
		DefaultBucket:  AWS_BUCKET_NAME,
		DefaultTimeout: 10,
		Provider:       bifrost.SimpleStorageService,
		EnableDebug:    true,
		PublicRead:     true,
		AccessKey:      AWS_ACCESS_KEY_ID,
		SecretKey:      AWS_SECRET_ACCESS_KEY,
		Region:         "us-east-2",
	})
	if err != nil {
		if err.(bifrost.Error).Code() == bifrost.ErrInvalidProvider {
			fmt.Println("Whoops, you didn't specify a valid provider!")
			return
		}
		fmt.Println(err.(bifrost.Error).Code(), err)
		return
	}
	defer bridge.Disconnect()

	fmt.Printf("Connected to %s\n", bridge.Config().Provider)
}

func TestUPload(t *testing.T) {
	setup()

	t.Run("Upload to S3", func(t *testing.T) {
		uploadedFile, err := bridge.UploadFile("./file.png", "file.png", map[string]interface{}{
			bifrost.OptACL: bifrost.OptPublicRead,
			bifrost.OptMetadata: map[string]string{
				"originalname": "file.png",
			},
		})
		if err != nil {
			fmt.Println(err.(bifrost.Error).Code(), err)
			return
		}
		fmt.Printf("Uploaded file: %s to %s\n", uploadedFile.Name, uploadedFile.Preview)
	})
}
