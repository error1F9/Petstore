package service

import (
	"Petstore/models/user/entity"
	"Petstore/models/user/repository"
	"Petstore/pkg/token"
	"context"
	"errors"
	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserServicer interface {
	Login(ctx context.Context, in LoginIn) LoginOut
	Logout(ctx context.Context) error
	CreateUser(ctx context.Context, user *entity.User) (uint64, error)
	CreateWithArray(ctx context.Context, users []*entity.User) error
	GetUser(ctx context.Context, username string) (*entity.User, error)
	UpdateUser(ctx context.Context, username string, updateData *entity.User) (*entity.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type UserService struct {
	repository   repository.Userer
	tokenService token.JWTTokenService
}

func NewUserService(repository repository.Userer, tokenService token.JWTTokenService) UserServicer {
	return &UserService{repository: repository, tokenService: tokenService}
}

func (s *UserService) Login(ctx context.Context, in LoginIn) LoginOut {
	user, err := s.repository.GetUser(ctx, in.Username)
	if err != nil {
		return LoginOut{"", err}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password))
	if err != nil {
		return LoginOut{"", err}
	}

	newToken, err := s.tokenService.GenerateToken(user.Username, user.ID)
	if err != nil {
		return LoginOut{"", err}
	}

	err = s.repository.Login(ctx, in.Username)
	return LoginOut{newToken, err}
}

func (s *UserService) Logout(ctx context.Context) error {
	_, claims, _ := jwtauth.FromContext(ctx)
	username, ok := claims["username"].(string)
	if !ok {
		return errors.New("username not found")
	}
	err := s.repository.Logout(ctx, username)
	return err
}

func (s *UserService) CreateUser(ctx context.Context, user *entity.User) (uint64, error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	user.Password = string(hashedPassword)

	id, err := s.repository.CreateUser(ctx, user)
	return id, err
}

func (s *UserService) CreateWithArray(ctx context.Context, users []*entity.User) error {
	for _, user := range users {
		_, err := s.CreateUser(ctx, user)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *UserService) GetUser(ctx context.Context, username string) (*entity.User, error) {
	user, err := s.repository.GetUser(ctx, username)
	return user, err
}

func (s *UserService) UpdateUser(ctx context.Context, username string, updateData *entity.User) (*entity.User, error) {
	existingUser, err := s.repository.GetUser(ctx, username)
	if err != nil {
		return nil, err
	}

	if updateData.Email != "" && updateData.Email != existingUser.Email {
		if s.repository.EmailExist(updateData.Email) {
			return nil, errors.New("email already exists")
		}
	}

	if updateData.Username != "" && updateData.Username != existingUser.Username {
		if s.repository.UsernameExist(updateData.Username) {
			return nil, errors.New("username already exists")
		}
	}

	updateData.ID = existingUser.ID

	if updateData.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateData.Password), 10)
		if err != nil {
			return nil, err
		}
		updateData.Password = string(hashedPassword)
	} else {
		updateData.Password = existingUser.Password
	}

	if err = s.repository.UpdateUser(ctx, updateData); err != nil {
		return nil, err
	}

	return updateData, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	err := s.repository.DeleteUser(ctx, id)
	return err
}
