package pinata

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/opensaucerer/bifrost/shared/config"
	"github.com/opensaucerer/bifrost/shared/errors"
	"github.com/opensaucerer/bifrost/shared/types"
)

type MetaData struct {
	Name string `json:"name"`
}

type SuccessResponse struct {
	IpfsHash  string `json:"IpfsHash"`
	Timestamp string `json:"Timestamp"`
	PinSize   int64  `json:"PinSize"`
}

/*
	Upload File to Pinata IPFS
*/

func (p PinataCloud) UploadFile(path, filename string, options map[string]interface{}) (*types.UploadedFile, error) {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("file does not exist: %s", path),
			ErrorCode: errors.ErrBadRequest,
		}
	}

	url := config.URLPinataPinFile
	method := "POST"

	payload := &bytes.Buffer{} // buffer writer
	writer := multipart.NewWriter(payload)
	file, err := os.Open(path)
	if err != nil {
		return nil, &errors.BifrostError{
			Err:       err,
			ErrorCode: errors.ErrFileOperationFailed,
		}
	}
	defer file.Close()

	if err != nil {
		return nil, &errors.BifrostError{
			Err:       err,
			ErrorCode: errors.ErrFileOperationFailed,
		}
	}

	part1, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("file does not exist: %s", path),
			ErrorCode: errors.ErrBadRequest,
		}
	}

	_, err = io.Copy(part1, file)
	if err != nil {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("failed to copy multipart data"),
			ErrorCode: errors.ErrBadRequest,
		}
	}

	metaData := map[string]string{"name": filename}
	pinataMetadata, err := json.Marshal(metaData)
	if err != nil {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("failed to marshal metadata"),
			ErrorCode: errors.ErrBadRequest,
		}
	}

	pinataOptions, err := json.Marshal(options)

	if err != nil {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("failed to convert metadata to JSON"),
			ErrorCode: errors.ErrBadRequest,
		}
	}

	_ = writer.WriteField("pinataOptions", string(pinataOptions))
	_ = writer.WriteField("pinataMetadata", string(pinataMetadata))

	err = writer.Close()
	if err != nil {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("failed to close multipart writer"),
			ErrorCode: errors.ErrBadRequest,
		}
	}

	req, err := http.NewRequest(method, url, payload)

	if p.DefaultTimeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.DefaultTimeout))
		defer cancel()
		req = req.WithContext(ctx)
	}

	client := &http.Client{}

	if err != nil {
		return nil, &errors.BifrostError{
			Err:       err,
			ErrorCode: errors.ErrBadRequest,
		}
	}
	req.Header.Add("Authorization", "Bearer "+p.Config().PinataJWT)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)

	if err != nil {
		return nil, &errors.BifrostError{
			Err:       err,
			ErrorCode: errors.ErrBadRequest,
		}
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, &errors.BifrostError{
			Err:       err,
			ErrorCode: errors.ErrBadRequest,
		}
	}
	var responseCode int = res.StatusCode

	if responseCode != errors.Status200 {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("upload request failed with status %d", responseCode),
			ErrorCode: errors.ErrFileOperationFailed,
		}
	}

	res_ := SuccessResponse{}
	json.Unmarshal([]byte(body), &res_)

	return &types.UploadedFile{
		Size:    res_.PinSize,
		Name:    res_.IpfsHash,
		Preview: fmt.Sprintf("https://gateway.pinata.cloud/ipfs/%v", res_.IpfsHash),
	}, nil

}

/*
Disconnect closes the Pinata connection and returns an error if one occurs.

Disconnect should only be called when the connection is no longer needed.
*/
func (g *PinataCloud) Disconnect() error {
	return nil
}

func (p *PinataCloud) Config() *types.BridgeConfig {
	return &types.BridgeConfig{
		Provider:       p.Provider,
		DefaultTimeout: p.DefaultTimeout,
		PinataJWT:      p.PinataJWT,
		EnableDebug:    p.EnableDebug,
		UseAsync:       p.UseAsync,
		PublicRead:     p.PublicRead,
	}
}

func (g *PinataCloud) PreFlight() error {

	url := config.URLPinataAuth
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
		return &errors.BifrostError{
			Err:       err,
			ErrorCode: errors.ErrBadRequest,
		}
	}
	req.Header.Add("Authorization", "Bearer PINATA_JWT")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return &errors.BifrostError{
			Err:       err,
			ErrorCode: errors.ErrBadRequest,
		}

	}
	defer res.Body.Close()
	_, err = io.ReadAll(res.Body)

	if err != nil {
		return &errors.BifrostError{
			Err:       err,
			ErrorCode: errors.ErrBadRequest,
		}
	}

	return nil

}
