package controller

import (
	"Petstore/models/pet/entity"
)

type PetResponse struct {
	Success bool `json:"success" example:"true"`
	Code    int  `json:"code,omitempty" example:"200"`
}

type PetResponseData struct {
	Success bool `json:"success" example:"true"`
	Code    int  `json:"code,omitempty" example:"200"`
	Data    Data `json:"data,omitempty"`
}

type PetFindByStatusResponse struct {
	Success bool     `json:"success" example:"true"`
	Code    int      `json:"code,omitempty" example:"200"`
	Data    DataPets `json:"data,omitempty"`
}

type Data struct {
	Message string     `json:"message,omitempty"`
	Pet     entity.Pet `json:"pet,omitempty"`
}

type DataPets struct {
	Message string       `json:"message,omitempty"`
	Pets    []entity.Pet `json:"pets,omitempty"`
}

type AddPetRequest struct {
	Name     string          `json:"name"`
	Category entity.Category `json:"category"`
	Status   string          `json:"status"`
}
