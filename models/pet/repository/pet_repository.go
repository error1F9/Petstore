package repository

import (
	"Petstore/models/pet/entity"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Peter interface {
	Add(ctx context.Context, pet *entity.Pet) error
	Update(ctx context.Context, pet *entity.Pet) error
	FindByStatus(ctx context.Context, status entity.PetStatus) ([]entity.Pet, error)
	FindById(ctx context.Context, id uint64) (entity.Pet, error)
	UpdateById(ctx context.Context, pet *entity.Pet) error
	Delete(ctx context.Context, id uint64) error
}

type PetRepository struct {
	db *gorm.DB
}

func NewPetRepository(db *gorm.DB) Peter {
	return &PetRepository{db: db}
}

func (p *PetRepository) categoryExist(category entity.Category) (*entity.Category, error) {
	var existingCategory entity.Category
	result := p.db.Where("name=?", category.Name).First(&existingCategory)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &existingCategory, nil
}

func (p *PetRepository) Add(ctx context.Context, pet *entity.Pet) error {
	existingCategory, err := p.categoryExist(pet.Category)
	if err != nil {
		return errors.New("Category not found")
	}

	var categoryID uint64
	if existingCategory == nil {
		newCategory := entity.Category{Name: pet.Category.Name}
		if err = p.db.Table("categories").Create(&newCategory).Error; err != nil {
			return errors.New("Category could not be created")
		}
		categoryID = newCategory.ID
	} else {
		categoryID = existingCategory.ID
	}

	petToCreate := entity.Pet{
		Name:       pet.Name,
		CategoryID: categoryID,
		Status:     pet.Status,
	}

	err = p.db.Create(&petToCreate).Error

	return nil
}

func (p *PetRepository) Update(ctx context.Context, pet *entity.Pet) error {
	var existingPet entity.Pet
	if err := p.db.Where("id = ?", pet.ID).First(&existingPet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("pet with ID %d not found", pet.ID)
		}
		return err
	}

	existingCategory, err := p.categoryExist(pet.Category)
	if err != nil {
		return errors.New("Category not found")
	}

	var categoryID uint64
	if existingCategory == nil {
		newCategory := entity.Category{Name: pet.Category.Name}
		if err = p.db.Table("categories").Create(&newCategory).Error; err != nil {
			return errors.New("Category could not be created")
		}
		categoryID = newCategory.ID
	} else {
		categoryID = existingCategory.ID
	}

	petToUpdate := entity.Pet{
		ID:         pet.ID,
		Name:       pet.Name,
		CategoryID: categoryID,
		Status:     pet.Status,
	}

	err = p.db.Save(&petToUpdate).Error
	return nil
}

func (p *PetRepository) FindByStatus(ctx context.Context, status entity.PetStatus) ([]entity.Pet, error) {
	pets := make([]entity.Pet, 0)
	if err := p.db.Preload("Category").Where("status=?", status).Find(&pets).Error; err != nil {
		return nil, err
	}

	return pets, nil
}

func (p *PetRepository) FindById(ctx context.Context, id uint64) (entity.Pet, error) {
	pet := entity.Pet{}
	err := p.db.Preload("Category").Where("id=?", id).First(&pet).Error
	return pet, err
}

func (p *PetRepository) UpdateById(ctx context.Context, pet *entity.Pet) error {
	err := p.Update(ctx, pet)
	return err

}

func (p *PetRepository) Delete(ctx context.Context, id uint64) error {
	pet, err := p.FindById(ctx, id)
	if err != nil {
		return err
	}

	now := time.Now()
	pet.DeletedAt = &now

	err = p.db.Save(&pet).Error
	return err
}
