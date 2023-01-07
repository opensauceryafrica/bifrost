package config

// URL constants.
const (
    // PinataPinFile is the API endpoint for pinning files and folders to pinata cloud
    PinataPinFile = "https://api.pinata.cloud/pinning/pinFileToIPFS"
    // This endpoint allows the sender to add and pin any JSON object they wish to Pinata's IPFS nodes.
    PinataPinJSON = "https://api.pinata.cloud/pinning/pinJSONToIPFS"
    // This endpoint allows the sender to pin Files to CID
    PinataPinCID = "https://api.pinata.cloud/pinning/pinByHash"
		// This endpoint allows for testing authentication against provided credentials
		PinataAuthentication = "https://api.pinata.cloud/data/PreFlight"
)