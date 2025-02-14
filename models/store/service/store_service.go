package service

import (
	entity2 "Petstore/models/store/entity"
	"Petstore/models/store/repository"
	"context"
)

type OrderServicer interface {
	Inventory(ctx context.Context) map[entity2.OrderStatus]int
	PlaceOrder(ctx context.Context, order entity2.Order) (entity2.Order, error)
	FindOrderById(ctx context.Context, id uint64) (entity2.Order, error)
	DeleteById(ctx context.Context, id uint64) error
}

type OrderService struct {
	repository repository.Storer
}

func NewOrderService(repository repository.Storer) *OrderService {
	return &OrderService{repository: repository}
}

func (s *OrderService) Inventory(ctx context.Context) map[entity2.OrderStatus]int {
	return s.repository.Inventory(ctx)
}

func (s *OrderService) PlaceOrder(ctx context.Context, order entity2.Order) (entity2.Order, error) {
	return s.repository.PlaceOrder(ctx, order)
}

func (s *OrderService) FindOrderById(ctx context.Context, id uint64) (entity2.Order, error) {
	return s.repository.FindOrderById(ctx, id)
}

func (s *OrderService) DeleteById(ctx context.Context, id uint64) error {
	return s.repository.DeleteById(ctx, id)
}
