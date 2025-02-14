package entity

import "time"

type Pet struct {
	ID         uint64     `json:"id"`
	Name       string     `json:"name"`
	CategoryID uint64     `json:"category_id"`
	Category   Category   `json:"category"`
	Status     PetStatus  `json:"status"`
	DeletedAt  *time.Time `json:"deleted_at"`
}

type Category struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}
