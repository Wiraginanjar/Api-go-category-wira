package repositories

import (
	"database/sql"
	"errors"
	"tugas-go/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (repo *ProductRepository) GetAll(name string) ([]models.Productsy, error) {
	query := "SELECT p.id, p.name, p.price, p.stock, c.name FROM product p JOIN category c ON p.category_id = c.id"
	args := []interface{}{}
	if name != "" {
		query += " WHERE p.name ILIKE $1"
		args = append(args, "%"+name+"%")
	}

	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Productsy, 0)
	for rows.Next() {
		var p models.Productsy
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.Category_name)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (repo *ProductRepository) Create(product *models.Product) error {
	query := "INSERT INTO product (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id"
	err := repo.db.QueryRow(query, product.Name, product.Price, product.Stock, product.Category_id).Scan(&product.ID)
	if err != nil {
		return err
	}
	queryGetCategory := "SELECT name FROM categories WHERE id = $1"
	err = repo.db.QueryRow(queryGetCategory, product.Category_id).Scan(&product.Category_name)

	if err != nil {
		product.Category_name = ""
	}
	return nil
}

// GetByID - ambil produk by ID
func (repo *ProductRepository) GetByID(id int) (*models.Productsy, error) {
	query := "SELECT p.id, p.name, p.price, p.stock, c.name FROM product p JOIN category c ON p.category_id = c.id WHERE id = $1"

	var p models.Productsy
	err := repo.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.Category_name)
	if err == sql.ErrNoRows {
		return nil, errors.New("produk tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (repo *ProductRepository) Update(product *models.Product) error {
	query := "UPDATE product SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5"
	result, err := repo.db.Exec(query, product.Name, product.Price, product.Stock, product.Category_id, product.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("produk tidak ditemukan")
	}

	queryGetCategory := "SELECT name FROM categories WHERE id = $1"
	err = repo.db.QueryRow(queryGetCategory, product.Category_id).Scan(&product.Category_name)

	if err != nil {
		product.Category_name = ""
	}
	return nil
}

func (repo *ProductRepository) Delete(id int) error {
	query := "DELETE FROM product WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("produk tidak ditemukan")
	}

	return err
}
