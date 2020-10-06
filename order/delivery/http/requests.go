package http

// PlaceOrderRequest represents the object of place order request params
type PlaceOrderRequest struct {
	Origin      []string `json:"origin"`
	Destination []string `json:"destination"`
}

// TakeOrderRequest represents the object of take order request params
type TakeOrderRequest struct {
	ID     int64  `uri:"id" valid:"int"`
	Status string `json:"status" valid:"-"`
}

// ListOrderRequest represents the object of list order request params
type ListOrderRequest struct {
	Page  int `form:"page" valid:"int"`
	Limit int `form:"limit" valid:"int"`
}
