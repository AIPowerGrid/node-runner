package models

type Node struct {
	ID         string  `json:"_id"`
	VRAM       float64 `json:"vram"`
	CudaDevice string  `json:"cuda_device"`
	GPUCount   int     `json:"gpus"`
	Type       string  `json:"mode"` // either image or text gen
	Model      string  `json:"model,omitempty"`
}
