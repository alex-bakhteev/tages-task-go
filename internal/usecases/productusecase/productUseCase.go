package productusecase

import (
	"context"
	"errors"
	"tages-task-go/internal/models/modelssvc"
	"tages-task-go/internal/models/modelsuc"
	"tages-task-go/internal/usecases"
	"tages-task-go/pkg/logging"
)

type ProductUsecase struct {
	repo   usecases.ProductRepo
	logger *logging.Logger
}

func New(repo usecases.ProductRepo, logger *logging.Logger) *ProductUsecase {
	return &ProductUsecase{repo: repo, logger: logger}
}

func (p *ProductUsecase) CreateProduct(ctx context.Context, product modelsuc.ProductUC) error {
	productSrv := &modelssvc.ProductSrv{
		Name:  product.Name,
		Price: product.Price,
	}

	if err := p.repo.CreateProduct(ctx, productSrv); err != nil {
		p.logger.ErrorCtx(ctx, "Failed to create product: %v", err)
		return errors.New("failed to create product: " + err.Error())
	}
	p.logger.InfoCtx(ctx, "Product created successfully: Name=%s", product.Name)
	return nil
}

func (p *ProductUsecase) GetProduct(ctx context.Context, id int) (modelsuc.ProductUC, error) {
	productSrv, err := p.repo.GetProductByID(ctx, id)
	if err != nil {
		p.logger.ErrorCtx(ctx, "Failed to get product by ID=%d: %v", id, err)
		return modelsuc.ProductUC{}, errors.New("failed to get product: " + err.Error())
	}

	productUC := modelsuc.ProductUC{
		ID:    productSrv.ID,
		Name:  productSrv.Name,
		Price: productSrv.Price,
	}
	p.logger.InfoCtx(ctx, "Product retrieved successfully by ID=%d", id)
	return productUC, nil
}

func (p *ProductUsecase) GetAllProducts(ctx context.Context) ([]modelsuc.ProductUC, error) {
	p.logger.InfoCtx(ctx, "Fetching all products in usecase")

	productsSrv, err := p.repo.GetAllProducts(ctx)
	if err != nil {
		p.logger.ErrorCtx(ctx, "Failed to get all products: %v", err)
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
	p.logger.InfoCtx(ctx, "All products retrieved successfully, count=%d", len(productsUC))
	return productsUC, nil
}
