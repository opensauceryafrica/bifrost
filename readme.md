# Bifrost

Rainbow bridge for shipping your files to any cloud storage service.

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
		DefaultBucket:   "bucket-name",
		DefaultTimeout:  10,
		CredentialsFile: "/path/to/json",
		Provider:        bifrost.GoogleCloudStorage,
		EnableDebug:     true,
	})
	if err != nil {
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
uploadedFile, err := bridge.UploadFile("./cmd/0000a_hair.jpg", "0000000_hair.jpg")
if err != nil {
    fmt.Println(err.(bifrost.Error).Code(), err)
    return
}

fmt.Printf("Uploaded file: %s to %s\n", uploadedFile.Name, uploadedFile.Preview)
```
