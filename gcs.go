// package bifrost provides a rainbow bridge for shipping files to cloud storage services.
package bifrost

// NewGoogleCloudStorage returns a new client for Google Cloud Storage.
func NewGoogleCloudStorage(g *GoogleCloudStorage) GoogleCloudStorage {
	return GoogleCloudStorage{
		DefaultBucket:   g.DefaultBucket,
		CredentialsFile: g.CredentialsFile,
		Project:         g.Project,
		Region:          g.Region,
		Zone:            g.Zone,
		DefaultTimeout:  g.DefaultTimeout,
	}
}
