package request

import (
	"fmt"
	"net/http"
	"time"
)

// BuildClient returns a new client for making requests
func BuildClient(url string, token string, timeout int64) *Client {
	h := &http.Client{}
	if timeout > 0 {
		h.Timeout = time.Duration(timeout) * time.Second
	}
	// these value will like be modified before any call to client.Do
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	return &Client{
		http:    h,
		request: req,
	}
}
