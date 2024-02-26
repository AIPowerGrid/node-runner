package models

type Config struct {
	Nodes     []Node `json:"nodes"`
	MachineID string `json:"machineId"`
	OwnerID   string `json:"ownerId"`
}
