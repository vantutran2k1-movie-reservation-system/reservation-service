package models

type ResponseMeta struct {
	Limit   int     `json:"limit"`
	Offset  int     `json:"offset"`
	Total   int     `json:"total"`
	NextUrl *string `json:"next_url,omitempty"`
	PrevUrl *string `json:"prev_url,omitempty"`
}
