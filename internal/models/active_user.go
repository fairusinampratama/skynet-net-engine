package models

type ActiveUser struct {
	Name      string `json:"name"`
	Address   string `json:"address"` // IP Address
	CallerID  string `json:"caller_id"` // MAC Address
	Uptime    string `json:"uptime"`
	RouterID  int    `json:"router_id"`
}
