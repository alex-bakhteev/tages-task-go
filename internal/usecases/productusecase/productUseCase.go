package productusecase

import (
	"context"
	"errors"
	"tages-task-go/internal/models/modelssvc"
	"tages-task-go/internal/models/modelsuc"
	"tages-task-go/pkg/logging"
)

type ProductRepo interface {
	CreateProduct(ctx context.Context, product modelssvc.ProductSrv) error
	GetProductByID(ctx context.Context, id int) (modelssvc.ProductSrv, error)
	GetAllProducts(ctx context.Context) ([]modelssvc.ProductSrv, error)
}

type ProductUsecase struct {
	repo   ProductRepo
	logger *logging.Logger
}

func New(repo ProductRepo, logger *logging.Logger) *ProductUsecase {
	return &ProductUsecase{repo: repo, logger: logger}
}

func (p *ProductUsecase) CreateProduct(ctx context.Context, product modelsuc.ProductUC) error {
	productSrv := modelssvc.ProductSrv{
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

func (p *ProductUsecase) GetProduct(ctx context.Context, id int) (modelsuc.ProductUC, error) {
	productSrv, err := p.repo.GetProductByID(ctx, id)
	if err != nil {
		p.logger.Error("Failed to get product by ID: ", err)
		return modelsuc.ProductUC{}, errors.New("failed to get product: " + err.Error())
	}

	productUC := modelsuc.ProductUC{
		ID:    productSrv.ID,
		Name:  productSrv.Name,
		Price: productSrv.Price,
	}
	p.logger.Info("Product retrieved successfully by ID:", id)
	return productUC, nil
}

func (p *ProductUsecase) GetAllProducts(ctx context.Context) ([]modelsuc.ProductUC, error) {
	productsSrv, err := p.repo.GetAllProducts(ctx)
	if err != nil {
		p.logger.Error("Failed to get all products: ", err)
		return nil, errors.New("failed to get products: " + err.Error())
	}

	var productsUC []modelsuc.ProductUC
	for _, productSrv := range productsSrv {
		productUC := modelsuc.ProductUC{
			ID:    productSrv.ID,
			Name:  productSrv.Name,
			Price: productSrv.Price,
		}
		productsUC = append(productsUC, productUC)
	}
	p.logger.Info("All products retrieved successfully")
	return productsUC, nil
}
