package pinata

// pinataIPFSStorage is the Pinata IPFS struct
type PinataCloud struct {
	// Provider is the name of the cloud storage service to use.
	Provider string
	// DefaultTimeout is the time-to-live for time-dependent Google Cloud Storage operations
	DefaultTimeout int64
		// PublicRead enables public read access to uploaded files.
	PublicRead bool
	// Pinata authorization JWT
	PinataJWT string
}
