package usecase

type StoreUseCase interface {
	OrderUseCase
	ProductUseCase
}

type storeUseCase struct {
	OrderUseCase
	ProductUseCase
}

func NewStoreUseCase(orderUC OrderUseCase, productUC ProductUseCase) StoreUseCase {
	return &storeUseCase{
		OrderUseCase:   orderUC,
		ProductUseCase: productUC,
	}
}
