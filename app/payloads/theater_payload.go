package payloads

type CreateTheaterRequest struct {
	Name string `json:"name" binding:"required,min=2,max=255"`
}
