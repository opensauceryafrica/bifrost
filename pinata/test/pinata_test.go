package pinata

import ("testing"
	"github.com/opensaucerer/bifrost"
	"github.com/joho/godotenv"
	"github.com/opensaucerer/bifrost/shared/types"
	"os"
	"log"
)

type RainbowBridge interface {
    UploadFile(path, filename string, options map[string]interface{}) (*types.UploadedFile, error)
    Disconnect() error
    Config() *types.BridgeConfig
}

func getRainboBridge() RainbowBridge {

	err := godotenv.Load("../../.env")

  if err != nil {
    log.Fatal("Error loading .env file")
  }

  PINATA_JWT := os.Getenv("PINATA_JWT")

	bridge, err := bifrost.NewRainbowBridge(&bifrost.BridgeConfig{
		PinataJWT: PINATA_JWT,
		Provider:  bifrost.PinataIPFSStorage,
	})

	if err != nil {
		panic(err)
	}

	return bridge
}
func TestUploadFileSucceeds(t *testing.T) {
	 
	bridge := getRainboBridge()

	 options := map[string]any{"cidVersion": 1,"wrapWithDirectory": true}
	 _, err := bridge.UploadFile("./hello.php", "hello.php", options)

	 if err != nil {
		t.Errorf("Failed to upload file: %v", err)
	}
}