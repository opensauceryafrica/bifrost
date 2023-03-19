package types

// PinataAuthResponse is the response from Pinata Cloud when authenticating.
type PinataAuthResponse struct {
	Error struct {
		Reason  string `json:"reason"`
		Details string `json:"details"`
	} `json:"error"`
	Message string `json:"message"`
}

// PinataPinFileResponse is the response from Pinata Cloud when pinning a file.
type PinataPinFileResponse struct {
	IpfsHash  string `json:"IpfsHash"`
	Timestamp string `json:"Timestamp"`
	PinSize   int64  `json:"PinSize"`
	Error     string `json:"error"`
}
