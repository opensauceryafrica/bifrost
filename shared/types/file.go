package types

import "errors"

// UploadedFile is the struct representing a completed file/files upload.
type UploadedFile struct {
	// Name is the name of the file.
	Name string
	// Bucket is the bucket the file was uploaded to.
	Bucket string
	// Path is the local path to the file.
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
	// Error is the error returned by the provider. This is only used for async operations and multi file uploads.
	Error error
}

// Param is the struct used to pass parameters to request methods.
type Param struct {
	// Files is a list of files to upload.
	Files []ParamFile
	// Data is a list of data to upload along with the files.
	Data []ParamData
}

// ParamData is the struct for uploading data along with files in a multipart request.
type ParamData struct {
	// Key is the key to use for the data.
	Key string
	// Value is the value to use for the data.
	Value string
}

// ParamFile is the struct for uploading a single file in a multipart request.
type ParamFile struct {
	// Name is the name of the file.
	Name string
	// Path is the path to the file.
	Path string
	// Key is the key to use for the file.
	Key string
}

// MultiFile is the struct for uploading multiple files.
// Along with options, you can also set global options that will be applied to all files.
type MultiFile struct {
	// Files is a list of files to upload.
	Files []File `json:"files"`
	// GlobalOptions is a map of options to store along with all the files.
	// say 3 of 4 files need to share the same option, you can set globally for those 3 files and set the 4th file's option separately, bifrost won't override the option
	GlobalOptions map[string]interface{} `json:"global_options"`
}

// Validate validates the MultiFile struct.
func (m *MultiFile) Validate() error {
	if len(m.Files) == 0 {
		return errors.New("no files to upload")
	}
	for _, file := range m.Files {
		if err := file.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// File is the struct for uploading a single file.
type File struct {
	// Path is the path to file.
	Path string `json:"path"`
	// Filename is the name to store the file as with the provider.
	Filename string `json:"filename"`
	// Options is a map of options to store along with each file.
	Options map[string]interface{} `json:"options"`
}

// Validate validates the File struct.
func (f *File) Validate() error {
	if f.Path == "" {
		return errors.New("file.Path is required")
	}
	return nil
}
