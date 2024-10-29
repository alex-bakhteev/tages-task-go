package usecase

import (
	"context"
	"errors"
	"tages-task-go/internal/service/db/postgresql"
	"tages-task-go/pkg/logging"
	"tages-task-go/pkg/models"
	"tages-task-go/pkg/models/usecase"
)

type OrderUseCase interface {
	CreateOrder(ctx context.Context, order usecase.OrderUC) error
	GetOrder(ctx context.Context, id int) (usecase.OrderUC, error)
	GetAllOrders(ctx context.Context) ([]usecase.OrderUC, error)
}

type orderUC struct {
	repo   postgresql.OrderRepository
	logger *logging.Logger
}

func NewOrderUseCase(repo postgresql.OrderRepository, logger *logging.Logger) OrderUseCase {
	return &orderUC{repo: repo, logger: logger}
}

func (o *orderUC) CreateOrder(ctx context.Context, order usecase.OrderUC) error {
	orderSrv := models.FromUseCaseToServiceOrder(order)

	if err := o.repo.CreateOrder(ctx, &orderSrv); err != nil {
		o.logger.Error("Failed to create order: ", err)
		return errors.New("failed to create order: " + err.Error())
	}
	o.logger.Info("Order created successfully")
	return nil
}

func (o *orderUC) GetOrder(ctx context.Context, id int) (usecase.OrderUC, error) {
	orderSrv, err := o.repo.GetOrderByID(ctx, id)
	if err != nil {
		o.logger.Error("Failed to get order by ID: ", err)
		return usecase.OrderUC{}, errors.New("failed to get order: " + err.Error())
	}

	orderUC := usecase.OrderUC{
		ID:        orderSrv.ID,
		ProductID: orderSrv.ProductID,
		Quantity:  orderSrv.Quantity,
	}
	o.logger.Info("Order retrieved successfully by ID:", id)
	return orderUC, nil
}

func (o *orderUC) GetAllOrders(ctx context.Context) ([]usecase.OrderUC, error) {
	ordersSrv, err := o.repo.GetAllOrders(ctx)
	if err != nil {
		o.logger.Error("Failed to get all orders: ", err)
		return nil, errors.New("failed to get orders: " + err.Error())
	}

	var ordersUC []usecase.OrderUC
	for _, orderSrv := range ordersSrv {
		orderUC := usecase.OrderUC{
			ID:        orderSrv.ID,
			ProductID: orderSrv.ProductID,
			Quantity:  orderSrv.Quantity,
		}
		ordersUC = append(ordersUC, orderUC)
	}
	o.logger.Info("All orders retrieved successfully")
	return ordersUC, nil
}
