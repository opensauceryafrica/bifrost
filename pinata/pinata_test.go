package pinata_test

import (
	"os"
	"testing"

	"github.com/opensaucerer/bifrost"
)

var (
	bridge     bifrost.RainbowBridge
	err        error
	PINATA_JWT = os.Getenv("PINATA_JWT")
)

func setup(t *testing.T) {
	bridge, err = bifrost.NewRainbowBridge(&bifrost.BridgeConfig{
		PinataJWT: PINATA_JWT,
		Provider:  bifrost.PinataCloud,
	})
	if err != nil {
		t.Error(err)
	}
}

func TestUploadFileSucceeds(t *testing.T) {
	setup(t)

	_, err = bridge.UploadFile("../image/file.png", "../image/file.png", map[string]interface{}{"cidVersion": 1, "wrapWithDirectory": true})

	if err != nil {
		t.Errorf("Failed to upload file: %v", err)
	}

	t.Logf("Uploaded file: %s to %s\n", "file.png", "file.png")
}
