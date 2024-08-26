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
type BridgeConfig types.BridgeConfig

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
	/*
	DeleteFile deletes a file from a bucket in provider's storage and returns an error if one occurs.

	Note: DeleteFile requires that a default bucket be set in bifrost.BridgeConfig.
	*/
	DeleteFile(fileFace interface{}) error
}

// BifrostError is the interface for errors returned by Bifrost.
type Error interface {
	Error() string
	Code() string
}

// PinataPinFileResponse is the response from Pinata Cloud when pinning a file.
type PinataPinFileResponse = types.PinataPinFileResponse

// MultiFile is the struct for uploading multiple files.
// Along with options, you can also set global options that will be applied to all files.
type MultiFile = types.MultiFile

// File is the struct for uploading a single file.
type File = types.File
