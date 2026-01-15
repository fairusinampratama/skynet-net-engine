package api

type CreateSecretRequest struct {
	User     string `json:"user" binding:"required"`
	Password string `json:"password" binding:"required"`
	Profile  string `json:"profile" binding:"required"`
	RemoteIP string `json:"remote_ip"`
	LocalIP  string `json:"local_ip"`
	Comment  string `json:"comment"`
}

type UpdatePlanRequest struct {
	Profile string `json:"profile" binding:"required"`
}
