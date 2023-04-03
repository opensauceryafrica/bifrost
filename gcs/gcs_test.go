package gcs_test

import (
	"log"
	"os"
	"testing"

	"github.com/opensaucerer/bifrost"
)

var (
	bridge bifrost.RainbowBridge
	err    error

	GOOGLE_BUCKET_NAME = os.Getenv("GOOGLE_BUCKET_NAME")
	// get full file path and join with os.Getenv("CREDENTIALS_FILE_PATH")
	CREDENTIALS_FILE_PATH = os.Getenv("CREDENTIALS_FILE_PATH")
)

func setup(t *testing.T) {

	log.Println(os.Getwd())

	bridge, err = bifrost.NewRainbowBridge(&bifrost.BridgeConfig{
		DefaultBucket:   GOOGLE_BUCKET_NAME,
		DefaultTimeout:  10,
		CredentialsFile: CREDENTIALS_FILE_PATH, // this is not required if you are using the default credentials
		Provider:        bifrost.GoogleCloudStorage,
		EnableDebug:     true,
		PublicRead:      true,
	})
	if err != nil {
		t.Error(err.(bifrost.Error).Code(), err)
		return
	}

	t.Logf("Connected to %s\n", bridge.Config().Provider)
}

func teardown() {
	bridge.Disconnect()
}

func TestGCS(t *testing.T) {
	setup(t)
	defer teardown()

	t.Run("Tests UploadFile method", func(t *testing.T) {
		o, err := bridge.UploadFile(bifrost.File{
			Path:     "../shared/image/aand.png",
			Filename: "a_and_ampersand.png",
			Options: map[string]interface{}{
				bifrost.OptMetadata: map[string]string{
					"originalname": "aand.png",
				},
			},
		})
		if err != nil {
			t.Errorf("Failed to upload file: %v", err)
			return
		}
		t.Logf("Uploaded file: %s to %s\n", o.Name, o.Preview)
	})

	t.Run("Tests UploadMultiFile method", func(t *testing.T) {
		o, err := bridge.UploadMultiFile(bifrost.MultiFile{
			Files: []bifrost.File{
				{
					Path:     "../shared/image/aand.png",
					Filename: "a_and_ampersand.png",
					Options: map[string]interface{}{
						bifrost.OptMetadata: map[string]string{
							"originalname": "aand.png",
						},
						bifrost.OptACL: bifrost.ACLPublicRead,
					},
				},
				{
					Path:     "../shared/image/hair.jpg",
					Filename: "hair_of_opensaucerer.jpg",
					Options: map[string]interface{}{
						bifrost.OptMetadata: map[string]string{
							"originalname": "hair.jpg",
						},
					},
				},
				{
					Path:     "../shared/image/bifrost.webp",
					Filename: "bifrost_bridge.webp",
					Options: map[string]interface{}{
						bifrost.OptMetadata: map[string]string{
							"originalname": "bifrost.jpg",
							"universe":     "Marvel",
						},
					},
				},
			},

			// say 3 of 4 files need to share the same option, you can set globally for those 3 files and set the 4th file's option separately, bifrost won't override the option
			GlobalOptions: map[string]interface{}{
				bifrost.OptACL: bifrost.ACLPrivate,
			},
		})
		if err != nil {
			t.Errorf("Failed to upload file: %v", err)
			return
		}

		for _, file := range o {
			t.Logf("Uploaded file: %s to %s\n", file.Name, file.Preview)
		}
	})

}
