package http

// PlaceOrderReponse represents the place order reponse body
type PlaceOrderReponse struct {
	ID       int    `json:"id"`
	Distance int    `json:"distance"`
	Status   string `json:"status"`
}

// TakeOrderResponse rrepresents the take order reponse body
type TakeOrderResponse struct {
	Status string `json:"status"`
}
