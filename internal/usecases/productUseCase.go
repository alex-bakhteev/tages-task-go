package usecases

import (
	"context"
	"errors"
	"tages-task-go/internal/models/modelssvc"
	"tages-task-go/internal/models/modelsuc"
)

func (u *Usecase) CreateProduct(ctx context.Context, product modelsuc.ProductUC) error {
	productSrv := &modelssvc.ProductSrv{
		Name:  product.Name,
		Price: product.Price,
	}

	if err := u.ProductUC.CreateProduct(ctx, productSrv); err != nil {
		u.Logger.ErrorCtx(ctx, "Failed to create product: %v", err)
		return errors.New("failed to create product: " + err.Error())
	}
	u.Logger.InfoCtx(ctx, "Product created successfully: Name=%s", product.Name)
	return nil
}

func (u *Usecase) GetProduct(ctx context.Context, id int) (modelsuc.ProductUC, error) {
	productSrv, err := u.ProductUC.GetProductByID(ctx, id)
	if err != nil {
		u.Logger.ErrorCtx(ctx, "Failed to get product by ID=%d: %v", id, err)
		return modelsuc.ProductUC{}, errors.New("failed to get product: " + err.Error())
	}

	productUC := modelsuc.ProductUC{
		ID:    productSrv.ID,
		Name:  productSrv.Name,
		Price: productSrv.Price,
	}
	u.Logger.InfoCtx(ctx, "Product retrieved successfully by ID=%d", id)
	return productUC, nil
}

func (u *Usecase) GetAllProducts(ctx context.Context) ([]modelsuc.ProductUC, error) {
	u.Logger.InfoCtx(ctx, "Fetching all products in usecase")

	productsSrv, err := u.ProductUC.GetAllProducts(ctx)
	if err != nil {
		u.Logger.ErrorCtx(ctx, "Failed to get all products: %v", err)
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
	u.Logger.InfoCtx(ctx, "All products retrieved successfully, count=%d", len(productsUC))
	return productsUC, nil
}
