package config

// Pinata URL constants.
const (
	// URLPinataPinFile is the endpoint for pinning files/folders to Pinata cloud.
	URLPinataPinFile = "https://api.pinata.cloud/pinning/pinFileToIPFS"
	// URLPinataPinJSON is the endpoint for pinning JSON objects to Pinata cloud.
	URLPinataPinJSON = "https://api.pinata.cloud/pinning/pinJSONToIPFS"
	// URLPinataPinCID is the endpoint for pinning CIDs to Pinata cloud.
	URLPinataPinCID = "https://api.pinata.cloud/pinning/pinByHash"
	// URLPinataAuth is the endpoint for testing authentication against provided Pinata credentials
	URLPinataAuth = "https://api.pinata.cloud/data/testAuthentication"
)
