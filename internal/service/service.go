package service

import "wb-task-L0/internal/repository"

type Order struct {
}

type Service struct {
	Order
}

func NewService(repos *repository.Repository) *Service {
	return &Service{}
}
