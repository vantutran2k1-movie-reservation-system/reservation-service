package payloads

type CreateCountryRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
	Code string `json:"code" binding:"required,len=2"`
}

type CreateStateRequest struct {
	Name string  `json:"name" binding:"required,min=2,max=100"`
	Code *string `json:"code" binding:"min=2,max=10"`
}

type CreateCityRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
}
