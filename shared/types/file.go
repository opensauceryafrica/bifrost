package types

import "cloud.google.com/go/storage"

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
	// Object is the object in the cloud.
	Object *storage.ObjectHandle
}
