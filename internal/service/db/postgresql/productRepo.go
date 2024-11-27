package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"tages-task-go/internal/models/modelssvc"
	"tages-task-go/pkg/logging"
)

type productRepository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

func NewProductRepository(db *pgxpool.Pool, logger *logging.Logger) *productRepository {
	return &productRepository{db: db, logger: logger}
}

// Создание продукта
func (r *productRepository) CreateProduct(ctx context.Context, product *modelssvc.ProductSrv) error {
	r.logger.InfoCtx(ctx, "Creating product with Name=%s, Price=%.2f", product.Name, product.Price)
	_, err := r.db.Exec(ctx, "INSERT INTO products (name, price) VALUES ($1, $2)", product.Name, product.Price)
	if err != nil {
		r.logPgError(ctx, "Failed to create product", err)
		return fmt.Errorf("failed to create product: %w", err)
	}
	r.logger.InfoCtx(ctx, "Product created successfully with Name=%s", product.Name)
	return nil
}

// Получение продукта по ID
func (r *productRepository) GetProductByID(ctx context.Context, id int) (*modelssvc.ProductSrv, error) {
	r.logger.InfoCtx(ctx, "Fetching product by ID=%d", id)
	var product modelssvc.ProductSrv
	err := r.db.QueryRow(ctx, "SELECT id, name, price FROM products WHERE id=$1", id).
		Scan(&product.ID, &product.Name, &product.Price)
	if err != nil {
		r.logPgError(ctx, fmt.Sprintf("Failed to fetch product by ID=%d", id), err)
		return nil, fmt.Errorf("failed to fetch product: %w", err)
	}
	r.logger.InfoCtx(ctx, "Product fetched successfully by ID=%d", id)
	return &product, nil
}

// Получение всех продуктов
func (r *productRepository) GetAllProducts(ctx context.Context) ([]*modelssvc.ProductSrv, error) {
	r.logger.InfoCtx(ctx, "Fetching all products")
	rows, err := r.db.Query(ctx, "SELECT id, name, price FROM products")
	if err != nil {
		r.logPgError(ctx, "Failed to fetch all products", err)
		return nil, fmt.Errorf("failed to fetch all products: %w", err)
	}
	defer rows.Close()

	var products []*modelssvc.ProductSrv
	for rows.Next() {
		product := &modelssvc.ProductSrv{}
		err = rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			r.logPgError(ctx, "Failed to scan product row", err)
			return nil, fmt.Errorf("failed to scan product row: %w", err)
		}
		products = append(products, product)
	}
	r.logger.InfoCtx(ctx, "All products fetched successfully, count=%d", len(products))
	return products, nil
}

// Логирование ошибок PostgreSQL
func (r *productRepository) logPgError(ctx context.Context, msg string, err error) {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		r.logger.ErrorCtx(ctx, "%s: SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
			msg, pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
	} else {
		r.logger.ErrorCtx(ctx, "%s: %v", msg, err)
	}
}
