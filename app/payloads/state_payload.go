package payloads

type CreateStateRequest struct {
	Name string  `json:"name" binding:"required,min=2,max=100"`
	Code *string `json:"code" binding:"min=2,max=10"`
}
