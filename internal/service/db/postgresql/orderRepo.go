package postgresql

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"tages-task-go/pkg/logging"
	"tages-task-go/pkg/models/service"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *service.OrderSrv) error
	GetOrderByID(ctx context.Context, id int) (*service.OrderSrv, error)
	GetAllOrders(ctx context.Context) ([]*service.OrderSrv, error) // Новый метод для всех заказов
}

type orderRepository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

func NewOrderRepository(db *pgxpool.Pool, logger *logging.Logger) *orderRepository {
	return &orderRepository{db: db, logger: logger}
}

// Получение всех заказов
func (r *orderRepository) GetAllOrders(ctx context.Context) ([]*service.OrderSrv, error) {
	rows, err := r.db.Query(ctx, "SELECT id, product_id, quantity, total_price FROM orders")
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			r.logger.Error(newErr)
			return nil, newErr
		}
		r.logger.Println("Error querying orders:", err)
		return nil, err
	}
	defer rows.Close()

	var orders []*service.OrderSrv
	for rows.Next() {
		// Инициализируем переменную order перед каждой итерацией
		order := &service.OrderSrv{}
		err = rows.Scan(&order.ID, &order.ProductID, &order.Quantity, &order.TotalPrice)
		if err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok {
				newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
					pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
				r.logger.Error(newErr)
				return nil, newErr
			}
			r.logger.Println("Error scanning order:", err)
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// Получение заказа по ID
func (r *orderRepository) GetOrderByID(ctx context.Context, id int) (*service.OrderSrv, error) {
	var order service.OrderSrv
	err := r.db.QueryRow(ctx, "SELECT id, product_id, quantity, total_price FROM orders WHERE id=$1", id).
		Scan(&order.ID, &order.ProductID, &order.Quantity, &order.TotalPrice)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			r.logger.Error(newErr) // Логируем детализированную ошибку
			return nil, newErr
		}
		r.logger.Println("Error fetching order by ID:", err) // Логируем общую ошибку
		return nil, err
	}
	return &order, nil
}

// Создание нового заказа с автоматическим расчетом total_price
func (r *orderRepository) CreateOrder(ctx context.Context, order *service.OrderSrv) error {
	// Получаем цену товара
	var productPrice float64
	err := r.db.QueryRow(ctx, "SELECT price FROM products WHERE id=$1", order.ProductID).Scan(&productPrice)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			r.logger.Error(newErr) // Логируем детализированную ошибку
			return newErr
		}
		r.logger.Println("Error fetching product price for order:", err) // Логируем общую ошибку
		return errors.New("product not found")
	}

	// Рассчитываем общую стоимость заказа
	totalPrice := productPrice * float64(order.Quantity)

	// Вставляем новый заказ
	_, err = r.db.Exec(ctx, "INSERT INTO orders (product_id, quantity, total_price) VALUES ($1, $2, $3)",
		order.ProductID, order.Quantity, totalPrice)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			r.logger.Error(newErr) // Логируем детализированную ошибку
			return newErr
		}
		r.logger.Println("Error creating order:", err) // Логируем общую ошибку
	}
	return err
}
