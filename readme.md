# Bifrost

Rainbow bridge for shipping your files to any cloud storage service with the same function calls.

<img src="https://user-images.githubusercontent.com/59074379/226159115-1cfcb221-127f-4574-87ed-b74b4b2c4591.png" width="1000" />

# Table of contents

- [Bifrost](#bifrost)
- [Problem Statement](#problem-statement)
  - [Google Cloud Storage(GCS)](#google-cloud-storagegcs)
  - [Pinata](#pinata)
- [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)
- [Changelog](#changelog)
- [Contributors](#contributors)

# Problem Statement
Many projects need to store files in the cloud, and different projects might use different cloud storage providers. Using different SDKs with different implementations for each provider can be tedious and time-consuming. Bifrost aims to simplify the process of working with multiple cloud storage providers by providing a consistent API for all of them. To better understand how Bifrost solves this problem, let's take a look at two separate code samples for GCS and Pinata using conventional means and how Bifrost eases the steps. 


## Google Cloud Storage(GCS)
Without Bifrost, the process of uploading a file to GCS using the Google Cloud Storage client library for Go would typically involve the following steps:
 ``` go
package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"cloud.google.com/go/storage"
)

func main() {
	ctx := context.Background()

	// create a client
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// open the file you want to upload
	file, err := os.Open("path/to/your/file")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// create a bucket object
	bucket := client.Bucket("your-bucket-name")

	// create an object handle
	object := bucket.Object("destination/file/name")

	// create a writer to upload the file
	writer := object.NewWriter(ctx)

	// copy the contents of the file to the object
	if _, err := io.Copy(writer, file); err != nil {
		log.Fatalf("Failed to upload file: %v", err)
	}

	// close the writer to finalize the upload
	if err := writer.Close(); err != nil {
		log.Fatalf("Failed to close writer: %v", err)
	}

	fmt.Println("File uploaded successfully!")
}

```
With Bifrost, the above process can be simplified to the following steps:
``` go
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
This Go code uploads a file named "aand.png" located at "../shared/image" to Google Cloud Storage via the Rainbow bridge.

## Pinata
If you don't use Bifrost, the usual way of uploading a file to Pinata involves going through the following steps:
``` go
package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

func main() {
	// Set the API key and secret key
	apiKey := "your-api-key"
	secretApiKey := "your-secret-api-key"

	// Open the file to be uploaded
	file, err := os.Open("path/to/file")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Prepare the request body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		fmt.Println("Error creating form file:", err)
		return
	}
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("Error copying file:", err)
		return
	}
	err = writer.Close()
	if err != nil {
		fmt.Println("Error closing writer:", err)
		return
	}

	// Prepare the request
	url := "https://api.pinata.cloud/pinning/pinFileToIPFS"
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("pinata_api_key", apiKey)
	req.Header.Add("pinata_secret_api_key", secretApiKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// Print the response
	fmt.Println(string(respBody))
}

```
With Bifrost, the above process can be simplified to the following steps:
``` go
package main

import (
	"fmt"
	"os"

	"github.com/opensaucerer/bifrost"
)

func main() {
	bridge, _ := bifrost.NewRainbowBridge(&bifrost.BridgeConfig{
		Provider:    bifrost.PinataCloud,
		PinataJWT:   os.Getenv("PINATA_JWT"),
		EnableDebug: true,
		PublicRead:  true,
	})
	defer bridge.Disconnect()
	fmt.Printf("Connected to %s\n", bridge.Config().Provider)
	// Upload a file
	uploadedFile, err := bridge.UploadFile(bifrost.File{
		Path:     "../shared/image/aand.png",
		Filename: "pinata_aand.png",
		Options: map[string]interface{}{
			bifrost.OptPinata: map[string]interface{}{
				"cidVersion": 1,
			},
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

# Installation
To install the Bifrost package, run the following command in your terminal:
```bash
go get github.com/opensaucerer/bifrost
```
# Usage
If you want to learn more about how Bifrost is creating different methods to make it easier to use different cloud providers, you can follow these links: 
- [Google Cloud Storage (GCS)](gcs\doc.md)
- [Amazon S3](s3\doc.md)
- [Pinata](pinata\doc.md)

# Contributing

Bifrost is an open source project and we welcome contributions of all kinds. Please read our [contributing guide](./contributing.md) to learn about our development process, how to propose bugfixes and improvements, and how to build and test your changes to Bifrost.

# License

Bifrost is [MIT licensed](./LICENSE).

# Changelog

See [changelog](./changelog.md) for more details.

# Contributors

<a href="https://github.com/opensaucerer/bifrost/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=opensaucerer/bifrost" />
</a>

Made with [contrib.rocks](https://contrib.rocks).
