package shortener

import "time"

type URL struct {
	ID          int64
	OriginalURL string
	ShortCode   string
	CreatedAt   time.Time
}
