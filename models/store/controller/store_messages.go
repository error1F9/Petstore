package controller

import (
	"Petstore/models/store/entity"
)

type PlaceOrderResponse struct {
	Success bool `json:"success" example:"true"`
	Code    int  `json:"code,omitempty" example:"200"`
	Data    Data `json:"data,omitempty"`
}
type Data struct {
	Message string       `json:"message,omitempty"`
	Order   entity.Order `json:"pet,omitempty"`
}

type DeleteOrderResponse struct {
	Success bool   `json:"success" example:"true"`
	Code    int    `json:"code,omitempty" example:"200"`
	Message string `json:"message,omitempty"`
}

type InventoryResponse struct {
	s map[entity.OrderStatus]int
}
