package orderusecase

import (
	"context"
	"errors"
	"tages-task-go/internal/models/modelssvc"
	"tages-task-go/internal/models/modelsuc"
	"tages-task-go/internal/usecases"
	"tages-task-go/pkg/logging"
)

type OrderUC struct {
	repo   usecases.OrderRepo
	logger *logging.Logger
}

func New(repo usecases.OrderRepo, logger *logging.Logger) *OrderUC {
	return &OrderUC{repo: repo, logger: logger}
}

func (o *OrderUC) CreateOrder(ctx context.Context, order modelsuc.OrderUC) error {
	orderSrv := modelssvc.OrderSrv{
		ID:        order.ID,
		ProductID: order.ProductID,
		Quantity:  order.Quantity,
	}

	if err := o.repo.CreateOrder(ctx, &orderSrv); err != nil {
		o.logger.Error("Failed to create order: ", err)
		return errors.New("failed to create order: " + err.Error())
	}
	o.logger.Info("Order created successfully")
	return nil
}

func (o *OrderUC) GetOrder(ctx context.Context, id int) (modelsuc.OrderUC, error) {
	orderSrv, err := o.repo.GetOrderByID(ctx, id)
	if err != nil {
		o.logger.Error("Failed to get order by ID: ", err)
		return modelsuc.OrderUC{}, errors.New("failed to get order: " + err.Error())
	}

	orderUC := modelsuc.OrderUC{
		ID:        orderSrv.ID,
		ProductID: orderSrv.ProductID,
		Quantity:  orderSrv.Quantity,
	}
	o.logger.Info("Order retrieved successfully by ID:", id)
	return orderUC, nil
}

func (o *OrderUC) GetAllOrders(ctx context.Context) ([]modelsuc.OrderUC, error) {
	ordersSrv, err := o.repo.GetAllOrders(ctx)
	if err != nil {
		o.logger.Error("Failed to get all orders: ", err)
		return nil, errors.New("failed to get orders: " + err.Error())
	}

	var ordersUC []modelsuc.OrderUC
	for _, orderSrv := range ordersSrv {
		orderUC := modelsuc.OrderUC{
			ID:        orderSrv.ID,
			ProductID: orderSrv.ProductID,
			Quantity:  orderSrv.Quantity,
		}
		ordersUC = append(ordersUC, orderUC)
	}
	o.logger.Info("All orders retrieved successfully")
	return ordersUC, nil
}
