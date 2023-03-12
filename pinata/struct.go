package pinata

import "github.com/opensaucerer/bifrost/request"

// pinataIPFSStorage is the Pinata IPFS struct
type PinataCloud struct {
	// Provider is the name of the cloud storage service to use.
	Provider string
	// DefaultTimeout is the time-to-live for time-dependent pinata operations
	DefaultTimeout int64
	// PublicRead enables public read access to uploaded files.
	PublicRead bool
	// Pinata authorization JWT
	PinataJWT string
	// UseAsync enables asynchronous operations with go routines.
	UseAsync bool
	// pinata request client
	Client *request.Client
	// EnableDebug enables debug logging.
	EnableDebug bool
}
