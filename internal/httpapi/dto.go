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

type urlStatsResponse struct {
	ShortCode      string     `json:"short_code"`
	Clicks         int64      `json:"clicks"`
	LastAccessedAt *time.Time `json:"last_accessed_at"`
}
