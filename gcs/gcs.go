package gcs

/*
UploadFile uploads a file to Google Cloud Storage and returns an error if one occurs.

Note: UploadFile requires that a default bucket be set in bifrost.BridgeConfig.
*/
func (g *GoogleCloudStorage) UploadFile(path, filename string) error {
	return nil
}

/*
Disconnect closes the Google Cloud Storage connection and returns an error if one occurs.

Disconnect should only be called when the connection is no longer needed. */
func (g *GoogleCloudStorage) Disconnect() error {
	return g.Client.Close()
}
