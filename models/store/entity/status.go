package entity

import (
	"encoding/json"
	"fmt"
)

type OrderStatus string

const (
	OrderStatusPlaced    = "placed"
	OrderStatusApproved  = "approved"
	OrderStatusDelivered = "delivered"
)

func (s OrderStatus) String() string {
	return string(s)
}

func (s *OrderStatus) UnmarshalJSON(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}

	switch str {
	case "placed":
		*s = OrderStatusPlaced
	case "approved":
		*s = OrderStatusApproved
	case "delivered":
		*s = OrderStatusDelivered
	default:
		return fmt.Errorf("invalid PetStatus: %s", str)
	}

	return nil
}

func (s OrderStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
