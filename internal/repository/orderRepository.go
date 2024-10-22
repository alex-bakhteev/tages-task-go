package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn" // Импортируем пакет для работы с ошибками
	"github.com/jackc/pgx/v5/pgxpool"
	"tages-task-go/internal/logger" // Импортируйте ваш пакет логгирования
	"tages-task-go/internal/model"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

// Получение всех заказов
func (r *OrderRepository) GetOrders(ctx context.Context) ([]model.Order, error) {
	rows, err := r.db.Query(ctx, "SELECT id, product_id, quantity, total_price FROM orders")
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			logger.Log.Error(newErr) // Логируем детализированную ошибку
			return nil, newErr
		}
		logger.Log.Println("Error querying orders:", err) // Логируем общую ошибку
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var order model.Order
		err = rows.Scan(&order.ID, &order.ProductID, &order.Quantity, &order.TotalPrice)
		if err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok {
				newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
					pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
				logger.Log.Error(newErr) // Логируем детализированную ошибку
				return nil, newErr
			}
			logger.Log.Println("Error scanning order:", err) // Логируем общую ошибку
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// Получение заказа по ID
func (r *OrderRepository) GetOrderByID(ctx context.Context, id int) (*model.Order, error) {
	var order model.Order
	err := r.db.QueryRow(ctx, "SELECT id, product_id, quantity, total_price FROM orders WHERE id=$1", id).
		Scan(&order.ID, &order.ProductID, &order.Quantity, &order.TotalPrice)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			logger.Log.Error(newErr) // Логируем детализированную ошибку
			return nil, newErr
		}
		logger.Log.Println("Error fetching order by ID:", err) // Логируем общую ошибку
		return nil, err
	}
	return &order, nil
}

// Создание нового заказа с автоматическим расчетом total_price
func (r *OrderRepository) CreateOrder(ctx context.Context, order model.Order) error {
	// Получаем цену товара
	var productPrice float64
	err := r.db.QueryRow(ctx, "SELECT price FROM products WHERE id=$1", order.ProductID).Scan(&productPrice)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			logger.Log.Error(newErr) // Логируем детализированную ошибку
			return newErr
		}
		logger.Log.Println("Error fetching product price for order:", err) // Логируем общую ошибку
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
			logger.Log.Error(newErr) // Логируем детализированную ошибку
			return newErr
		}
		logger.Log.Println("Error creating order:", err) // Логируем общую ошибку
	}
	return err
}
