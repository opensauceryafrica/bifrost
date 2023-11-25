# How to Use Bifrost with Google Cloud Storage (GCS)
Welcome to the Bifrost documentation for Google Cloud Storage (GCS)! This guide will walk you through the steps of using Bifrost to upload files to GCS.

## Overview
Google Cloud Storage is a popular cloud storage service that allows you to store and access your data on Google's infrastructure. Bifrost provides a simple and intuitive way to upload files to GCS without having to write complex code.

## Prerequisites
Before you can start using Bifrost to upload files to Google Cloud Storage, you'll need to make sure you have the following:
- A Google Cloud account with GCS access
- A GCS bucket to upload files to
- Bifrost installed on your local machine

### Steps:
#### Create a GCS bucket
1. Login to your Google Cloud account and navigate to the GCS console
2. Click on the "Create bucket" button and follow the prompts to create a new bucket
_Note the name of your bucket as you will need it later_

#### Set up Google Cloud credentials
1. Create a new service account by navigating to the IAM & Admin console and selecting "Service accounts" from the left-hand menu
2. Click on the "Create Service Account" button and follow the prompts to create a new service account
3. Once you have created the service account, navigate to the "Keys" tab and click on the "Create Key" button
4. Select "JSON" as the key type and download the JSON key file
_Note the path to your JSON key file as you will need it later_

## Mounting a Bifrost Bridge to GCS
1. Install Bifrost using: ```go get github.com/bifrost-cloud/bifrost```
2. Initialize a new Bifrost client and mount a GCS bridge using the following code:
```go
package main

import (
	"fmt"

	"github.com/opensaucerer/bifrost"
)

func main() {
	bridge, err := bifrost.NewRainbowBridge(&bifrost.BridgeConfig{
		DefaultBucket:   "bifrost",
		DefaultTimeout:  10,
		Provider:        bifrost.GoogleCloudStorage,
		CredentialsFile: "/path/to/service/account/json", // this is not required if you are using google's default credentials
		EnableDebug:     true,
		PublicRead:      true,
	})
	if err != nil {
		// bifrost comes with some error codes
		if err.(bifrost.Error).Code() == bifrost.ErrInvalidProvider {
			fmt.Println("Whoops, you didn't specify a valid provider!")
			return
		}
		fmt.Println(err.(bifrost.Error).Code(), err)
		return
	}
	defer bridge.Disconnect()
	fmt.Printf("Connected to %s\n", bridge.Config().Provider)
}
```
And that's it! You have now mounted a Bifrost bridge to your GCS account and can start uploading files via this bridge.

## Shipping a file to Google Cloud Storage via the rainbow bridge
Now that you have mounted a Bifrost bridge to GCS, you can use Bifrost to upload files to GCS using the following code:
```go
// Upload a file
uploadedFile, err := bridge.UploadFile(bifrost.File{
	Path:     "../shared/image/aand.png",
	Filename: "a_and_ampersand.png",
	Options: map[string]interface{}{
		bifrost.OptMetadata: map[string]string{
			"originalname": "aand.png",
		},
	},
})
if err != nil {
	fmt.Println(err)
	return
}
fmt.Printf("Uploaded file: %s to %s\n", uploadedFile.Name, uploadedFile.Preview)
```
As you can see, uploading a file to GCS using Bifrost is as simple as calling the UploadFile method on the Bifrost client with the path to the file on your local machine and the name to give the file on GCS.

## Shipping multiple files to Google Cloud Storage via the rainbow bridge
If you want to upload multiple files using Bifrost with GCS, you can use the UploadMultiFile method provided by the GCS bridge in Bifrost. Here is an example code snippet:

```go
// Upload multiple files
uploadedFiles, err := bridge.UploadMultiFile(bifrost.MultiFile{
	Files: []bifrost.File{
		{
			Path:     "../shared/image/aand.png",
			Filename: "a_and_ampersand.png",
			Options: map[string]interface{}{
				bifrost.OptMetadata: map[string]string{
					"originalname": "aand.png",
				},
				bifrost.OptACL: bifrost.ACLPublicRead,
			},
		},
		{
			Path:     "../shared/image/hair.jpg",
			Filename: "hair_of_opensaucerer.jpg",
			Options: map[string]interface{}{
				bifrost.OptMetadata: map[string]string{
					"originalname": "hair.jpg",
				},
			},
		},
		{
			Path:     "../shared/image/bifrost.webp",
			Filename: "bifrost_bridge.webp",
			Options: map[string]interface{}{
				bifrost.OptMetadata: map[string]string{
					"originalname": "bifrost.jpg",
					"universe":     "Marvel",
				},
			},
		},
	},

	// since we want both files to be public, we can set the global options rather than setting it for each file
	// say 3 of 4 files need to share the same option, you can set globally for those 3 files and set the 4th file's option separately, bifrost won't override the option
	GlobalOptions: map[string]interface{}{
		bifrost.OptACL: bifrost.ACLPrivate,
	},
})
if err != nil {
	fmt.Println(err)
	return
}

for _, file := range uploadedFiles {
	fmt.Printf("Uploaded file: %s to %s\n", file.Name, file.Preview)
}
```

## Additional Resources
- [Google Cloud Storage Documentation](https://cloud.google.com/storage/docs)
- [Setting up GCS](https://cloud.google.com/storage/docs/creating-buckets)
- [Setting up gcloud CLI](https://cloud.google.com/sdk/docs/install)

We hope you found this guide helpful in using Bifrost with GCS. If you have any questions or feedback, please don't hesitate to open an [issue](https://github.com/opensaucerer/bifrost/issues)!
