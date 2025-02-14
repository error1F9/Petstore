package repository

import (
	entity2 "Petstore/models/store/entity"
	"context"
	"gorm.io/gorm"
	"log"
)

type Storer interface {
	Inventory(ctx context.Context) map[entity2.OrderStatus]int
	PlaceOrder(ctx context.Context, order entity2.Order) (entity2.Order, error)
	FindOrderById(ctx context.Context, id uint64) (entity2.Order, error)
	DeleteById(ctx context.Context, id uint64) error
}

type StoreRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) Storer {
	return &StoreRepository{db: db}
}

func (r *StoreRepository) Inventory(ctx context.Context) map[entity2.OrderStatus]int {
	orders := make(map[entity2.OrderStatus]int)

	var results []struct {
		Status entity2.OrderStatus `json:"status"`
		Count  int                 `json:"count"`
	}

	err := r.db.Table("orders").Select("status", "COUNT(status) as count").Group("status").Scan(&results).Error
	if err != nil {
		log.Println(err)
		return orders
	}
	for _, result := range results {
		orders[result.Status] = result.Count
	}
	return orders
}

func (r *StoreRepository) PlaceOrder(ctx context.Context, order entity2.Order) (entity2.Order, error) {
	err := r.db.Table("orders").Create(&order).Error
	return order, err
}

func (r *StoreRepository) FindOrderById(ctx context.Context, id uint64) (entity2.Order, error) {
	var order entity2.Order
	err := r.db.First(&order, id).Error
	return order, err
}
func (r *StoreRepository) DeleteById(ctx context.Context, id uint64) error {
	err := r.db.Delete(&entity2.Order{}, id).Error
	return err
}
