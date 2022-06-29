//go:build integration
// +build integration

package integrationtests_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/go-resty/resty/v2"

	"github.com/imylam/delivery-test/order"
	"github.com/imylam/delivery-test/order/api/rest"
)

func Test_ListOrders(t *testing.T) {

	client := resty.New()

	t.Run("GIVEN_no_orders_WHEN_list_order_THEN_empty_array_json_should_be_returned", func(t *testing.T) {

		resp := listOrders(1, 5, client)

		assert.Equal(t, 200, resp.StatusCode())
		assert.Equal(t, "200", resp.Header().Get("HTTP"))
		assert.Equal(t, "[]", string(resp.Body()))
	})

	t.Run("GIVEN_ten_orders_WHEN_list_order_with_page_and_limit_THEN_orders_according_to_page_limit_settings_should_be_returned", func(t *testing.T) {

		for i := 0; i < 10; i++ {
			placeOrder(&rest.PlaceOrderReponse{}, client)
		}

		// Test limit
		resp := listOrders(1, 5, client)
		var orders []order.Order
		_ = json.Unmarshal(resp.Body(), &orders)

		assert.Equal(t, 5, len(orders))

		// Test page
		firstOrderId := orders[0].ID
		resp2 := listOrders(2, 5, client)
		var orders2 []order.Order
		_ = json.Unmarshal(resp2.Body(), &orders2)

		assert.Equal(t, 5, len(orders2))
		assert.Equal(t, firstOrderId+5, orders2[0].ID)
	})
}

func Test_PlaceOrders(t *testing.T) {

	client := resty.New()

	t.Run("GIVEN_vaild_PlaceOrderRequest_body_WHEN_place_order_THEN_success_placeOrderResponse_should_be_returned", func(t *testing.T) {

		placeOrderResponose := &rest.PlaceOrderReponse{}

		resp := placeOrder(placeOrderResponose, client)

		assert.Equal(t, 200, resp.StatusCode())
		assert.Equal(t, "200", resp.Header().Get("HTTP"))
		assert.Equal(t, true, placeOrderResponose.ID > 0)
		assert.Equal(t, placeOrderResponose.Status, "UNASSIGNED")
	})
}

func Test_TakeOrder(t *testing.T) {

	client := resty.New()

	t.Run("GIVEN_order_untaken_WHEN_take_order_THEN_success_take_order_response_should_be_returned", func(t *testing.T) {

		placeOrderResponose := &rest.PlaceOrderReponse{}
		placeOrder(placeOrderResponose, client)

		orderId := placeOrderResponose.ID
		takeOrderResponse := &rest.TakeOrderResponse{}
		resp := takeOrder(orderId, takeOrderResponse, client)

		assert.Equal(t, 200, resp.StatusCode())
		assert.Equal(t, "200", resp.Header().Get("HTTP"))
		assert.Equal(t, "SUCCESS", takeOrderResponse.Status)
	})

	t.Run("GIVEN_order_taken_WHEN_take_order_THEN_conflict_take_order_response_should_be_returned", func(t *testing.T) {

		placeOrderResponose := &rest.PlaceOrderReponse{}
		placeOrder(placeOrderResponose, client)

		orderId := placeOrderResponose.ID
		takeOrder(orderId, &rest.TakeOrderResponse{}, client)

		takeOrderResponse := &rest.TakeOrderResponse{}
		resp := takeOrder(orderId, takeOrderResponse, client)

		assert.Equal(t, 409, resp.StatusCode())
		assert.Equal(t, "409", resp.Header().Get("HTTP"))
		assert.Equal(t, `{"error":"order taken, you are too late"}`, string(resp.Body()))
	})
}

func listOrders(page int, limit int, client *resty.Client) (resp *resty.Response) {
	resp, _ = client.R().
		SetHeader("Content-Type", "application/json").
		Get(fmt.Sprintf("%s/orders?page=%d&limit=%d", getBaseUrl(), page, limit))

	return
}

func placeOrder(placeOrderResponose *rest.PlaceOrderReponse, client *resty.Client) (resp *resty.Response) {
	resp, _ = client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"origin": ["0.00", "0.00"], "destination": ["1.00", "0.00"]}`).
		SetResult(placeOrderResponose).
		Post(fmt.Sprintf("%s/orders", getBaseUrl()))

	return
}

func takeOrder(orderId int, takeOrderResponse *rest.TakeOrderResponse, client *resty.Client) (resp *resty.Response) {
	resp, _ = client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"status":"TAKEN"}`).
		SetResult(takeOrderResponse).
		Patch(fmt.Sprintf("%s/orders/%d", getBaseUrl(), orderId))

	return
}

func getBaseUrl() string {
	appUrl := "http://localhost:8080"

	if appUrlFromEnv, isFound := os.LookupEnv("APP_URL"); isFound {
		appUrl = appUrlFromEnv
	}

	return appUrl
}
