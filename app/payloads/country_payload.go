package payloads

type CreateCountryRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
	Code string `json:"code" binding:"required,len=2"`
}
