package entity

import (
	"encoding/json"
	"fmt"
)

type PetStatus string

const (
	PetStatusAvailable PetStatus = "available"
	PetStatusPending   PetStatus = "pending"
	PetStatusSold      PetStatus = "sold"
)

func (s PetStatus) String() string {
	return string(s)
}

func (s *PetStatus) UnmarshalJSON(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}

	switch str {
	case "available":
		*s = PetStatusAvailable
	case "pending":
		*s = PetStatusPending
	case "sold":
		*s = PetStatusSold
	default:
		return fmt.Errorf("invalid PetStatus: %s", str)
	}

	return nil
}

func (s PetStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
