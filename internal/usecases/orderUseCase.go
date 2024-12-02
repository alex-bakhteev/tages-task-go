package usecases

import (
	"context"
	"errors"
	"tages-task-go/internal/models/modelssvc"
	"tages-task-go/internal/models/modelsuc"
)

func (u *Usecase) CreateOrder(ctx context.Context, order modelsuc.OrderUC) error {
	orderSrv := modelssvc.OrderSrv{
		ID:        order.ID,
		ProductID: order.ProductID,
		Quantity:  order.Quantity,
	}

	if err := u.OrderUC.CreateOrder(ctx, &orderSrv); err != nil {
		u.Logger.ErrorCtx(ctx, "Failed to create order: %v", err)
		return errors.New("failed to create order: " + err.Error())
	}
	u.Logger.InfoCtx(ctx, "Order created successfully: ID=%d", order.ID)
	return nil
}

func (u *Usecase) GetOrder(ctx context.Context, id int) (modelsuc.OrderUC, error) {
	orderSrv, err := u.OrderUC.GetOrderByID(ctx, id)
	if err != nil {
		u.Logger.ErrorCtx(ctx, "Failed to get order by ID=%d: %v", id, err)
		return modelsuc.OrderUC{}, errors.New("failed to get order: " + err.Error())
	}

	orderUC := modelsuc.OrderUC{
		ID:        orderSrv.ID,
		ProductID: orderSrv.ProductID,
		Quantity:  orderSrv.Quantity,
	}
	u.Logger.InfoCtx(ctx, "Order retrieved successfully by ID=%d", id)
	return orderUC, nil
}

func (u *Usecase) GetAllOrders(ctx context.Context) ([]modelsuc.OrderUC, error) {
	u.Logger.InfoCtx(ctx, "Fetching all orders in usecase")

	ordersSrv, err := u.OrderUC.GetAllOrders(ctx)
	if err != nil {
		u.Logger.ErrorCtx(ctx, "Failed to get all orders: %v", err)
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
	u.Logger.InfoCtx(ctx, "All orders retrieved successfully, count=%d", len(ordersUC))
	return ordersUC, nil
}
