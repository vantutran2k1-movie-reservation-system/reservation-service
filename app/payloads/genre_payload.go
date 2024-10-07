package payloads

type CreateGenreRequest struct {
	Name string `json:"name" binding:"required,min=1,max=255"`
}
