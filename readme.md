# Bifrost

Rainbow bridge for shipping your files to any cloud storage service.

![image](https://user-images.githubusercontent.com/59074379/226159115-1cfcb221-127f-4574-87ed-b74b4b2c4591.png)

# Problem Statement

You might just want to use 3 different cloud storage providers in your project. This means you'll need 3 different SDKs with 3 different implementations. That's just too much learning curve.

How about you ride with Thor on the Bifrost and easily transport your files to any cloud storage with the exact same function calls.

# Installation

```bash
go get github.com/opensaucerer/bifrost
```

# Usage

### Mounting a rainbow bridge to link with Google Cloud Storage

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

### Shipping a file to Google Cloud Storage via the rainbow bridge

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

### Shipping multiple files to Google Cloud Storage via the rainbow bridge

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

### Mount a rainbow bridge and ship a file to Amazon S3

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

### Shipping multiple files to Amazon S3 via the rainbow bridge

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

### Mount a rainbow bridge and ship a file to Pinata

```go
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

### Shipping multiple files to Pinata via the rainbow bridge

```go
// Upload multiple files
uploadedFiles, err := bridge.UploadMultiFile(bifrost.MultiFile{
	Files: []bifrost.File{
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

	GlobalOptions: map[string]interface{}{
		bifrost.OptPinata: map[string]interface{}{
			"cidVersion": 1,
		},
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
