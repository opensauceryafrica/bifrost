package request

import (
	"fmt"
	"net/http"
	"time"

	"github.com/opensaucerer/bifrost/shared/config"
)

// NewClient returns a new client for making requests
func NewClient(url string, token string, timeout int64) *Client {
	h := &http.Client{}
	if timeout > 0 {
		h.Timeout = time.Duration(timeout) * time.Second
	}
	// these value will like be modified before any call to client.Do
	req, err := http.NewRequest(config.MethodGet, url, nil)
	if err != nil {
		return nil
	}
	req.Header.Add(config.ReqAuth, fmt.Sprintf(config.ReqBearer, token))
	return &Client{
		Http:    h,
		Request: req,
	}
}
