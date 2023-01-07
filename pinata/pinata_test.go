package pinata_test

import (
	"github.com/opensaucerer/bifrost"
	"github.com/opensaucerer/bifrost/shared/types"
	"os"
	"testing"
)

type RainbowBridge interface {
	UploadFile(path, filename string, options map[string]interface{}) (*types.UploadedFile, error)
	Disconnect() error
	Config() *types.BridgeConfig
}

func getRainboBridge() RainbowBridge {

	PINATA_JWT := os.Getenv("PINATA_JWT")

	bridge, err := bifrost.NewRainbowBridge(&bifrost.BridgeConfig{
		PinataJWT: PINATA_JWT,
		Provider:  bifrost.PinataCloud,
	})

	if err != nil {
		panic(err)
	}

	return bridge
}
func TestUploadFileSucceeds(t *testing.T) {

	bridge := getRainboBridge()

	options := map[string]any{"cidVersion": 1, "wrapWithDirectory": true}
	_, err := bridge.UploadFile("./image/file.png","./image/file.png", options)

	if err != nil {
		t.Errorf("Failed to upload file: %v", err)
	}
}
