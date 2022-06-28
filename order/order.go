package order

import "time"

const (
	StatusUnassigned string = "UNASSIGNED"
	StatusTaken      string = "TAKEN"
)

// Order struct to represents an Order
type Order struct {
	ID        int64     `json:"id" db:"id"`
	Distance  int       `json:"distance" db:"distance"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

// OrderUsecase represents Order Usecase
type OrderUsecase interface {
	PlaceOrder([]string, []string) (*Order, error)
	TakeOrder(int64) (string, error)
	ListOrders(int, int) (*[]Order, error)
}

// OrderRepository represents Order Repository
type OrderRepository interface {
	Create(*Order) error
	UpdateStatusByID(int64) error
	FindByID(int64) (*Order, error)
	FindRange(int, int) (*[]Order, error)
}
