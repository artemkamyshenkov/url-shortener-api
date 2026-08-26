package httpapi

import "time"

type createURLRequest struct {
	URL string `json:"url"`
}

type urlResponse struct {
	ID        int64     `json:"id"`
	URL       string    `json:"url"`
	ShortCode string    `json:"short_code"`
	ShortURL  string    `json:"short_url"`
	CreatedAt time.Time `json:"created_at"`
}
