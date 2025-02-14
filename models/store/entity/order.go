package entity

import "time"

type Order struct {
	ID       uint64      `json:"id"`
	PetID    uint64      `json:"pet_id"`
	Quantity int         `json:"quantity"`
	ShipDate *time.Time  `json:"ship_date"`
	Status   OrderStatus `json:"status"`
	Complete bool        `json:"complete"`
}
