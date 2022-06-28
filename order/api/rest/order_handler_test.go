package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/go-sql-driver/mysql"
	"github.com/imylam/delivery-test/common/middleware"
	"github.com/imylam/delivery-test/logger"
	"github.com/imylam/delivery-test/order"
	"github.com/imylam/delivery-test/order/mocks"
	"github.com/stretchr/testify/mock"
)

func TestPlaceOrder(t *testing.T) {
	logger.Init()

	httpMethod := "POST"
	httpPath := "/orders"

	t.Run("invalid-coordinate-missing-longitude", func(t *testing.T) {
		placeOrderReq := createMockPlaceOrderRequest([]string{"22.300789"}, createValidDestination())
		jsonBytes, _ := json.Marshal(placeOrderReq)

		mockOrderUC := new(mocks.OrderUsecase)
		router := createGinRouter()
		NewOrderHandler(router, mockOrderUC)

		req, _ := http.NewRequest(httpMethod, httpPath, bytes.NewReader(jsonBytes))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "400", w.Header().Get("HTTP"))
	})

	t.Run("invalid-coordinate-not-number", func(t *testing.T) {
		placeOrderReq := createMockPlaceOrderRequest([]string{"22.300789", "abc"}, createValidDestination())
		jsonBytes, _ := json.Marshal(placeOrderReq)

		mockOrderUC := new(mocks.OrderUsecase)
		router := createGinRouter()
		NewOrderHandler(router, mockOrderUC)

		req, _ := http.NewRequest(httpMethod, httpPath, bytes.NewReader(jsonBytes))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "400", w.Header().Get("HTTP"))
	})

	t.Run("invalid-latitude", func(t *testing.T) {
		placeOrderReq := createMockPlaceOrderRequest([]string{"92.300789", "-114.167815"}, createValidDestination())
		placeOrderReq.Origin = []string{"92.300789", "-114.167815"}
		jsonBytes, _ := json.Marshal(placeOrderReq)

		mockOrderUC := new(mocks.OrderUsecase)
		router := createGinRouter()
		NewOrderHandler(router, mockOrderUC)

		req, _ := http.NewRequest(httpMethod, httpPath, bytes.NewReader(jsonBytes))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "400", w.Header().Get("HTTP"))
	})

	t.Run("invalid-longitude", func(t *testing.T) {
		placeOrderReq := createMockPlaceOrderRequest([]string{"22.300789", "-214.167815"}, createValidDestination())
		jsonBytes, _ := json.Marshal(placeOrderReq)

		mockOrderUC := new(mocks.OrderUsecase)
		router := createGinRouter()
		NewOrderHandler(router, mockOrderUC)

		req, _ := http.NewRequest(httpMethod, httpPath, bytes.NewReader(jsonBytes))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "400", w.Header().Get("HTTP"))
	})

	t.Run("db-error", func(t *testing.T) {
		tempMockRequest := createValidPlaceOrderRequest()
		jsonBytes, _ := json.Marshal(tempMockRequest)

		mockOrderUC := new(mocks.OrderUsecase)
		mockOrderUC.On("PlaceOrder", mock.AnythingOfType("[]string"),
			mock.AnythingOfType("[]string")).Return(nil, &mysql.MySQLError{})
		router := createGinRouter()
		NewOrderHandler(router, mockOrderUC)

		req, _ := http.NewRequest(httpMethod, httpPath, bytes.NewReader(jsonBytes))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "500", w.Header().Get("HTTP"))
		mockOrderUC.AssertExpectations(t)
	})
}

func TestTakeOrder(t *testing.T) {
	httpMethod := "PATCH"
	httpPath := "/orders/1"

	t.Run("uri-param-not-digit", func(t *testing.T) {
		tempHTTPPath := "/orders/aa"

		mockRequest := createMockTakeOrderRequest()
		jsonBytes, _ := json.Marshal(mockRequest)

		mockOrderUC := new(mocks.OrderUsecase)
		router := createGinRouter()
		NewOrderHandler(router, mockOrderUC)

		req, _ := http.NewRequest(httpMethod, tempHTTPPath, bytes.NewReader(jsonBytes))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "400", w.Header().Get("HTTP"))
	})

	t.Run("invalid-request-body", func(t *testing.T) {
		mockRequest := createMockTakeOrderRequest()
		mockRequest.Status = "HELLO"
		jsonBytes, _ := json.Marshal(mockRequest)

		mockOrderUC := new(mocks.OrderUsecase)
		router := createGinRouter()
		NewOrderHandler(router, mockOrderUC)

		req, _ := http.NewRequest(httpMethod, httpPath, bytes.NewReader(jsonBytes))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "400", w.Header().Get("HTTP"))
	})

	t.Run("db-error", func(t *testing.T) {
		mockRequest := createMockTakeOrderRequest()
		jsonBytes, _ := json.Marshal(mockRequest)

		mockOrderUC := new(mocks.OrderUsecase)
		mockOrderUC.On("TakeOrder", mock.AnythingOfType("int64")).Return("", &mysql.MySQLError{})
		router := createGinRouter()
		NewOrderHandler(router, mockOrderUC)

		req, _ := http.NewRequest(httpMethod, httpPath, bytes.NewReader(jsonBytes))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "500", w.Header().Get("HTTP"))
		mockOrderUC.AssertExpectations(t)
	})
}

func TestListOrders(t *testing.T) {
	httpMethod := "GET"
	httpPath := "/orders"

	t.Run("empty-result", func(t *testing.T) {
		qParams := buildListOrderQueryParams(100, 100)

		expJSONRespBytes, _ := json.Marshal([]string{})

		mockOrderUC := new(mocks.OrderUsecase)
		mockOrderUC.On("ListOrders", mock.AnythingOfType("int"),
			mock.AnythingOfType("int")).Return(&[]order.Order{}, nil)
		router := createGinRouter()
		NewOrderHandler(router, mockOrderUC)

		req, _ := http.NewRequest(httpMethod, httpPath+qParams, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "200", w.Header().Get("HTTP"))
		assert.Equal(t, string(expJSONRespBytes), w.Body.String())
		mockOrderUC.AssertExpectations(t)
	})

	t.Run("invalid-qparams", func(t *testing.T) {
		mockPage := "aaa"
		mockLimit := 4
		qParams := fmt.Sprintf("?page=%s&limit=%d", mockPage, mockLimit)

		mockOrderUC := new(mocks.OrderUsecase)
		router := createGinRouter()
		NewOrderHandler(router, mockOrderUC)

		req, _ := http.NewRequest(httpMethod, httpPath+qParams, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "400", w.Header().Get("HTTP"))
	})

	t.Run("db-error", func(t *testing.T) {
		qParams := buildListOrderQueryParams(1, 4)

		mockOrderUC := new(mocks.OrderUsecase)
		mockOrderUC.On("ListOrders", mock.AnythingOfType("int"),
			mock.AnythingOfType("int")).Return(nil, &mysql.MySQLError{})
		router := createGinRouter()
		NewOrderHandler(router, mockOrderUC)

		req, _ := http.NewRequest(httpMethod, httpPath+qParams, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "500", w.Header().Get("HTTP"))
		mockOrderUC.AssertExpectations(t)
	})
}

func TestValidatePlaceOrder(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRequest := PlaceOrderRequest{
			Origin:      []string{"22.300789", "114.167815"},
			Destination: []string{"22.33540", "114.176155"},
		}

		isValid, _ := validatePlaceOrder(mockRequest)

		assert.Equal(t, true, isValid)
	})

	t.Run("coordinate-not-two", func(t *testing.T) {
		mockRequest := PlaceOrderRequest{
			Origin:      []string{"22.300789", "114.167815", "114.167815"},
			Destination: []string{"22.33540", "114.176155"},
		}

		isValid, s := validatePlaceOrder(mockRequest)

		assert.Equal(t, false, isValid)
		assert.Equal(t, errInvalidCoordinates, s)
	})

	t.Run("coordinate-not-digit", func(t *testing.T) {
		mockRequest := PlaceOrderRequest{
			Origin:      []string{"22.300789", "114.167815"},
			Destination: []string{"22.33540", "aaa"},
		}

		isValid, s := validatePlaceOrder(mockRequest)

		assert.Equal(t, false, isValid)
		assert.Equal(t, errInvalidCoordinates, s)
	})

	t.Run("invalid-latitude", func(t *testing.T) {
		mockRequest := PlaceOrderRequest{
			Origin:      []string{"122.300789", "114.167815"},
			Destination: []string{"22.33540", "114.176155"},
		}

		isValid, s := validatePlaceOrder(mockRequest)

		assert.Equal(t, false, isValid)
		assert.Equal(t, errInvalidCoordinates, s)
	})

	t.Run("invalid-longtitude", func(t *testing.T) {
		mockRequest := PlaceOrderRequest{
			Origin:      []string{"22.300789", "114.167815"},
			Destination: []string{"22.33540", "-214.176155"},
		}

		isValid, s := validatePlaceOrder(mockRequest)

		assert.Equal(t, false, isValid)
		assert.Equal(t, errInvalidCoordinates, s)
	})
}

func createGinRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(middleware.HandleRestError)

	return router
}

func createValidPlaceOrderRequest() PlaceOrderRequest {
	return createMockPlaceOrderRequest(
		[]string{"22.300789", "114.167815"},
		[]string{"22.33540", "114.176155"},
	)
}

func createMockPlaceOrderRequest(ori []string, dest []string) PlaceOrderRequest {
	return PlaceOrderRequest{
		Origin:      ori,
		Destination: dest,
	}
}

func createValidOrigin() []string {
	return []string{"22.300789", "114.167815"}
}

func createValidDestination() []string {
	return []string{"22.33540", "114.176155"}
}

func createMockTakeOrderRequest() TakeOrderRequest {
	return TakeOrderRequest{Status: "TAKEN"}
}

func buildListOrderQueryParams(page int, limit int) string {
	return fmt.Sprintf("?page=%d&limit=%d", page, limit)
}
