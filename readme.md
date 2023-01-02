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

```go
package main

import (
	"fmt"

	"github.com/opensaucerer/bifrost"
)

func main() {
	gcs := bifrost.NewGoogleCloudStorage(&bifrost.GoogleCloudStorage{
		DefaultBucket: "default-bucket",
	})

	// upload a file to the default bucket
    _, err := gcs.Upload("file.txt", "file.txt")
    if err != nil {
        panic(err)
    }
}
```
