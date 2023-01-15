package pinata_test

import (
	"github.com/opensaucerer/bifrost"
	"github.com/opensaucerer/bifrost/shared/types"
	"os"
	"testing"
)


func TestUploadFileSucceeds(t *testing.T) {

        PINATA_JWT := os.Getenv("PINATA_JWT")

	bridge, err := bifrost.NewRainbowBridge(&bifrost.BridgeConfig{
		PinataJWT: PINATA_JWT,
		Provider:  bifrost.PinataCloud,
	})

	if err != nil {
		panic(err)
	}

	options := map[string]interface{"cidVersion": 1, "wrapWithDirectory": true}
	_, err := bridge.UploadFile("./image/file.png", "./image/file.png", options)

	if err != nil {
		t.Errorf("Failed to upload file: %v", err)
	}
}
