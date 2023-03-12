package types

// UploadedFile is the struct returned by the UploadFile function.
type UploadedFile struct {
	// Name is the name of the file.
	Name string
	// Bucket is the bucket the file was uploaded to.
	Bucket string
	// Local path is the path to the file.
	Path string
	// Size is the size of the file in bytes.
	Size int64
	// URL is the location of the file in the cloud.
	URL string
	// Preview is the URL to a preview of the file.
	Preview string
	// ProviderObject is the object returned by the cloud storage provider.
	// You need to type assert this to the correct type to use it.
	ProviderObject interface{}
	// Done sends a message to signal when an async process is complete.
	Done chan bool
	// Quit receives a message to signal for the exit of an async process.
	Quit chan bool
	// CID is the content identifier for the file.
	// This is only implemented by some providers (e.g. Pinata Cloud).
	CID string
}

// Param is the struct used to pass parameters to request methods.
type Param struct {
	// Files is a list of files to upload.
	Files []ParamFile
	Data  []ParamData
}

type ParamData struct {
	// Key is the key to use for the data.
	Key string
	// Value is the value to use for the data.
	Value string
}

type ParamFile struct {
	// Name is the name of the file.
	Name string
	// Path is the path to the file.
	Path string
	// Key is the key to use for the file.
	Key string
}
