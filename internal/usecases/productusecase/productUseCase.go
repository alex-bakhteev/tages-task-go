package productusecase

import (
	"context"
	"errors"
	"tages-task-go/internal/models/modelssvc"
	"tages-task-go/internal/models/modelsuc"
	"tages-task-go/pkg/logging"
)

type ProductRepo interface {
	CreateProduct(ctx context.Context, product *modelssvc.ProductSrv) error
	GetProductByID(ctx context.Context, id int) (*modelssvc.ProductSrv, error)
	GetAllProducts(ctx context.Context) ([]*modelssvc.ProductSrv, error)
}

type ProductUC struct {
	Repo   ProductRepo
	Logger *logging.Logger
}

func (p *ProductUC) CreateProduct(ctx context.Context, product modelsuc.ProductUC) error {
	productSrv := &modelssvc.ProductSrv{
		Name:  product.Name,
		Price: product.Price,
	}

	if err := p.Repo.CreateProduct(ctx, productSrv); err != nil {
		p.Logger.ErrorCtx(ctx, "Failed to create product: %v", err)
		return errors.New("failed to create product: " + err.Error())
	}
	p.Logger.InfoCtx(ctx, "Product created successfully: Name=%s", product.Name)
	return nil
}

func (p *ProductUC) GetProduct(ctx context.Context, id int) (modelsuc.ProductUC, error) {
	productSrv, err := p.Repo.GetProductByID(ctx, id)
	if err != nil {
		p.Logger.ErrorCtx(ctx, "Failed to get product by ID=%d: %v", id, err)
		return modelsuc.ProductUC{}, errors.New("failed to get product: " + err.Error())
	}

	productUC := modelsuc.ProductUC{
		ID:    productSrv.ID,
		Name:  productSrv.Name,
		Price: productSrv.Price,
	}
	p.Logger.InfoCtx(ctx, "Product retrieved successfully by ID=%d", id)
	return productUC, nil
}

func (p *ProductUC) GetAllProducts(ctx context.Context) ([]modelsuc.ProductUC, error) {
	p.Logger.InfoCtx(ctx, "Fetching all products in usecase")

	productsSrv, err := p.Repo.GetAllProducts(ctx)
	if err != nil {
		p.Logger.ErrorCtx(ctx, "Failed to get all products: %v", err)
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
	p.Logger.InfoCtx(ctx, "All products retrieved successfully, count=%d", len(productsUC))
	return productsUC, nil
}
