package model

type ServerResponse struct {
	Message      string `json:"message,omitempty"`
	ResponseCode int    `json:"responseCode,omitempty"`
}
