package config

// URL constants.
const (
	// URLPinataPinFile is the endpoint for pinning files/folders to Pinata cloud.
	URLPinataPinFile = "https://api.pinata.cloud/pinning/pinFileToIPFS"

	// URLPinataPinJSON is the endpoint for pinning JSON objects to Pinata cloud.
	URLPinataPinJSON = "https://api.pinata.cloud/pinning/pinJSONToIPFS"

	// URLPinataPinCID is the endpoint for pinning CIDs to Pinata cloud.
	URLPinataPinCID = "https://api.pinata.cloud/pinning/pinByHash"

	// URLPinataAuth is the endpoint for testing authentication against provided Pinata credentials
	URLPinataAuth = "https://api.pinata.cloud/data/testAuthentication"

	// URLPinataGateway is the public gateway for Pinata cloud.
	URLPinataGateway = "https://gateway.pinata.cloud/ipfs/%v"

	// URLGoogleCloudStorage is the public gateway for Google Cloud Storage.
	URLGoogleCloudStorage = "https://storage.googleapis.com/%s/%s"

	// URLSimpleStorageService is the public gateway for Simple Storage Service.
	URLSimpleStorageService = "https://%s.s3.%s.amazonaws.com/%s"

	//URLWasabiCloudStorage is the public gateway for Wasabi Cloud Storage.
	URLWasabiCloudStorage = "https://%s.s3.%s.wasabisys.com/%s"

	// URLWasabiEndpoint is the endpoint for Wasabi Cloud Storage.
	URLWasabiEndpoint = "https://s3.%s.wasabisys.com"
)
