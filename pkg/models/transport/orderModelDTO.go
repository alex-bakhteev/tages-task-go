package transport

type OrderDTO struct {
	ID        int `json:"id"`
	ProductID int `json:"productId"`
	Quantity  int `json:"quantity"`
}
