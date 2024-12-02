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

func TestOrderUsecase_CreateOrder(t *testing.T) {
	mockOrderUsecase := mocks.NewOrderUsecase(t)

	tests := []struct {
		name          string
		input         modelsuc.OrderUC
		mockError     error
		expectedError string
	}{
		{
			name: "Valid input",
			input: modelsuc.OrderUC{
				ID:        1,
				ProductID: 1001,
				Quantity:  2,
			},
			mockError:     nil,
			expectedError: "",
		},
		{
			name: "Invalid input - missing ProductID",
			input: modelsuc.OrderUC{
				ID:       2,
				Quantity: 3,
			},
			mockError:     errors.New("invalid ProductID"),
			expectedError: "invalid ProductID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockOrderUsecase.On("CreateOrder", mock.Anything, tt.input).
				Return(tt.mockError).
				Once()

			err := mockOrderUsecase.CreateOrder(context.Background(), tt.input)

			if tt.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
			}
			mockOrderUsecase.AssertExpectations(t)
		})
	}
}

func TestOrderUsecase_GetOrder(t *testing.T) {
	mockOrderUsecase := mocks.NewOrderUsecase(t)

	tests := []struct {
		name          string
		orderID       int
		mockOrder     modelsuc.OrderUC
		mockError     error
		expectedOrder modelsuc.OrderUC
		expectedError string
	}{
		{
			name:    "Valid input",
			orderID: 1,
			mockOrder: modelsuc.OrderUC{
				ID:        1,
				ProductID: 1001,
				Quantity:  2,
			},
			mockError:     nil,
			expectedOrder: modelsuc.OrderUC{ID: 1, ProductID: 1001, Quantity: 2},
			expectedError: "",
		},
		{
			name:          "Invalid input - negative ID",
			orderID:       -1,
			mockOrder:     modelsuc.OrderUC{},
			mockError:     errors.New("invalid order ID"),
			expectedOrder: modelsuc.OrderUC{},
			expectedError: "invalid order ID",
		},
		{
			name:          "Error response - order not found",
			orderID:       999,
			mockOrder:     modelsuc.OrderUC{},
			mockError:     errors.New("order not found"),
			expectedOrder: modelsuc.OrderUC{},
			expectedError: "order not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockOrderUsecase.On("GetOrder", mock.Anything, tt.orderID).
				Return(tt.mockOrder, tt.mockError).
				Once()

			order, err := mockOrderUsecase.GetOrder(context.Background(), tt.orderID)

			if tt.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOrder, order)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Equal(t, tt.expectedOrder, order)
			}
			mockOrderUsecase.AssertExpectations(t)
		})
	}
}
