package types

// UploadedFile is the struct returned by the UploadFile function.
type UploadedFile struct {
	// Name is the name of the file.
	Name string
	// Local path is the path to the file.
	Path string
	// Size is the size of the file in bytes.
	Size int64
	// Location is the location of the file in the cloud.
	Location string
}
