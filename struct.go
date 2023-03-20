package bifrost

import "github.com/opensaucerer/bifrost/shared/types"

/*
At a point, you might wonder why we have some structs and constants duplicated in the root package and in the subpackages.
This is because we want to keep the imports as simple as possible for the end user.
No need to import subpackages, just import the root package and you're good to go.

So, if you need to use the BridgeConfig struct, you can just import the root package and use it. And if you need to assert an error, you can just import the root package and use it.

It's just a design choice, others might oppose it, that's fine. But keeping the learning curve as low as possible is a priority for me.
*/

// BridgeConfig is the configuration for the rainbow bridge.
type BridgeConfig struct {
	// Provider is the name of the cloud storage service to use.
	Provider string
	// Zone is the service zone to use for storage.
	// This is only implemented by some providers (e.g. S3).
	Zone string
	// DefaultBucket is the default storage bucket to use for storing.
	// This is only implemented by some providers (e.g. Google Cloud Storage, S3).
	DefaultBucket string
	// CredentialsFile is the path to the credentials file.
	// This is only implemented by some providers (e.g. Google Cloud Storage).
	CredentialsFile string
	// SecretKey is the secret key for IAM authentication.
	SecretKey string
	// AccessKey is the access key for IAM authentication.
	AccessKey string
	// Region is the service region to use for storing.
	// This is only implemented by some providers (e.g. S3, Google Cloud Storage).
	Region string
	// DefaultTimeout is the time-to-live for time-dependent storage operations.
	DefaultTimeout int64
	// EnableDebug enables debug logging.
	EnableDebug bool
	// Project is the cloud project to use for storage.
	// This is only implemented by some providers (e.g. Google Cloud Storage).
	Project string
	// PublicRead enables public read access to uploaded files.
	PublicRead bool
	// UseAsync enables asynchronous operations with go routines.
	UseAsync bool
	// PinataJWT is the JWT generated for your Pinata cloud account
	PinataJWT string
}

type RainbowBridge interface {
	/*
		UploadFile uploads a file to the provider storage and returns an error if one occurs.

		Note: for some providers, UploadFile requires that a default bucket be set in bifrost.BridgeConfig.
	*/
	UploadFile(fileFace interface{}) (*types.UploadedFile, error)
	/*
		UploadMultiFile uploads mutliple files to the provider storage and returns an error if one occurs. If any of the uploads fail, the error is appended
		to the []UploadedFile.Error and also logged when debug is enabled while the rest of the uploads continue.

		Note: for some providers, UploadMultiFile requires that a default bucket be set in bifrost.BridgeConfig.
	*/
	UploadMultiFile(multiFace interface{}) ([]*types.UploadedFile, error)
	/*
		Disconnect closes the provider client connection and returns an error if one occurs.

		Disconnect should only be called when the connection is no longer needed.
	*/
	Disconnect() error
	// Config returns the provider configuration.
	Config() *types.BridgeConfig
	// IsConnected returns true if there is an active connection to the provider.
	IsConnected() bool
	/*
		UploadFolder uploads a folder to the provider storage and returns an error if one occurs.

		Note: for some providers, UploadFolder requires that a default bucket be set in bifrost.BridgeConfig.
	*/
	UploadFolder(foldFace interface{}) ([]*types.UploadedFile, error)
}

// BifrostError is the interface for errors returned by Bifrost.
type Error interface {
	Error() string
	Code() string
}

// PinataPinFileResponse is the response from Pinata Cloud when pinning a file.
type PinataPinFileResponse struct {
	IpfsHash  string
	Timestamp string
	PinSize   int64
	Error     string
}

// MultiFile is the struct for uploading multiple files.
// Along with options, you can also set global options that will be applied to all files.
type MultiFile struct {
	Files []File `json:"files"`
	// GlobalOptions is a map of options to store along with all the files.
	// say 3 of 4 files need to share the same option, you can set globally for those 3 files and set the 4th file's option separately, bifrost won't override the option
	GlobalOptions map[string]interface{} `json:"global_options"`
}

type File struct {
	// Path is the path to file.
	Path string `json:"path"`
	// Filename is the name to store the file as with the provider.
	Filename string `json:"filename"`
	// Options is a map of options to store along with each file.
	Options map[string]interface{} `json:"options"`
}
