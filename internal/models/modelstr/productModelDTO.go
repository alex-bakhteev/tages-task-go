package modelstr

type ProductDTO struct {
	ID    int     `json:"id"`
	Name  string  `json:"name" validate:"max=100"`
	Price float64 `json:"price"`
}
