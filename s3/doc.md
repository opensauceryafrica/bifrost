# How to use Bifrost with Amazon S3
Welcome to the Bifrost documentation for Amazon S3! In this guide, we will show you how to use Bifrost to upload files to Amazon S3, one of the most popular cloud storage services out there.

## Overview
Amazon S3 is a highly scalable, secure, and durable object storage service that is used by millions of customers worldwide. Bifrost provides a straightforward way to upload files to Amazon S3 without the need to write complex code.

## Prerequisites
Before you can start using Bifrost to upload files to Amazon S3, you'll need to make sure you have the following:
- An AWS account with S3 access
- AWS Access Key ID and Secret Access Key
- An S3 bucket to upload files to
- Bifrost installed on your local machine

### Steps:
#### Create an S3 bucket
1. Login to your AWS account and navigate to the S3 console
2. Click on the "Create bucket" button and follow the prompts to create a new bucket
_Note the name of your bucket as you will need it later_

#### Set up AWS credentials
1. Create a new access key for your AWS account by navigating to the IAM console and selecting "Users" from the left-hand menu
2. Select your user and click on the "Security credentials" tab
3. Click on the "Create access key" button and note down the Access Key ID and Secret Access Key

## Mount a Bifrost bridge to your S3 account
1. Install Bifrost using: ```go get github.com/bifrost-cloud/bifrost```
2. Create a new Bifrost client and mount an S3 bridge using the following code:
```go
package main

import (
	"fmt"
	"github.com/opensaucerer/bifrost"
)

func main() {
	bridge, _ := bifrost.NewRainbowBridge(&bifrost.BridgeConfig{
		DefaultBucket: "default-bucket",
		Provider:      bifrost.SimpleStorageService,
		AccessKey:     "access-key",
		SecretKey:     "secret-key",
		EnableDebug:   true,
		PublicRead:    true,
		Region: "ap-northeast-1",
	})
	defer bridge.Disconnect()
	fmt.Printf("Connected to %s\n", bridge.Config().Provider)
```

## Upload files to S3 using Bifrost
``` go
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
And that's it! You have now successfully uploaded a file to Amazon S3 using Bifrost.

## Shipping multiple files to Amazon S3 via the rainbow bridge
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
- [Amazon S3 Documentation](https://docs.aws.amazon.com/s3/index.html)
- [Bifrost GitHub Repository](https://github.com/opensaucerer/bifrost)

We hope you found this guide helpful in using Bifrost with Amazon S3. If you have any questions or feedback, please don't hesitate to reach out to us!

