package rest

import (
	"log"
	"net/http"

	resterrors "github.com/imylam/delivery-test/common/rest_errors"
	"github.com/imylam/delivery-test/domain"
	"github.com/imylam/delivery-test/order/usecase"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

const (
	errInvalidCoordinates    string = "invalid coordinates"
	errInvalidResquestParams string = "invalid request params"
	errInternalServer        string = "internal server error"
)

// orderHandler represents the httphandler for handling requests relating to Orders
type orderHandler struct {
	orderUC domain.OrderUsecase
}

// NewOrderHandler will initialize the Order endpoints
func NewOrderHandler(g *gin.Engine, orderUC domain.OrderUsecase) {
	handler := &orderHandler{
		orderUC: orderUC,
	}

	g.POST("/orders", handler.placeOrder)
	g.PATCH("/orders/:id", handler.takeOrder)
	g.GET("/orders", handler.listOrder)
}

func (h *orderHandler) placeOrder(c *gin.Context) {
	var req PlaceOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(resterrors.NewBadRequestError(errInvalidResquestParams))
		return
	}

	isValid, errMsg := validatePlaceOrder(req)
	if !isValid {
		c.Error(resterrors.NewBadRequestError(errMsg))
		return
	}

	order, err := h.orderUC.PlaceOrder(req.Origin, req.Destination)
	if err != nil {
		log.Print(err.Error())

		c.Error(resterrors.NewInternalServerError(errInternalServer))
		return
	}

	c.Header("HTTP", "200")
	c.JSON(http.StatusOK, order)
}

func (h *orderHandler) takeOrder(c *gin.Context) {
	var req TakeOrderRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.Error(resterrors.NewBadRequestError(errInvalidResquestParams))
		return
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(resterrors.NewBadRequestError(errInvalidResquestParams))
		return
	}

	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		c.Error(resterrors.NewBadRequestError(errInvalidResquestParams))
		return
	}

	if req.Status != "TAKEN" {
		c.Error(resterrors.NewBadRequestError(errInvalidResquestParams))
		return
	}

	status, err := h.orderUC.TakeOrder(req.ID)
	if err != nil {
		if err.Error() != usecase.ErrorOrderTaken {
			log.Print(err.Error())

			c.Error(resterrors.NewInternalServerError(errInternalServer))
			return
		}

		c.Header("HTTP", "409")
		c.JSON(http.StatusConflict, gin.H{"error": usecase.ErrorOrderTaken})
		log.Print(err.Error() == usecase.ErrorOrderTaken)
		return
	}

	c.Header("HTTP", "200")
	c.JSON(http.StatusOK, gin.H{"status": status})
}

func (h *orderHandler) listOrder(c *gin.Context) {
	var req ListOrderRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(resterrors.NewBadRequestError(errInvalidResquestParams))
		return
	}

	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		c.Error(resterrors.NewBadRequestError(err.Error()))
		return
	}

	orders, err := h.orderUC.ListOrders(req.Page, req.Limit)
	if err != nil {
		log.Print(err.Error())
		c.Error(resterrors.NewInternalServerError(errInternalServer))
		return
	}

	if len(*orders) == 0 {
		c.Header("HTTP", "200")
		c.JSON(http.StatusOK, []string{})
		return
	}

	c.Header("HTTP", "200")
	c.JSON(http.StatusOK, orders)
}

// validatePlaceOrder checks where coodinates are string that can be converted to float64
func validatePlaceOrder(req PlaceOrderRequest) (bool, string) {
	originInterface := make([]interface{}, len(req.Origin))
	for i, v := range req.Origin {
		originInterface[i] = v
	}

	destInterface := make([]interface{}, len(req.Destination))
	for i, v := range req.Destination {
		destInterface[i] = v
	}

	var fn govalidator.ConditionIterator = func(value interface{}, index int) bool {
		s, ok := value.(string)
		if !ok {
			return false
		}

		// first number of coordinate is latitude
		if index == 0 {
			return govalidator.IsLatitude(s)
		}

		return govalidator.IsLongitude(s)
	}

	if govalidator.Count(originInterface, fn) != 2 {
		return false, errInvalidCoordinates
	}
	if !govalidator.ValidateArray(originInterface, fn) {
		return false, errInvalidCoordinates
	}

	if govalidator.Count(destInterface, fn) != 2 {
		return false, errInvalidCoordinates
	}
	if !govalidator.ValidateArray(destInterface, fn) {
		return false, errInvalidCoordinates
	}

	return true, ""
}
