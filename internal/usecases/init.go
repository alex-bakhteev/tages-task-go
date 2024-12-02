package usecases

import (
	"context"
	"tages-task-go/internal/models/modelssvc"
	"tages-task-go/pkg/logging"
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
	ProductUC ProductRepo
	OrderUC   OrderRepo
	Logger    *logging.Logger
}

func NewUsecase(
	productRepo ProductRepo,
	orderRepo OrderRepo,
	logger *logging.Logger,
) *Usecase {
	return &Usecase{
		ProductUC: productRepo,
		OrderUC:   orderRepo,
		Logger:    logger,
	}
}
