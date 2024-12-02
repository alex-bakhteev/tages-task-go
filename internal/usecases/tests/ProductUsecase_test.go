package tests

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"tages-task-go/internal/models/modelsuc"
	"tages-task-go/internal/usecases/mocks"
	"testing"
)

func TestProductUsecase_CreateProduct(t *testing.T) {
	mockProductUsecase := mocks.NewProductUsecase(t)

	tests := []struct {
		name    string
		product modelsuc.ProductUC
		mockErr error
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid input",
			product: modelsuc.ProductUC{
				ID:    1,
				Name:  "Product A",
				Price: 100.0,
			},
			mockErr: nil,
			wantErr: false,
		},
		{
			name: "Invalid input - missing Name",
			product: modelsuc.ProductUC{
				ID:    2,
				Price: 50.0,
			},
			mockErr: errors.New("missing product name"),
			wantErr: true,
			errMsg:  "missing product name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProductUsecase.On("CreateProduct", mock.Anything, tt.product).Return(tt.mockErr).Once()

			err := mockProductUsecase.CreateProduct(context.Background(), tt.product)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
			mockProductUsecase.AssertExpectations(t)
		})
	}
}

func TestProductUsecase_GetAllProducts(t *testing.T) {
	mockProductUsecase := mocks.NewProductUsecase(t)

	tests := []struct {
		name       string
		mockResult []modelsuc.ProductUC
		mockErr    error
		wantErr    bool
		errMsg     string
	}{
		{
			name: "Valid response",
			mockResult: []modelsuc.ProductUC{
				{ID: 1, Name: "Product A", Price: 100.0},
				{ID: 2, Name: "Product B", Price: 200.0},
			},
			mockErr: nil,
			wantErr: false,
		},
		{
			name:       "Error response",
			mockResult: nil,
			mockErr:    errors.New("database error"),
			wantErr:    true,
			errMsg:     "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProductUsecase.On("GetAllProducts", mock.Anything).Return(tt.mockResult, tt.mockErr).Once()

			products, err := mockProductUsecase.GetAllProducts(context.Background())

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
				assert.Nil(t, products)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockResult, products)
			}
			mockProductUsecase.AssertExpectations(t)
		})
	}
}

func TestProductUsecase_GetProduct(t *testing.T) {
	mockProductUsecase := mocks.NewProductUsecase(t)

	tests := []struct {
		name       string
		id         int
		mockResult modelsuc.ProductUC
		mockErr    error
		wantErr    bool
		errMsg     string
	}{
		{
			name: "Valid input",
			id:   1,
			mockResult: modelsuc.ProductUC{
				ID:    1,
				Name:  "Product A",
				Price: 100.0,
			},
			mockErr: nil,
			wantErr: false,
		},
		{
			name:       "Invalid input - negative ID",
			id:         -1,
			mockResult: modelsuc.ProductUC{},
			mockErr:    errors.New("invalid product ID"),
			wantErr:    true,
			errMsg:     "invalid product ID",
		},
		{
			name:       "Error response - product not found",
			id:         999,
			mockResult: modelsuc.ProductUC{},
			mockErr:    errors.New("product not found"),
			wantErr:    true,
			errMsg:     "product not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProductUsecase.On("GetProduct", mock.Anything, tt.id).Return(tt.mockResult, tt.mockErr).Once()

			product, err := mockProductUsecase.GetProduct(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
				assert.Equal(t, modelsuc.ProductUC{}, product)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockResult, product)
			}
			mockProductUsecase.AssertExpectations(t)
		})
	}
}
