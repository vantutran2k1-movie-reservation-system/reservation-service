package payloads

type CreateCityRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
}
