package repository

import (
	"Petstore/models/user/entity"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type Userer interface {
	CreateUser(ctx context.Context, user *entity.User) (uint64, error)
	GetUser(ctx context.Context, username string) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, id string) error
	Login(ctx context.Context, username string) error
	Logout(ctx context.Context, username string) error
	UserExist(user entity.User) error
	EmailExist(email string) bool
	UsernameExist(username string) bool
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) Userer {
	return &UserRepository{db: db}
}

func (r *UserRepository) UserExist(user entity.User) error {
	var count int64
	if err := r.db.Table("users").Where("email = ?", user.Email).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("user with email %s already exists", user.Email)
	}

	if err := r.db.Table("users").Where("username = ?", user.Username).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("user with email %s already exists", user.Email)
	}
	return nil
}

func (r *UserRepository) EmailExist(email string) bool {
	err := r.db.Where("email=?", email).First(&entity.User{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}

func (r *UserRepository) UsernameExist(username string) bool {
	err := r.db.Where("username=?", username).First(&entity.User{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}

func (r *UserRepository) CreateUser(ctx context.Context, user *entity.User) (uint64, error) {
	if err := r.UserExist(*user); err != nil {
		return 0, err
	}
	if err := r.db.Create(&user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (r *UserRepository) GetUser(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	err := r.db.Table("users").Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	if err := r.db.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, username string) error {
	user, err := r.GetUser(ctx, username)
	if err != nil {
		return err
	}
	err = r.db.Delete(user, "username = ?", username).Error
	return err
}

func (r *UserRepository) Login(ctx context.Context, username string) error {
	err := r.db.Table("users").Where("username = ?", username).Update("user_status", 1).Error
	return err
}

func (r *UserRepository) Logout(ctx context.Context, username string) error {
	err := r.db.Table("users").Where("username = ?", username).Update("user_status", 0).Error
	return err
}
