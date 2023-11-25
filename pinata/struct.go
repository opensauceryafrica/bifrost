package pinata

import (
	"github.com/opensaucerer/bifrost/shared/request"
	"github.com/opensaucerer/bifrost/shared/types"
)

// PinataCloud is the Pinata IPFS struct
type PinataCloud struct {
	// Provider is the name of the cloud storage service to use.
	Provider types.Provider
	// DefaultTimeout is the time-to-live for time-dependent pinata operations
	DefaultTimeout int64
	// PublicRead enables public read access to uploaded files.
	PublicRead bool
	// Pinata authorization JWT
	PinataJWT string
	// UseAsync enables asynchronous operations with go routines.
	UseAsync bool
	// Pinata request client
	Client *request.Client
	// EnableDebug enables debug logging.
	EnableDebug bool
}
