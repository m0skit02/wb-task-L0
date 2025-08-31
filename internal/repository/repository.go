package repository

import (
	"github.com/jmoiron/sqlx"
	"wb-task-L0/internal/models"
)

type Order interface {
	Create(order models.Order) (string, error)
	GetAll() ([]models.Order, error)
	GetByID(string) (*models.Order, error)
	Delete(string) error
}

type Repository struct {
	Order
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{}
}
