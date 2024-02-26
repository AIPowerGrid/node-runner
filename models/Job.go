package models

type Job struct {
	ID     string `json:"_id"`
	Type   string `json:"type"`
	Prompt string `json:"prompt"`
	Seed   int64  `json:"seed,omitempty"`
	Model  string `json:"model"`
	NodeID string `json:"nodeID,omitempty"`
}
