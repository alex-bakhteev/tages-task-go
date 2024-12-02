package usecases

import (
	"tages-task-go/internal/usecases/orderusecase"
	"tages-task-go/internal/usecases/productusecase"
	"tages-task-go/pkg/logging"
)

type Usecase struct {
	ProductUC *productusecase.ProductUC
	OrderUC   *orderusecase.OrderUC
}

func NewUsecase(
	productRepo productusecase.ProductRepo,
	orderRepo orderusecase.OrderRepo,
	logger *logging.Logger,
) *Usecase {
	return &Usecase{
		ProductUC: &productusecase.ProductUC{
			Repo:   productRepo,
			Logger: logger,
		},
		OrderUC: &orderusecase.OrderUC{
			Repo:   orderRepo,
			Logger: logger,
		},
	}
}
