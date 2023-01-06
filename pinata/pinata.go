package pinata

import (
	"bytes"
	"fmt"
	"github.com/opensaucerer/bifrost/shared/errors"
	"github.com/opensaucerer/bifrost/shared/types"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	 "encoding/json"
)

type MetaData struct {
    Name   string      `json:"name"`
}

type SuccessResponse struct {
	 IpfsHash   string `json:"IpfsHash"`
	 Timestamp 	string `json:"Timestamp"`
	 PinSize 		int64 `json:"PinSize"`
}

/*
	Upload File to pinata IPFS
*/

func (p PinataIPFSStorage) UploadFile(path, filename string, options map[string]interface{}) (*types.UploadedFile, error) {

  // validate JWT	
	success,err := p.TestAuthentication()
	if err != nil {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("problem authenticating: %v", err),
			ErrorCode: errors.ErrBadRequest,
		}
	}
	if !success {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("authentication failed"),
			ErrorCode: errors.ErrBadRequest,
		}
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("file does not exist: %s", path),
			ErrorCode: errors.ErrBadRequest,
		}
	}

	url := "https://api.pinata.cloud/pinning/pinFileToIPFS"
	method := "POST"

	payload := &bytes.Buffer{} // buffer writer
	writer := multipart.NewWriter(payload)
	file, fileOpenError := os.Open(path)

	if fileOpenError != nil {
		return nil, &errors.BifrostError{
			Err:       fileOpenError,
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

	_, errFile1 := io.Copy(part1, file)
	if errFile1 != nil {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("failed to copy multipart data"),
			ErrorCode: errors.ErrBadRequest,
		}
	}

	defer file.Close()

	 metaData := map[string]string{"name": filename}
	 pinataMetadata,err := json.Marshal(metaData)
	 if err!= nil {
		 return nil, &errors.BifrostError{
        Err:       fmt.Errorf("failed to marshal metadata"),
        ErrorCode: errors.ErrBadRequest,
      }
	 }

	 pinataOptions,err := json.Marshal(options)


	 if err!= nil {
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

	client := &http.Client {
  }
  req, err := http.NewRequest(method, url, payload)

  if err != nil {
    return nil, &errors.BifrostError{
			Err:       err,
      ErrorCode: errors.ErrBadRequest,
		}
  }
	 req.Header.Add("Authorization", "Bearer " + p.Config().PinataJWT)
	 

  req.Header.Set("Content-Type", writer.FormDataContentType())
  res, err := client.Do(req)
	fmt.Println(req.Body)

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
	fmt.Println(responseCode)

	if responseCode != 200 {
		return nil, &errors.BifrostError{
      Err:       fmt.Errorf("upload request failed with status %d", responseCode),
      ErrorCode: errors.ErrBadRequest,
    }
	}

		res_ := SuccessResponse{}
    json.Unmarshal([]byte(body), &res_)

		return &types.UploadedFile{
			Size: res_.PinSize,
			Name: res_.IpfsHash,
			Preview: "https://gateway.pinata.cloud/ipfs/"+res_.IpfsHash ,
	}, nil

}

/*

 */
func (g *PinataIPFSStorage) Disconnect() error {
	return nil
}

func (g *PinataIPFSStorage) Config() *types.BridgeConfig {
	return &types.BridgeConfig{
		Provider:        g.Provider,
		DefaultBucket:   g.DefaultBucket,
		CredentialsFile: g.CredentialsFile,
		Project:         g.Project,
		DefaultTimeout:  g.DefaultTimeout,
		EnableDebug:     g.EnableDebug,
		UseAsync:        g.UseAsync,
		PinataJWT:			 g.PinataJWT,
	}
}

func (g *PinataIPFSStorage) TestAuthentication() (bool, error) {

	url := "https://api.pinata.cloud/data/testAuthentication"
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
		return false, &errors.BifrostError{
			Err:       fmt.Errorf("failed to copy multipart data"),
			ErrorCode: errors.ErrBadRequest,
		}
	}
	req.Header.Add("Authorization", "Bearer PINATA_JWT")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false, &errors.BifrostError{
			Err:       fmt.Errorf("failed to copy multipart data"),
			ErrorCode: errors.ErrBadRequest,
		}

	}
	defer res.Body.Close()
	_, err = io.ReadAll(res.Body)

	if err != nil {
		fmt.Println(err)
		return false, &errors.BifrostError{
			Err:       fmt.Errorf("failed to read body stream"),
			ErrorCode: errors.ErrBadRequest,
		}
	}

	return true, nil

}
