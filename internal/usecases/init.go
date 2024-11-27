package usecases

import (
	"context"
	"tages-task-go/internal/models/modelssvc"
)

type OrderRepo interface {
	GetAllOrders(ctx context.Context) ([]*modelssvc.OrderSrv, error)
	GetOrderByID(ctx context.Context, id int) (*modelssvc.OrderSrv, error)
	CreateOrder(ctx context.Context, order *modelssvc.OrderSrv) error
}

type ProductRepo interface {
	CreateProduct(ctx context.Context, product *modelssvc.ProductSrv) error
	GetProductByID(ctx context.Context, id int) (*modelssvc.ProductSrv, error)
	GetAllProducts(ctx context.Context) ([]*modelssvc.ProductSrv, error)
}

type Usecase struct {
	OrderRepo   OrderRepo
	ProductRepo ProductRepo
}
