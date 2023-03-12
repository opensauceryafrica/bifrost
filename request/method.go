package request

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/opensaucerer/bifrost/shared/config"
	"github.com/opensaucerer/bifrost/shared/errors"
	"github.com/opensaucerer/bifrost/shared/types"
)

func (c *Client) PinataPreflight() error {
	c.request.URL, _ = c.request.URL.Parse(config.URLPinataAuth)
	c.request.Method = config.MethodGet
	res, err := c.http.Do(c.request)
	if err != nil {
		fmt.Println(err)
		return &errors.BifrostError{
			Err:       err,
			ErrorCode: errors.ErrBadRequest,
		}

	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return &errors.BifrostError{
			Err:       err,
			ErrorCode: errors.ErrBadRequest,
		}
	}
	var par types.PinataAuthResponse
	err = json.Unmarshal(b, &par)
	if err != nil {
		return &errors.BifrostError{
			Err:       err,
			ErrorCode: errors.ErrBadRequest,
		}
	}
	if par.Message == "" {
		return &errors.BifrostError{
			Err:       fmt.Errorf(par.Error.Reason),
			ErrorCode: errors.ErrUnauthorized,
		}
	}
	return nil
}
