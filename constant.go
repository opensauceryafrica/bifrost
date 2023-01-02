package bifrost

var (
	providers = map[string]string{
		"s3":  "Simple Storage Service",
		"gcs": "Google Cloud Storage",
	}
)

const (
	// SimpleStorageService is the identifier of the S3 provider
	SimpleStorageService = "s3"
	// GoogleCloudStorage is the identifier of the Google Cloud Storage provider
	GoogleCloudStorage = "gcs"

	// BridgeConfigType is the type of the bridge configuration
	bridgeConfigType = "BridgeConfig"
)
