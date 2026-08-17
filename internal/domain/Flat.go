package domain

import "time"

type Flat struct {
	ID         int64     `json:"id"`
	HouseID    int64     `json:"house_id"`
	Price      int       `json:"price"`
	Rooms      int       `json:"rooms"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	FlatNumber int       `json:"flat_number"`
}
