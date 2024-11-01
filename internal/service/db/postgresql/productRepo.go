package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"tages-task-go/pkg/logging"
	"tages-task-go/pkg/models/service"
)

//type ProductRepository interface {
//	CreateProduct(ctx context.Context, product service.ProductSrv) error
//	GetProductByID(ctx context.Context, id int) (service.ProductSrv, error)
//	GetAllProducts(ctx context.Context) ([]service.ProductSrv, error)
//}

type productRepository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

func NewProductRepository(db *pgxpool.Pool, logger *logging.Logger) *productRepository {
	return &productRepository{db: db, logger: logger}
}

// Создание нового продукта
func (r *productRepository) CreateProduct(ctx context.Context, product service.ProductSrv) error {
	query := `INSERT INTO products (name, price) VALUES ($1, $2)`
	_, err := r.db.Exec(ctx, query, product.Name, product.Price)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			r.logger.Error(newErr)
			return newErr
		}
		r.logger.Println("Error creating product:", err)
	}
	return err
}

// Получение продукта по ID
func (r *productRepository) GetProductByID(ctx context.Context, id int) (service.ProductSrv, error) {
	var product service.ProductSrv
	query := `SELECT id, name, price FROM products WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).Scan(&product.ID, &product.Name, &product.Price)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			r.logger.Error(newErr)
			return product, newErr
		}
		r.logger.Println("Error fetching product by ID:", err)
	}
	return product, err
}

// Получение всех продуктов
func (r *productRepository) GetAllProducts(ctx context.Context) ([]service.ProductSrv, error) {
	query := `SELECT id, name, price FROM products`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			r.logger.Error(newErr)
			return nil, newErr
		}
		r.logger.Println("Error querying products:", err)
		return nil, err
	}
	defer rows.Close()

	var products []service.ProductSrv
	for rows.Next() {
		var product service.ProductSrv
		err = rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok {
				newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
					pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
				r.logger.Error(newErr)
				return nil, newErr
			}
			r.logger.Println("Error scanning product:", err)
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}
