package orderusecase

import (
	"context"
	"errors"
	"tages-task-go/internal/models/modelssvc"
	"tages-task-go/internal/models/modelsuc"
	"tages-task-go/pkg/logging"
)

type OrderRepo interface {
	GetAllOrders(ctx context.Context) ([]*modelssvc.OrderSrv, error)
	GetOrderByID(ctx context.Context, id int) (*modelssvc.OrderSrv, error)
	CreateOrder(ctx context.Context, order *modelssvc.OrderSrv) error
}

type OrderUC struct {
	Repo   OrderRepo
	Logger *logging.Logger
}

func (o *OrderUC) CreateOrder(ctx context.Context, order modelsuc.OrderUC) error {
	orderSrv := modelssvc.OrderSrv{
		ID:        order.ID,
		ProductID: order.ProductID,
		Quantity:  order.Quantity,
	}

	if err := o.Repo.CreateOrder(ctx, &orderSrv); err != nil {
		o.Logger.ErrorCtx(ctx, "Failed to create order: %v", err)
		return errors.New("failed to create order: " + err.Error())
	}
	o.Logger.InfoCtx(ctx, "Order created successfully: ID=%d", order.ID)
	return nil
}

func (o *OrderUC) GetOrder(ctx context.Context, id int) (modelsuc.OrderUC, error) {
	orderSrv, err := o.Repo.GetOrderByID(ctx, id)
	if err != nil {
		o.Logger.ErrorCtx(ctx, "Failed to get order by ID=%d: %v", id, err)
		return modelsuc.OrderUC{}, errors.New("failed to get order: " + err.Error())
	}

	orderUC := modelsuc.OrderUC{
		ID:        orderSrv.ID,
		ProductID: orderSrv.ProductID,
		Quantity:  orderSrv.Quantity,
	}
	o.Logger.InfoCtx(ctx, "Order retrieved successfully by ID=%d", id)
	return orderUC, nil
}

func (o *OrderUC) GetAllOrders(ctx context.Context) ([]modelsuc.OrderUC, error) {
	o.Logger.InfoCtx(ctx, "Fetching all orders in usecase")

	ordersSrv, err := o.Repo.GetAllOrders(ctx)
	if err != nil {
		o.Logger.ErrorCtx(ctx, "Failed to get all orders: %v", err)
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
	o.Logger.InfoCtx(ctx, "All orders retrieved successfully, count=%d", len(ordersUC))
	return ordersUC, nil
}
