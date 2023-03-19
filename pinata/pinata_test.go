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

func TestPinataUploadFile(t *testing.T) {
	setup(t)

	t.Run("Tests Pinata UploadFile method", func(t *testing.T) {
		o, err := bridge.UploadFile("../image/file.png", "pinata_file.png", map[string]interface{}{
			bifrost.OptPinata: map[string]interface{}{
				"cidVersion": 1,
			},
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
