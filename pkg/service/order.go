package service

import (
	"context"
	"wb-task-L0/pkg/cache"
	"wb-task-L0/pkg/models"
	"wb-task-L0/pkg/repository"
)

type OrderService struct {
	repo  repository.Order
	cache *cache.OrderCache
}

func NewOrderService(repo repository.Order, cache *cache.OrderCache) *OrderService {
	return &OrderService{
		repo:  repo,
		cache: cache,
	}
}

func (s *OrderService) Create(order *models.Order) (*models.Order, error) {
	uid, err := s.repo.Create(order)
	if err != nil {
		return nil, err
	}

	order.OrderUID = uid
	s.cache.Set(*order)

	return order, nil
}

func (s *OrderService) CreateOrderWithAssociations(ctx context.Context, order *models.Order) (*models.Order, error) {
	if err := s.repo.CreateOrderWithAssociations(ctx, order); err != nil {
		return nil, err
	}

	s.cache.Set(*order)

	return order, nil
}

func (s *OrderService) GetAll() ([]models.Order, error) {
	return s.repo.GetAll()
}

func (s *OrderService) GetByID(id string) (models.Order, error) {
	if order, ok := s.cache.Get(id); ok {
		return order, nil
	}

	order, err := s.repo.GetByID(id)
	if err != nil {
		return models.Order{}, err
	}

	s.cache.Set(order)

	return s.repo.GetByID(id)
}

func (s *OrderService) Delete(id string) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	s.cache.Delete(id)

	return nil
}
