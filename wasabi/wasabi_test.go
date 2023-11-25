package wasabi_test

import (
	"os"
	"testing"

	"github.com/opensaucerer/bifrost"
)

var (
	bridge bifrost.RainbowBridge
	err    error

	API_KEY     = os.Getenv("API_KEY")
	API_SECRET  = os.Getenv("API_SECRET")
	BUCKET_NAME = os.Getenv("BUCKET_NAME")
	REGION      = os.Getenv("REGION")
)

func setup(t *testing.T) {

	bridge, err = bifrost.NewRainbowBridge(&bifrost.BridgeConfig{
		DefaultBucket:  BUCKET_NAME,
		DefaultTimeout: 10,
		Provider:       bifrost.WasabiCloudStorage,
		EnableDebug:    true,
		PublicRead:     true,
		AccessKey:      API_KEY,
		SecretKey:      API_SECRET,
		Region:         REGION,
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

// go test -v -run TestWasabi ./...
func TestWasabi(t *testing.T) {
	setup(t)
	defer teardown()

	t.Run("Tests UploadFile method", func(t *testing.T) {

		f, _ := os.Open("../shared/image/aand.png")

		o, err := bridge.UploadFile(bifrost.File{
			// Path:     "../shared/image/aand.png",
			Handle:   f,
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

		f, _ := os.Open("../shared/image/hair.jpg")

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
					Path:     "../shared/image/bifrost.webp",
					Filename: "bifrost_bridge.webp",
					Options: map[string]interface{}{
						bifrost.OptMetadata: map[string]string{
							"originalname": "bifrost.jpg",
							"universe":     "Marvel",
						},
					},
				},
				{
					Path:     "",
					Handle:   f,
					Filename: "sammy.jpg",
					Options: map[string]interface{}{
						bifrost.OptMetadata: map[string]string{
							"originalname": "hair.jpg",
							"specie":       "Human",
						},
						bifrost.OptACL: bifrost.ACLPublicRead,
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
