package service

import (
	"Petstore/models/pet/entity"
	"Petstore/models/pet/repository"
	"context"
)

type PetService struct {
	repository repository.Peter
}

func NewPetService(repository repository.Peter) *PetService {
	return &PetService{repository: repository}
}

func (p *PetService) Add(ctx context.Context, in PetAddIn) PetAddOut {
	pet := &entity.Pet{
		Name:     in.Name,
		Category: in.Category,
		Status:   in.Status,
	}

	err := p.repository.Add(ctx, pet)
	return PetAddOut{err}
}

func (p *PetService) Update(ctx context.Context, in PetUpdateIn) PetUpdateOut {
	pet := in.Pet
	err := p.repository.Update(ctx, &pet)
	return PetUpdateOut{err}
}

func (p *PetService) FindByStatus(ctx context.Context, in PetFindByStatusIn) PetFindByStatusOut {
	pets, err := p.repository.FindByStatus(ctx, in.Status)
	return PetFindByStatusOut{Pets: pets, Err: err}
}

func (p *PetService) FindById(ctx context.Context, in PetFindByIdIn) PetFindByIdOut {
	pet, err := p.repository.FindById(ctx, in.PetID)
	return PetFindByIdOut{pet, err}
}

func (p *PetService) UpdateById(ctx context.Context, in PetUpdateByIdIn) PetUpdateByIdOut {
	err := p.repository.UpdateById(ctx, &in.Pet)
	return PetUpdateByIdOut{err}
}

func (p *PetService) Delete(ctx context.Context, in PetDeleteIn) PetDeleteOut {
	err := p.repository.Delete(ctx, in.PetID)
	return PetDeleteOut{err}
}
