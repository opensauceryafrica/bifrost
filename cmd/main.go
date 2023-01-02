package main

import (
	"fmt"

	"github.com/opensaucerer/bifrost"
)

func main() {
	gcs := bifrost.NewGoogleCloudStorage(&bifrost.GoogleCloudStorage{
		DefaultBucket: "default-bucket",
	})

	fmt.Println(gcs)
}
