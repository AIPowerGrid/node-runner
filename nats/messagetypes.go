package nats

type ModelResponse struct {
	Model string `json:"model"`
	Type  string `json:"type"`
}

type JSResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}
