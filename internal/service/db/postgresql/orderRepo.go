package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"tages-task-go/internal/models/modelssvc"
	"tages-task-go/pkg/logging"
)

type orderRepository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

func NewOrderRepository(db *pgxpool.Pool, logger *logging.Logger) *orderRepository {
	return &orderRepository{db: db, logger: logger}
}

// Получение всех заказов
func (r *orderRepository) GetAllOrders(ctx context.Context) ([]*modelssvc.OrderSrv, error) {
	r.logger.InfoCtx(ctx, "Fetching all orders")
	rows, err := r.db.Query(ctx, "SELECT id, product_id, quantity, total_price FROM orders")
	if err != nil {
		r.logPgError(ctx, "Failed to fetch all orders", err)
		return nil, fmt.Errorf("failed to fetch all orders: %w", err)
	}
	defer rows.Close()

	var orders []*modelssvc.OrderSrv
	for rows.Next() {
		order := &modelssvc.OrderSrv{}
		err = rows.Scan(&order.ID, &order.ProductID, &order.Quantity, &order.TotalPrice)
		if err != nil {
			r.logPgError(ctx, "Failed to scan order row", err)
			return nil, fmt.Errorf("failed to scan order row: %w", err)
		}
		orders = append(orders, order)
	}
	r.logger.InfoCtx(ctx, "All orders fetched successfully, count=%d", len(orders))
	return orders, nil
}

// Получение заказа по ID
func (r *orderRepository) GetOrderByID(ctx context.Context, id int) (*modelssvc.OrderSrv, error) {
	r.logger.InfoCtx(ctx, "Fetching order by ID=%d", id)
	var order modelssvc.OrderSrv
	err := r.db.QueryRow(ctx, "SELECT id, product_id, quantity, total_price FROM orders WHERE id=$1", id).
		Scan(&order.ID, &order.ProductID, &order.Quantity, &order.TotalPrice)
	if err != nil {
		r.logPgError(ctx, fmt.Sprintf("Failed to fetch order by ID=%d", id), err)
		return nil, fmt.Errorf("failed to fetch order: %w", err)
	}
	r.logger.InfoCtx(ctx, "Order fetched successfully by ID=%d", id)
	return &order, nil
}

// Создание нового заказа с автоматическим расчетом total_price
func (r *orderRepository) CreateOrder(ctx context.Context, order *modelssvc.OrderSrv) error {
	r.logger.InfoCtx(ctx, "Creating order for ProductID=%d with Quantity=%d", order.ProductID, order.Quantity)
	// Получаем цену товара
	var productPrice float64
	err := r.db.QueryRow(ctx, "SELECT price FROM products WHERE id=$1", order.ProductID).Scan(&productPrice)
	if err != nil {
		r.logPgError(ctx, "Failed to fetch product price", err)
		return fmt.Errorf("failed to fetch product price: %w", err)
	}

	// Рассчитываем общую стоимость заказа
	totalPrice := productPrice * float64(order.Quantity)

	// Вставляем новый заказ
	_, err = r.db.Exec(ctx, "INSERT INTO orders (product_id, quantity, total_price) VALUES ($1, $2, $3)",
		order.ProductID, order.Quantity, totalPrice)
	if err != nil {
		r.logPgError(ctx, "Failed to create order", err)
		return fmt.Errorf("failed to create order: %w", err)
	}

	r.logger.InfoCtx(ctx, "Order created successfully for ProductID=%d, Quantity=%d, TotalPrice=%.2f",
		order.ProductID, order.Quantity, totalPrice)
	return nil
}

// Логирование ошибок PostgreSQL
func (r *orderRepository) logPgError(ctx context.Context, msg string, err error) {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		r.logger.ErrorCtx(ctx, "%s: SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
			msg, pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
	} else {
		r.logger.ErrorCtx(ctx, "%s: %v", msg, err)
	}
}
