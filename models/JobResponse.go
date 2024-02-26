package models

type TextStream struct {
	RequestID    string `json:"_id"`
	ResponseType string `json:"response_type"`
	Chunk        string `json:"chunk"`
	Response     string `json:"response,omitempty"`
}
