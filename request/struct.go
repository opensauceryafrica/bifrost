package request

import "net/http"

// Client is the interface for the request client
type Client struct {
	http    *http.Client
	request *http.Request
}
