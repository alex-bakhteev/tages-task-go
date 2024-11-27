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
		o.logger.ErrorCtx(ctx, "Failed to create order: %v", err)
		return errors.New("failed to create order: " + err.Error())
	}
	o.logger.InfoCtx(ctx, "Order created successfully: ID=%d", order.ID)
	return nil
}

func (o *OrderUC) GetOrder(ctx context.Context, id int) (modelsuc.OrderUC, error) {
	orderSrv, err := o.repo.GetOrderByID(ctx, id)
	if err != nil {
		o.logger.ErrorCtx(ctx, "Failed to get order by ID=%d: %v", id, err)
		return modelsuc.OrderUC{}, errors.New("failed to get order: " + err.Error())
	}

	orderUC := modelsuc.OrderUC{
		ID:        orderSrv.ID,
		ProductID: orderSrv.ProductID,
		Quantity:  orderSrv.Quantity,
	}
	o.logger.InfoCtx(ctx, "Order retrieved successfully by ID=%d", id)
	return orderUC, nil
}

func (o *OrderUC) GetAllOrders(ctx context.Context) ([]modelsuc.OrderUC, error) {
	o.logger.InfoCtx(ctx, "Fetching all orders in usecase")

	ordersSrv, err := o.repo.GetAllOrders(ctx)
	if err != nil {
		o.logger.ErrorCtx(ctx, "Failed to get all orders: %v", err)
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
	o.logger.InfoCtx(ctx, "All orders retrieved successfully, count=%d", len(ordersUC))
	return ordersUC, nil
}
