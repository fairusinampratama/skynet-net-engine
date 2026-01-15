package models

type SystemResource struct {
	Uptime       string `json:"uptime"`
	CPU          string `json:"cpu"` // Load in %
	TotalMemory  int64  `json:"total_memory"`
	FreeMemory   int64  `json:"free_memory"`
	BoardName    string `json:"board_name"`
	Version      string `json:"version"`
}

type TrafficStats struct {
	Name string `json:"name"` // User or Interface name
	RX   int64  `json:"rx"`   // bps
	TX   int64  `json:"tx"`   // bps
}
