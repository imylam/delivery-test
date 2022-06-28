package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/imylam/delivery-test/domain"
	"github.com/imylam/delivery-test/domain/mocks"
	"github.com/stretchr/testify/mock"
)

func TestPlaceOrder(t *testing.T) {
	httpMethod := "POST"
	httpPath := "/orders"
	mockRequest := PlaceOrderRequest{
		Origin:      []string{"22.300789", "114.167815"},
		Destination: []string{"22.33540", "114.176155"},
	}
	mockOrder := domain.Order{
		ID:       1,
		Distance: 888,
		Status:   domain.StatusUnassigned,
	}

	t.Run("success", func(t *testing.T) {
		tempMockOrder := mockOrder
		expJSONRespBytes, _ := json.Marshal(tempMockOrder)

		tempMockRequest := mockRequest
		jsonBytes, _ := json.Marshal(tempMockRequest)

		mockOrderUC := new(mocks.OrderUsecase)
		mockOrderUC.On("PlaceOrder", mock.AnythingOfType("[]string"),
			mock.AnythingOfType("[]string")).Return(&tempMockOrder, nil)

		gin.SetMode(gin.TestMode)
		router := gin.Default()

		NewOrderHandler(router, mockOrderUC)

		req, _ := http.NewRequest(httpMethod, httpPath, bytes.NewReader(jsonBytes))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("TestPlaceOrder() fails, expect response code %d, got: %d", http.StatusOK, w.Code)
		}
		if w.Header().Get("HTTP") != "200" {
			t.Errorf("TestPlaceOrder() fails, expect header HTTP: %s, got: %s", "200", w.Header().Get("HTTP"))
		}
		if w.Body.String() != string(expJSONRespBytes) {
			t.Errorf("TestRegister() fails, expect response: %s, got: %s", string(expJSONRespBytes), w.Body.String())
		}

		mockOrderUC.AssertExpectations(t)
	})

	t.Run("invalid-coordinate-missing-longitude", func(t *testing.T) {
		tempMockRequest := mockRequest
		tempMockRequest.Origin = []string{"22.300789"}
		jsonBytes, _ := json.Marshal(tempMockRequest)

		mockOrderUC := new(mocks.OrderUsecase)

		gin.SetMode(gin.TestMode)
		router := gin.Default()

		NewOrderHandler(router, mockOrderUC)

		req, _ := http.NewRequest(httpMethod, httpPath, bytes.NewReader(jsonBytes))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("TestPlaceOrder() fails, expect response code %d, got: %d", http.StatusBadRequest, w.Code)
		}
		if w.Header().Get("HTTP") != "400" {
			t.Errorf("TestPlaceOrder() fails, expect header HTTP: %s, got: %s", "400", w.Header().Get("HTTP"))
		}
	})

	t.Run("invalid-coordinate-not-number", func(t *testing.T) {
		tempMockRequest := mockRequest
		tempMockRequest.Origin = []string{"22.300789", "abc"}
		jsonBytes, _ := json.Marshal(tempMockRequest)

		mockOrderUC := new(mocks.OrderUsecase)

		gin.SetMode(gin.TestMode)
		router := gin.Default()

		NewOrderHandler(router, mockOrderUC)

		req, _ := http.NewRequest(httpMethod, httpPath, bytes.NewReader(jsonBytes))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("TestPlaceOrder() fails, expect response code %d, got: %d", http.StatusBadRequest, w.Code)
		}
		if w.Header().Get("HTTP") != "400" {
			t.Errorf("TestPlaceOrder() fails, expect header HTTP: %s, got: %s", "400", w.Header().Get("HTTP"))
		}
	})

	t.Run("invalid-latitude", func(t *testing.T) {
		tempMockRequest := mockRequest
		tempMockRequest.Origin = []string{"92.300789", "-114.167815"}
		jsonBytes, _ := json.Marshal(tempMockRequest)

		mockOrderUC := new(mocks.OrderUsecase)

		gin.SetMode(gin.TestMode)
		router := gin.Default()

		NewOrderHandler(router, mockOrderUC)

		req, _ := http.NewRequest(httpMethod, httpPath, bytes.NewReader(jsonBytes))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("TestPlaceOrder() fails, expect response code %d, got: %d", http.StatusBadRequest, w.Code)
		}
		if w.Header().Get("HTTP") != "400" {
			t.Errorf("TestPlaceOrder() fails, expect header HTTP: %s, got: %s", "400", w.Header().Get("HTTP"))
		}
	})

	t.Run("invalid-longitude", func(t *testing.T) {
		tempMockRequest := mockRequest
		tempMockRequest.Origin = []string{"22.300789", "-214.167815"}
		jsonBytes, _ := json.Marshal(tempMockRequest)

		mockOrderUC := new(mocks.OrderUsecase)

		gin.SetMode(gin.TestMode)
		router := gin.Default()

		NewOrderHandler(router, mockOrderUC)

		req, _ := http.NewRequest(httpMethod, httpPath, bytes.NewReader(jsonBytes))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("TestPlaceOrder() fails, expect response code %d, got: %d", http.StatusBadRequest, w.Code)
		}
		if w.Header().Get("HTTP") != "400" {
			t.Errorf("TestPlaceOrder() fails, expect header HTTP: %s, got: %s", "400", w.Header().Get("HTTP"))
		}
	})

	t.Run("db-error", func(t *testing.T) {
		tempMockRequest := mockRequest
		jsonBytes, _ := json.Marshal(tempMockRequest)

		mockOrderUC := new(mocks.OrderUsecase)
		mockOrderUC.On("PlaceOrder", mock.AnythingOfType("[]string"),
			mock.AnythingOfType("[]string")).Return(nil, &mysql.MySQLError{})

		gin.SetMode(gin.TestMode)
		router := gin.Default()

		NewOrderHandler(router, mockOrderUC)

		req, _ := http.NewRequest(httpMethod, httpPath, bytes.NewReader(jsonBytes))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("TestPlaceOrder() fails, expect response code %d, got: %d", http.StatusInternalServerError, w.Code)
		}
		if w.Header().Get("HTTP") != "500" {
			t.Errorf("TestPlaceOrder() fails, expect header HTTP: %s, got: %s", "500", w.Header().Get("HTTP"))
		}

		mockOrderUC.AssertExpectations(t)
	})
}

func TestTakeOrder(t *testing.T) {
	httpMethod := "PATCH"
	httpPath := "/orders/1"
	mockRequest := TakeOrderRequest{
		Status: "TAKEN",
	}

	t.Run("success", func(t *testing.T) {
		resp := make(map[string]interface{}, 1)
		resp["status"] = "SUCCESS"
		expJSONRespBytes, _ := json.Marshal(resp)

		tempMockRequest := mockRequest
		jsonBytes, _ := json.Marshal(tempMockRequest)

		mockOrderUC := new(mocks.OrderUsecase)
		mockOrderUC.On("TakeOrder", mock.AnythingOfType("int64")).Return("SUCCESS", nil)

		gin.SetMode(gin.TestMode)
		router := gin.Default()

		NewOrderHandler(router, mockOrderUC)

		req, _ := http.NewRequest(httpMethod, httpPath, bytes.NewReader(jsonBytes))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("TestPlaceOrder() fails, expect response code %d, got: %d", http.StatusOK, w.Code)
		}
		if w.Header().Get("HTTP") != "200" {
			t.Errorf("TestPlaceOrder() fails, expect header HTTP: %s, got: %s", "200", w.Header().Get("HTTP"))
		}
		if w.Body.String() != string(expJSONRespBytes) {
			t.Errorf("TestRegister() fails, expect response: %s, got: %s", string(expJSONRespBytes), w.Body.String())
		}

		mockOrderUC.AssertExpectations(t)
	})

	t.Run("uri-param-not-digit", func(t *testing.T) {
		tempHTTPPath := "/orders/aa"

		tempMockRequest := mockRequest
		jsonBytes, _ := json.Marshal(tempMockRequest)

		mockOrderUC := new(mocks.OrderUsecase)

		gin.SetMode(gin.TestMode)
		router := gin.Default()

		NewOrderHandler(router, mockOrderUC)

		req, _ := http.NewRequest(httpMethod, tempHTTPPath, bytes.NewReader(jsonBytes))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("TestPlaceOrder() fails, expect response code %d, got: %d", http.StatusBadRequest, w.Code)
		}
		if w.Header().Get("HTTP") != "400" {
			t.Errorf("TestPlaceOrder() fails, expect header HTTP: %s, got: %s", "400", w.Header().Get("HTTP"))
		}
	})

	t.Run("invalid-request-body", func(t *testing.T) {
		tempMockRequest := mockRequest
		tempMockRequest.Status = "HELLO"
		jsonBytes, _ := json.Marshal(tempMockRequest)

		mockOrderUC := new(mocks.OrderUsecase)
		// mockOrderUC.On("TakeOrder", mock.AnythingOfType("int64")).Return("SUCCESS", nil)

		gin.SetMode(gin.TestMode)
		router := gin.Default()

		NewOrderHandler(router, mockOrderUC)

		req, _ := http.NewRequest(httpMethod, httpPath, bytes.NewReader(jsonBytes))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("TestPlaceOrder() fails, expect response code %d, got: %d", http.StatusBadRequest, w.Code)
		}
		if w.Header().Get("HTTP") != "400" {
			t.Errorf("TestPlaceOrder() fails, expect header HTTP: %s, got: %s", "400", w.Header().Get("HTTP"))
		}
	})

	t.Run("db-error", func(t *testing.T) {
		tempMockRequest := mockRequest
		jsonBytes, _ := json.Marshal(tempMockRequest)

		mockOrderUC := new(mocks.OrderUsecase)
		mockOrderUC.On("TakeOrder", mock.AnythingOfType("int64")).Return("", &mysql.MySQLError{})

		gin.SetMode(gin.TestMode)
		router := gin.Default()

		NewOrderHandler(router, mockOrderUC)

		req, _ := http.NewRequest(httpMethod, httpPath, bytes.NewReader(jsonBytes))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("TestPlaceOrder() fails, expect response code %d, got: %d", http.StatusInternalServerError, w.Code)
		}
		if w.Header().Get("HTTP") != "500" {
			t.Errorf("TestPlaceOrder() fails, expect header HTTP: %s, got: %s", "500", w.Header().Get("HTTP"))
		}

		mockOrderUC.AssertExpectations(t)
	})
}

func TestListOrders(t *testing.T) {
	httpMethod := "GET"
	httpPath := "/orders"
	mockOrders := []domain.Order{
		domain.Order{ID: 1, Distance: 100, Status: domain.StatusTaken},
		domain.Order{ID: 2, Distance: 200, Status: domain.StatusUnassigned},
		domain.Order{ID: 3, Distance: 300, Status: domain.StatusUnassigned},
		domain.Order{ID: 4, Distance: 400, Status: domain.StatusTaken},
	}

	t.Run("success", func(t *testing.T) {
		mockPage := 1
		mockLimit := 4
		qParams := fmt.Sprintf("?page=%d&limit=%d", mockPage, mockLimit)

		tempMockOrders := mockOrders
		expJSONRespBytes, _ := json.Marshal(tempMockOrders)

		mockOrderUC := new(mocks.OrderUsecase)
		mockOrderUC.On("ListOrders", mock.AnythingOfType("int"),
			mock.AnythingOfType("int")).Return(&mockOrders, nil)

		gin.SetMode(gin.TestMode)
		router := gin.Default()

		NewOrderHandler(router, mockOrderUC)

		req, _ := http.NewRequest(httpMethod, httpPath+qParams, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("TestListOrders() fails, expect response code %d, got: %d", http.StatusOK, w.Code)
		}
		if w.Header().Get("HTTP") != "200" {
			t.Errorf("TestPlaceOrder() fails, expect header HTTP: %s, got: %s", "200", w.Header().Get("HTTP"))
		}
		if w.Body.String() != string(expJSONRespBytes) {
			t.Errorf("TestListOrders() fails, expect response: %s, got: %s",
				string(expJSONRespBytes), w.Body.String())
		}

		mockOrderUC.AssertExpectations(t)
	})

	t.Run("empty-result", func(t *testing.T) {
		mockPage := 100
		mockLimit := 100
		qParams := fmt.Sprintf("?page=%d&limit=%d", mockPage, mockLimit)

		expJSONRespBytes, _ := json.Marshal([]string{})

		mockOrderUC := new(mocks.OrderUsecase)
		mockOrderUC.On("ListOrders", mock.AnythingOfType("int"),
			mock.AnythingOfType("int")).Return(&[]domain.Order{}, nil)

		gin.SetMode(gin.TestMode)
		router := gin.Default()

		NewOrderHandler(router, mockOrderUC)

		req, _ := http.NewRequest(httpMethod, httpPath+qParams, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("TestListOrders() fails, expect response code %d, got: %d", http.StatusOK, w.Code)
		}
		if w.Header().Get("HTTP") != "200" {
			t.Errorf("TestListOrders() fails, expect header HTTP: %s, got: %s", "200", w.Header().Get("HTTP"))
		}
		if w.Body.String() != string(expJSONRespBytes) {
			t.Errorf("TestListOrders() fails, expect response: %s, got: %s",
				string(expJSONRespBytes), w.Body.String())
		}

		mockOrderUC.AssertExpectations(t)
	})

	t.Run("invalid-qparams", func(t *testing.T) {
		mockPage := "aaa"
		mockLimit := 4
		qParams := fmt.Sprintf("?page=%s&limit=%d", mockPage, mockLimit)

		mockOrderUC := new(mocks.OrderUsecase)

		gin.SetMode(gin.TestMode)
		router := gin.Default()

		NewOrderHandler(router, mockOrderUC)

		req, _ := http.NewRequest(httpMethod, httpPath+qParams, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("TestListOrders() fails, expect response code %d, got: %d", http.StatusBadRequest, w.Code)
		}
		if w.Header().Get("HTTP") != "400" {
			t.Errorf("TestPlaceOrder() fails, expect header HTTP: %s, got: %s", "400", w.Header().Get("HTTP"))
		}
	})

	t.Run("invalid-qparams", func(t *testing.T) {
		mockPage := 1
		mockLimit := 4
		qParams := fmt.Sprintf("?page=%d&limit=%d", mockPage, mockLimit)

		mockOrderUC := new(mocks.OrderUsecase)
		mockOrderUC.On("ListOrders", mock.AnythingOfType("int"),
			mock.AnythingOfType("int")).Return(nil, &mysql.MySQLError{})

		gin.SetMode(gin.TestMode)
		router := gin.Default()

		NewOrderHandler(router, mockOrderUC)

		req, _ := http.NewRequest(httpMethod, httpPath+qParams, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("TestPlaceOrder() fails, expect response code %d, got: %d", http.StatusInternalServerError, w.Code)
		}
		if w.Header().Get("HTTP") != "500" {
			t.Errorf("TestPlaceOrder() fails, expect header HTTP: %s, got: %s", "500", w.Header().Get("HTTP"))
		}

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

		if !isValid {
			t.Errorf("TestValidatePlaceOrder(() fails, expected valid, but got invalid")
		}
	})

	t.Run("coordinate-not-two", func(t *testing.T) {
		mockRequest := PlaceOrderRequest{
			Origin:      []string{"22.300789", "114.167815", "114.167815"},
			Destination: []string{"22.33540", "114.176155"},
		}

		isValid, s := validatePlaceOrder(mockRequest)

		if isValid {
			t.Errorf("TestValidatePlaceOrder(() fails, expected imvalid, but got valid")
		}
		if s != errInvalidCoordinates {
			t.Errorf("TestValidatePlaceOrder(() fails, expected %s, but got %s", errInvalidCoordinates, s)
		}
	})

	t.Run("coordinate-not-digit", func(t *testing.T) {
		mockRequest := PlaceOrderRequest{
			Origin:      []string{"22.300789", "114.167815"},
			Destination: []string{"22.33540", "aaa"},
		}

		isValid, s := validatePlaceOrder(mockRequest)

		if isValid {
			t.Errorf("TestValidatePlaceOrder(() fails, expected imvalid, but got valid")
		}
		if s != errInvalidCoordinates {
			t.Errorf("TestValidatePlaceOrder(() fails, expected %s, but got %s", errInvalidCoordinates, s)
		}
	})

	t.Run("invalid-latitude", func(t *testing.T) {
		mockRequest := PlaceOrderRequest{
			Origin:      []string{"122.300789", "114.167815"},
			Destination: []string{"22.33540", "114.176155"},
		}

		isValid, s := validatePlaceOrder(mockRequest)

		if isValid {
			t.Errorf("TestValidatePlaceOrder(() fails, expected imvalid, but got valid")
		}
		if s != errInvalidCoordinates {
			t.Errorf("TestValidatePlaceOrder(() fails, expected %s, but got %s", errInvalidCoordinates, s)
		}
	})

	t.Run("invalid-longtitude", func(t *testing.T) {
		mockRequest := PlaceOrderRequest{
			Origin:      []string{"22.300789", "114.167815"},
			Destination: []string{"22.33540", "-214.176155"},
		}

		isValid, s := validatePlaceOrder(mockRequest)

		if isValid {
			t.Errorf("TestValidatePlaceOrder(() fails, expected imvalid, but got valid")
		}
		if s != errInvalidCoordinates {
			t.Errorf("TestValidatePlaceOrder(() fails, expected %s, but got %s", errInvalidCoordinates, s)
		}
	})
}
