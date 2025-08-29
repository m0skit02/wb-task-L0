package repository

import "github.com/jmoiron/sqlx"

type Order struct {
}

type Repository struct {
	Order
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{}
}
