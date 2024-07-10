// Bifrost interface for Pinata Cloud
package pinata

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/opensaucerer/bifrost/shared/config"
	"github.com/opensaucerer/bifrost/shared/errors"
	"github.com/opensaucerer/bifrost/shared/types"
)

/*
UploadFile uploads a file to Pinata and returns an error if one occurs.
*/
func (p *PinataCloud) UploadFile(fileFace interface{}) (*types.UploadedFile, error) {

	// assert that the fileFace is of type bifrost.File
	bFile, ok := fileFace.(types.File)
	if !ok {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("argument must be of type bifrost.File"),
			ErrorCode: errors.ErrBadRequest,
		}
	}

	// validate struct
	if err := bFile.Validate(); err != nil {
		return nil, &errors.BifrostError{
			Err:       err,
			ErrorCode: errors.ErrInvalidParameters,
		}
	}

	if !p.IsConnected() {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("no active Pinata client"),
			ErrorCode: errors.ErrClientError,
		}
	}

	if bFile.Path != "" {
		// verify that file exists
		if _, err := os.Stat(bFile.Path); os.IsNotExist(err) {
			return nil, &errors.BifrostError{
				Err:       fmt.Errorf("file does not exist: %s", bFile.Path),
				ErrorCode: errors.ErrBadRequest,
			}
		}

		// build the request params
		if bFile.Filename == "" {
			bFile.Filename = filepath.Base(bFile.Path)
		}
	}

	var param types.Param = types.Param{
		Files: []types.ParamFile{
			{
				Path:   bFile.Path,
				Handle: bFile.Handle,
				Key:    "file",
				Name:   bFile.Filename,
			},
		},
		Data: []types.ParamData{},
	}

	// configure upload options
	for k, v := range bFile.Options {
		switch k {
		// pinataOptions
		case config.OptPinata:
			if v, ok := v.(map[string]interface{}); ok {
				opt, err := json.Marshal(v)
				if err != nil {
					return nil, &errors.BifrostError{
						Err:       fmt.Errorf("failed to marshal pinata options: %s", err.Error()),
						ErrorCode: errors.ErrBadRequest,
					}
				}
				param.Data = append(param.Data, types.ParamData{
					Key:   config.OptPinata,
					Value: string(opt),
				})
			}

		// pinataMetadata
		case config.PinataCloud + strings.ToUpper(config.OptMetadata[:1]) + config.OptMetadata[1:]:
			if v, ok := v.(map[string]interface{}); ok {
				m, err := json.Marshal(v)
				if err != nil {
					return nil, &errors.BifrostError{
						Err:       fmt.Errorf("failed to marshal pinata metadata: %s", err.Error()),
						ErrorCode: errors.ErrBadRequest,
					}
				}
				param.Data = append(param.Data, types.ParamData{
					Key:   config.PinataCloud + strings.ToUpper(config.OptMetadata[:1]) + config.OptMetadata[1:],
					Value: string(m),
				})
			}
		}
	}

	res, err := p.Client.PostForm(config.URLPinataPinFile, param)
	if err != nil {
		return nil, &errors.BifrostError{
			Err:       err,
			ErrorCode: errors.ErrBadRequest,
		}
	}

	var obj types.PinataPinFileResponse
	if err := json.Unmarshal(res, &obj); err != nil {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("failed to unmarshal response: %s", err.Error()),
			ErrorCode: errors.ErrBadRequest,
		}
	}

	if obj.Error != "" {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("failed to upload file: %s", obj.Error),
			ErrorCode: errors.ErrBadRequest,
		}
	}

	return &types.UploadedFile{
		Size:           obj.PinSize,
		CID:            obj.IpfsHash,
		Preview:        fmt.Sprintf(config.URLPinataGateway, obj.IpfsHash),
		ProviderObject: obj,
		Name:           filepath.Base(bFile.Path),
		Path:           bFile.Path,
		URL:            fmt.Sprintf(config.URLPinataGateway, obj.IpfsHash),
	}, nil
}

/*
Disconnect closes the Pinata connection and returns an error if one occurs.

Disconnect should only be called when the connection is no longer needed.
*/
func (p *PinataCloud) Disconnect() error {
	if p.IsConnected() {
		p.Client = nil
	}
	return nil
}

// Config returns the Pinata Cloud configuration.
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

// Preflight attempts to authenticate with Pinata and returns an error if one occurs.
func (p *PinataCloud) Preflight() error {
	if !p.IsConnected() {
		return &errors.BifrostError{
			Err:       fmt.Errorf("no active Pinata client"),
			ErrorCode: errors.ErrClientError,
		}
	}
	// copy the request
	req := p.Client.Request.Clone(p.Client.Request.Context())
	req.URL, _ = req.URL.Parse(config.URLPinataAuth)
	req.Method = config.MethodGet
	res, err := p.Client.Http.Do(req)
	if err != nil {
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

// IsConnected returns true if the Pinata Cloud connection is non nil.
func (p *PinataCloud) IsConnected() bool {
	return p.Client != nil
}

/*
UploadFolder uploads a folder to the provider storage and returns an error if one occurs.

Note: for some providers, UploadFolder requires that a default bucket be set in bifrost.BridgeConfig.
*/
func (p *PinataCloud) UploadFolder(foldFace interface{}) ([]*types.UploadedFile, error) {
	return nil, nil
}

// UploadMultiFile
func (p *PinataCloud) UploadMultiFile(multiFace interface{}) ([]*types.UploadedFile, error) {

	// assert that the multiFace is of type bifrost.File
	multiFile, ok := multiFace.(types.MultiFile)
	if !ok {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("argument must be of type bifrost.MultiFile"),
			ErrorCode: errors.ErrBadRequest,
		}
	}

	// validate struct
	if err := multiFile.Validate(); err != nil {
		return nil, &errors.BifrostError{
			Err:       err,
			ErrorCode: errors.ErrInvalidParameters,
		}
	}

	if !p.IsConnected() {
		return nil, &errors.BifrostError{
			Err:       fmt.Errorf("no active Google Cloud Storage client"),
			ErrorCode: errors.ErrClientError,
		}
	}

	uploadedFiles := make([]*types.UploadedFile, 0, len(multiFile.Files))

	// @TODO: add concurrency when UseAsync is true
	for _, file := range multiFile.Files {
		if multiFile.GlobalOptions != nil {
			// merge global options with file options
			for k, v := range multiFile.GlobalOptions {
				// don't override if file has option already
				if _, ok := file.Options[k]; !ok {
					file.Options[k] = v
				}
			}
		}

		uploadedFile, err := p.UploadFile(file)
		if err != nil {
			if p.EnableDebug {
				// log failed file and continue
				log.Printf("Upload for file at path %s failed with err: %s\n", file.Path, err.Error())
			}
			uploadedFiles = append(uploadedFiles, &types.UploadedFile{Error: err, Name: file.Filename, Path: file.Path})
			continue
		}
		uploadedFiles = append(uploadedFiles, uploadedFile)
	}

	return uploadedFiles, nil
}

/*
DeleteObject deletes an object from an array of buckets in the provider's storage and returns an error if one occurs.

Note: DeleteObject requires that an object and an array of buckets to be set in bifrost.BridgeConfig.
*/
func (p *PinataCloud) DeleteObject() error {
	return nil
}
