package types

type PinataAuthResponse struct {
	Error struct {
		Reason  string `json:"reason"`
		Details string `json:"details"`
	} `json:"error"`
	Message string `json:"message"`
}
