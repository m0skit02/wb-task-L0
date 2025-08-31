package repository

import (
	"context"
	"gorm.io/gorm"
	"wb-task-L0/pkg/models"
)

type Order interface {
	Create(order *models.Order) (string, error)
	GetAll() ([]models.Order, error)
	GetByID(id string) (models.Order, error)
	Delete(id string) error
	CreateOrderWithAssociations(context.Context, *models.Order) error
}

type Repository struct {
	Order
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Order: NewOrderRepo(db),
	}
}
