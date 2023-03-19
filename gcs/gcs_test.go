package gcs_test

import (
	"os"
	"testing"

	"github.com/opensaucerer/bifrost"
)

var (
	bridge bifrost.RainbowBridge
	err    error

	GOOGLE_BUCKET_NAME    = os.Getenv("GOOGLE_BUCKET_NAME")
	CREDENTIALS_FILE_PATH = os.Getenv("CREDENTIALS_FILE_PATH")
)

func setup(t *testing.T) {

	bridge, err = bifrost.NewRainbowBridge(&bifrost.BridgeConfig{
		DefaultBucket:   GOOGLE_BUCKET_NAME,
		DefaultTimeout:  10,
		CredentialsFile: CREDENTIALS_FILE_PATH, // this is not required if you are using the default credentials
		Provider:        bifrost.GoogleCloudStorage,
		EnableDebug:     true,
		PublicRead:      true,
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

func TestGCSUploadFile(t *testing.T) {
	setup(t)

	t.Run("Tests Google Cloud Storage UploadFile method", func(t *testing.T) {
		o, err := bridge.UploadFile("../image/file.png", "file.png", map[string]interface{}{
			bifrost.OptACL: bifrost.ACLPublicRead,
			bifrost.OptMetadata: map[string]string{
				"originalname": "file.png",
			},
		})
		if err != nil {
			t.Errorf("Failed to upload file: %v", err)
		}
		t.Logf("Uploaded file: %s to %s\n", o.Name, o.Preview)
	})
}
