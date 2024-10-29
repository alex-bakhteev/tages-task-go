package usecase

import (
	"context"
	"errors"
	"tages-task-go/internal/service/db/postgresql"
	"tages-task-go/pkg/logging"
	"tages-task-go/pkg/models/service"
	"tages-task-go/pkg/models/usecase"
)

type ProductUseCase interface {
	CreateProduct(ctx context.Context, product usecase.ProductUC) error
	GetProduct(ctx context.Context, id int) (usecase.ProductUC, error)
	GetAllProducts(ctx context.Context) ([]usecase.ProductUC, error)
}

type productUsecase struct {
	repo   postgresql.ProductRepository
	logger *logging.Logger
}

func NewProductUseCase(repo postgresql.ProductRepository, logger *logging.Logger) ProductUseCase {
	return &productUsecase{repo: repo, logger: logger}
}

func (p *productUsecase) CreateProduct(ctx context.Context, product usecase.ProductUC) error {
	productSrv := service.ProductSrv{
		Name:  product.Name,
		Price: product.Price,
	}

	if err := p.repo.CreateProduct(ctx, productSrv); err != nil {
		p.logger.Error("Failed to create product: ", err)
		return errors.New("failed to create product: " + err.Error())
	}
	p.logger.Info("Product created successfully")
	return nil
}

func (p *productUsecase) GetProduct(ctx context.Context, id int) (usecase.ProductUC, error) {
	productSrv, err := p.repo.GetProductByID(ctx, id)
	if err != nil {
		p.logger.Error("Failed to get product by ID: ", err)
		return usecase.ProductUC{}, errors.New("failed to get product: " + err.Error())
	}

	productUC := usecase.ProductUC{
		ID:    productSrv.ID,
		Name:  productSrv.Name,
		Price: productSrv.Price,
	}
	p.logger.Info("Product retrieved successfully by ID:", id)
	return productUC, nil
}

func (p *productUsecase) GetAllProducts(ctx context.Context) ([]usecase.ProductUC, error) {
	productsSrv, err := p.repo.GetAllProducts(ctx)
	if err != nil {
		p.logger.Error("Failed to get all products: ", err)
		return nil, errors.New("failed to get products: " + err.Error())
	}

	var productsUC []usecase.ProductUC
	for _, productSrv := range productsSrv {
		productUC := usecase.ProductUC{
			ID:    productSrv.ID,
			Name:  productSrv.Name,
			Price: productSrv.Price,
		}
		productsUC = append(productsUC, productUC)
	}
	p.logger.Info("All products retrieved successfully")
	return productsUC, nil
}
