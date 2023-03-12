package request

import "net/http"

// Client is the interface for the request client
// Client is safe for concurrent use by multiple goroutines.
type Client struct {
	Http    *http.Client
	Request *http.Request
}
