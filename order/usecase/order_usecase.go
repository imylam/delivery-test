package usecase

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/imylam/delivery-test/configs"
	"github.com/imylam/delivery-test/order"
	"github.com/imylam/delivery-test/order/infrastructure/googlemap"
)

const (
	ErrorOrderTaken                string = "order taken, you are too late"
	statusUpdateOrderStatusSuccess string = "SUCCESS"
)

type orderUsecase struct {
	orderRepo order.OrderRepository
	mapClient googlemap.MapClient
}

// NewOrderUsecase will create new a orderUsecase object representation of order.OrderUsecase interface
func NewOrderUsecase(userRepo order.OrderRepository,
	mapClient googlemap.MapClient) order.OrderUsecase {

	return &orderUsecase{
		orderRepo: userRepo,
		mapClient: mapClient,
	}
}

func (uc *orderUsecase) PlaceOrder(origins,
	destinations []string) (newOrder *order.Order, err error) {

	origin := strings.Join(origins, ",")
	dest := strings.Join(destinations, ",")

	dist, err := getDistance(origin, dest, uc.mapClient)
	if err != nil {
		return
	}

	newOrder = &order.Order{Distance: dist, Status: order.StatusUnassigned}
	err = uc.orderRepo.Create(newOrder)
	if err != nil {
		return
	}

	return
}

func (uc *orderUsecase) TakeOrder(id int64) (status string, err error) {
	orderFound, err := uc.orderRepo.FindByID(id)
	if err != nil {
		return
	}
	if orderFound.Status == order.StatusTaken {
		err = errors.New(ErrorOrderTaken)
		return
	}

	err = uc.orderRepo.UpdateStatusByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.New(ErrorOrderTaken)
			return
		}
		return
	}

	status = "SUCCESS"
	return
}

func (uc *orderUsecase) ListOrders(page, limit int) (orders *[]order.Order, err error) {

	offset := (page - 1) * limit
	orders, err = uc.orderRepo.FindRange(limit, offset)

	return
}

func getDistance(origin, dest string, mapClient googlemap.MapClient) (int, error) {

	if configs.Get(configs.KeyAppEnv) == "integration-test" {
		return 10, nil
	} else {
		dist, err := mapClient.GetDistance(origin, dest)
		if err != nil {
			return 0, err
		}

		return dist, nil
	}
}
