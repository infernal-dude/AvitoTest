package domain

import (
	"time"
)

type House struct {
	ID              int64     `json:"id"`
	Address         string    `json:"address"`
	Year            int       `json:"year"`
	Developer       *string   `json:"developer"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	LastFlatAddedAt time.Time `json:"last_flat_added_at"`
}
