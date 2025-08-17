package repository

import "database/sql"

type ProductRepository struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (pr *ProductRepository) Create() {

}
