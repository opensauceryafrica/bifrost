package gdrive

import (
	"google.golang.org/api/drive/v3"
)

type GoogleDriveStorage struct {
	Provider string

	DefaultBucket string

	CredentialsFile string

	Project string

	DefaultTimeout int64

	Client *drive.Service

	EnableDebug bool

	PublicRead bool
}
