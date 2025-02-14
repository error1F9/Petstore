package service

import (
	entity2 "Petstore/models/pet/entity"
	"context"
)

type PetServicer interface {
	Add(ctx context.Context, in PetAddIn) PetAddOut
	Update(ctx context.Context, in PetUpdateIn) PetUpdateOut
	FindByStatus(ctx context.Context, in PetFindByStatusIn) PetFindByStatusOut
	FindById(ctx context.Context, in PetFindByIdIn) PetFindByIdOut
	UpdateById(ctx context.Context, id PetUpdateByIdIn) PetUpdateByIdOut
	Delete(ctx context.Context, id PetDeleteIn) PetDeleteOut
}

type PetAddIn struct {
	Name     string
	Category entity2.Category
	Status   entity2.PetStatus
}

type PetAddOut struct {
	Err error
}

type PetUpdateIn struct {
	Pet entity2.Pet
}

type PetUpdateOut struct {
	Err error
}

type PetFindByStatusIn struct {
	Status entity2.PetStatus
}

type PetFindByStatusOut struct {
	Pets []entity2.Pet
	Err  error
}

type PetFindByIdIn struct {
	PetID uint64
}

type PetFindByIdOut struct {
	Pet entity2.Pet
	Err error
}

type PetUpdateByIdIn struct {
	Pet entity2.Pet
}

type PetUpdateByIdOut struct {
	Err error
}

type PetDeleteIn struct {
	PetID uint64
}

type PetDeleteOut struct {
	Err error
}
