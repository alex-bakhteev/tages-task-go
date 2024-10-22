package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn" // Импортируем пакет для работы с ошибками
	"github.com/jackc/pgx/v5/pgxpool"
	"tages-task-go/internal/logger" // Импортируйте ваш пакет логгирования
	"tages-task-go/internal/model"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetAllProducts(ctx context.Context) ([]model.Product, error) {
	rows, err := r.db.Query(ctx, "SELECT id, name, price FROM products")
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			logger.Log.Error(newErr) // Логируем детализированную ошибку
			return nil, newErr
		}
		logger.Log.Println("Error querying all products:", err) // Логируем общую ошибку
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var product model.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok {
				newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
					pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
				logger.Log.Error(newErr) // Логируем детализированную ошибку
				return nil, newErr
			}
			logger.Log.Println("Error scanning product:", err) // Логируем общую ошибку
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (r *ProductRepository) CreateProduct(ctx context.Context, product model.Product) error {
	_, err := r.db.Exec(ctx, "INSERT INTO products (name, price) VALUES ($1, $2)", product.Name, product.Price)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			logger.Log.Error(newErr) // Логируем детализированную ошибку
			return newErr
		}
		logger.Log.Println("Error creating product:", err) // Логируем общую ошибку
	}
	return err
}

func (r *ProductRepository) GetProductByID(ctx context.Context, id int) (model.Product, error) {
	var product model.Product
	err := r.db.QueryRow(ctx, "SELECT id, name, price FROM products WHERE id=$1", id).Scan(&product.ID, &product.Name, &product.Price)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			logger.Log.Error(newErr) // Логируем детализированную ошибку
			return product, newErr
		}
		logger.Log.Println("Error fetching product by ID:", err) // Логируем общую ошибку
		return product, err
	}
	return product, nil
}
