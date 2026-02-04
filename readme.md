# Bifrost

Rainbow bridge for shipping your files to any cloud storage service with the same function calls.

<img src="https://user-images.githubusercontent.com/59074379/226159115-1cfcb221-127f-4574-87ed-b74b4b2c4591.png" width="1000" />

# Table of contents

- [Bifrost](#bifrost)
- [Problem Statement](#problem-statement)
- [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)
- [Changelog](#changelog)
- [Contributors](#contributors)

# Problem Statement

Many projects need to store files in the cloud and different projects might use different cloud storage providers or, sometimes, multiple cloud providers all at once. Using different SDKs with different implementations for each provider can be tedious and time-consuming. Bifrost aims to simplify the process of working with multiple cloud storage providers by providing a consistent API for all of them.

# Installation

To install the Bifrost package, run the following command in your terminal:

```bash
go get github.com/opensaucerer/bifrost
```

## Usage

```go
package main

import (
	"fmt"
	"os"
	"github.com/opensaucerer/bifrost"
)

// mount a bridge to gcs
bridge, _ := bifrost.NewRainbowBridge(&bifrost.BridgeConfig{
	DefaultBucket:   "bifrost",
	DefaultTimeout:  10,
	Provider:        bifrost.GoogleCloudStorage,
	CredentialsFile: "/path/to/service/account/json", // this is not required if you are using google's default credentials
	EnableDebug:     true,
	PublicRead:      true,
})
defer bridge.Disconnect()
fmt.Printf("Connected to %s\n", bridge.Config().Provider)

// upload a file to gcs using the bridge
guf, _ := bridge.UploadFile(bifrost.File{
	Path:     "../shared/image/aand.png",
	Filename: "a_and_ampersand.png",
	Options: map[string]interface{}{
		bifrost.OptMetadata: map[string]string{
			"originalname": "aand.png",
		},
	},
})
fmt.Printf("Uploaded file %s to GCS at: %s\n", guf.Name, guf.Preview)
```

The above example clearly demonstrates the speed, simplicity, and ease of use that Bifrost offers. Now you know what it feels like to ride with Thor!

If you want to learn more about how Bifrost is creating different methods to make it easier to use different cloud providers, you can follow these links:

- [Google Cloud Storage (GCS)](gcs/doc.md)
- [Amazon S3](s3/doc.md)
- [Pinata Cloud](pinata/doc.md)
- [Wasabi Cloud](wasabi/doc.md)

# Variants

Bifrost also exists in other forms and languages and you are free to start a new variant of bifrost in any other form or language of your choice. For now, below are the know variants of bifrost.

- [x] [Bifrost CLI](https://github.com/showbaba/bifrost-cli)
- [x] [Bifrost in Python](https://github.com/ifihan/byfrost)

# Contributing

Bifrost is an open source project and we welcome contributions of all kinds. Please read our [contributing guide](./contributing.md) to learn about our development process, how to propose bug fixes and improvements, and how to build and test your changes to Bifrost.

# License

Bifrost is [MIT licensed](./LICENSE).

# Changelog

See [changelog](./changelog.md) for more details.

# Contributors

<a href="https://github.com/opensaucerer/bifrost/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=opensaucerer/bifrost" />
</a>

Made with [contrib.rocks](https://contrib.rocks).
