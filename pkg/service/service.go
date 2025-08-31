package service

import (
	"context"
	"wb-task-L0/pkg/cache"
	"wb-task-L0/pkg/models"
	"wb-task-L0/pkg/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Order interface {
	Create(order *models.Order) (*models.Order, error)
	GetByID(id string) (models.Order, error)
	GetAll() ([]models.Order, error)
	Delete(id string) error
	CreateOrderWithAssociations(context.Context, *models.Order) error
}

type Service struct {
	Order
}

func NewService(repos *repository.Repository, cache *cache.OrderCache) *Service {
	return &Service{
		Order: NewOrderService(repos.Order, cache),
	}
}
