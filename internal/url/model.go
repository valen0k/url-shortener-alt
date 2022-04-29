package url

import "time"

type Url struct {
	ID          string    `json:"-"`
	HashUrl     string    `json:"hash_url"`
	OriginalUrl string    `json:"original_url"`
	CreatedAt   time.Time `json:"-"`
}
