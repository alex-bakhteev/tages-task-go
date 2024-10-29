package models

import (
	modelsSrv "tages-task-go/pkg/models/service"
	modelsDTO "tages-task-go/pkg/models/transport"
	modelsUC "tages-task-go/pkg/models/usecase"
)

// FromDtoToUsecase - преобразует транспортную модель OrderDTO в модель usecase.OrderUC
func FromDtoToUseCaseOrder(orderDTO modelsDTO.OrderDTO) modelsUC.OrderUC {
	return modelsUC.OrderUC{
		ID:        orderDTO.ID,
		ProductID: orderDTO.ProductID,
		Quantity:  orderDTO.Quantity,
	}
}

// FromUsecaseToDto - преобразует модель usecase.OrderUC обратно в транспортную модель OrderDTO для ответа клиенту
func FromUseCaseToDtoOrder(orderUC modelsUC.OrderUC) modelsDTO.OrderDTO {
	return modelsDTO.OrderDTO{
		ID:        orderUC.ID,
		ProductID: orderUC.ProductID,
		Quantity:  orderUC.Quantity,
	}
}

// MapToUsecaseProduct - преобразует транспортную модель ProductDTO в usecase.ProductUC
func FromDtoToUseCaseProduct(productDTO modelsDTO.ProductDTO) modelsUC.ProductUC {
	return modelsUC.ProductUC{
		ID:    productDTO.ID,
		Name:  productDTO.Name,
		Price: productDTO.Price,
	}
}

// MapToTransportProduct - преобразует модель usecase.ProductUC в транспортную модель ProductDTO
func FromUseCaseToDtoProduct(productUC modelsUC.ProductUC) modelsDTO.ProductDTO {
	return modelsDTO.ProductDTO{
		ID:    productUC.ID,
		Name:  productUC.Name,
		Price: productUC.Price,
	}
}

// FromDtoToUsecase - преобразует транспортную модель OrderDTO в модель usecase.OrderUC
func FromServiceToUseCaseOrder(orderSrv modelsSrv.OrderSrv) modelsUC.OrderUC {
	return modelsUC.OrderUC{
		ID:        orderSrv.ID,
		ProductID: orderSrv.ProductID,
		Quantity:  orderSrv.Quantity,
	}
}

// MapToUsecaseProduct - преобразует транспортную модель ProductDTO в usecase.ProductUC
func FromServiceToUseCaseProduct(productSrv modelsSrv.ProductSrv) modelsUC.ProductUC {
	return modelsUC.ProductUC{
		ID:    productSrv.ID,
		Name:  productSrv.Name,
		Price: productSrv.Price,
	}
}

// FromUsecaseToDto - преобразует модель usecase.OrderUC обратно в транспортную модель OrderDTO для ответа клиенту
func FromUseCaseToServiceOrder(orderUC modelsUC.OrderUC) modelsSrv.OrderSrv {
	return modelsSrv.OrderSrv{
		ID:         orderUC.ID,
		ProductID:  orderUC.ProductID,
		Quantity:   orderUC.Quantity,
		TotalPrice: 0,
	}
}

// MapToTransportProduct - преобразует модель usecase.ProductUC в транспортную модель ProductDTO
func FromUseCaseToServiceProduct(productUC modelsUC.ProductUC) modelsSrv.ProductSrv {
	return modelsSrv.ProductSrv{
		ID:    productUC.ID,
		Name:  productUC.Name,
		Price: productUC.Price,
	}
}
