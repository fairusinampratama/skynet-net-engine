package api

type IsolateRequest struct {
	IP     string `json:"ip" binding:"required"`
	Action string `json:"action" binding:"required,oneof=add remove"` // add or remove
	List   string `json:"list"` // Default to "ISOLATED" if empty
	Comment string `json:"comment"`
}
