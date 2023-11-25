# How to use Bifrost with Wasabi Cloud Storage

Welcome to the Bifrost documentation for Wasabi Cloud Storage! In this guide, we will show you how to use Bifrost to upload files to Wasabi Cloud Storage, one of the most popular cloud storage services out there.

## Overview

Wasabi Cloud Storage is a highly scalable, secure, and durable object storage service that is used by millions of customers worldwide. Bifrost provides a straightforward way to upload files to Wasabi Cloud Storage without the need to write complex code.

## Prerequisites

Before you can start using Bifrost to upload files to Wasabi Cloud Storage, you'll need to make sure you have the following:

- A Wasabi account
- Wasabi Access Key ID and Secret Access Key
- A Wasabi bucket to upload files to
- Bifrost installed on your local machine

### Steps:

#### Create an Wasabi bucket

1. Login to your Wasabi account and navigate to the Buckets page
2. Click on the "Create bucket" button and follow the prompts to create a new bucket
   _Note the name of your bucket as you will need it later_

#### Set up Wasabi credentials

1. Create a new access key for your Wasabi account by navigating to the Access Keys page from the left-hand menu
2. Click on the "Create access key" button, follow the prompt, and note down the Access Key ID and Secret Access Key

## Mount a Bifrost bridge to your Wasabi account

1. Install Bifrost using: `go get github.com/bifrost-cloud/bifrost`
2. Create a new Bifrost client and mount an Wasabi bridge using the following code:

```go
package main

import (
	"fmt"
	"github.com/opensaucerer/bifrost"
)

func main() {
	bridge, _ := bifrost.NewRainbowBridge(&bifrost.BridgeConfig{
		DefaultBucket: "default-bucket",
		Provider:      bifrost.WasabiCloudStorage,
		AccessKey:     "access-key",
		SecretKey:     "secret-key",
		EnableDebug:   true,
		PublicRead:    true,
		Region: "ap-northeast-1",
	})
	defer bridge.Disconnect()
	fmt.Printf("Connected to %s\n", bridge.Config().Provider)
```

## Upload files to Wasabi using Bifrost

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
}

```

And that's it! You have now mounted a Bifrost bridge to your Wasabi Wasabi account and can start uploading files via this bridge.

## Shipping multiple files to Wasabi Cloud Storage via the rainbow bridge

```go
// Upload multiple files

f, _ := os.Open("../shared/image/hair.jpg")

uploadedFiles, err := bridge.UploadMultiFile(bifrost.MultiFile{
	Files: []bifrost.File{
		{
			Path:     "../shared/image/aand.png",
			Filename: "a_and_ampersand.png",
			Options: map[string]interface{}{
				bifrost.OptMetadata: map[string]string{
					"originalname": "aand.png",
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
        {
            Path:     "",
            Handle:   f,
            Filename: "sammy.jpg",
            Options: map[string]interface{}{
                bifrost.OptMetadata: map[string]string{
                    "originalname": "hair.jpg",
                    "specie":       "Human",
                },
                bifrost.OptACL: bifrost.ACLPublicRead,
            },
        },
	},
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

- [Wasabi Cloud Storage Documentation](https://docs.wasabi.com/)

We hope you found this guide helpful in using Bifrost with Wasabi Cloud Storage. If you have any questions or feedback, please don't hesitate to open an [issue](https://github.com/opensaucerer/bifrost/issues)!
