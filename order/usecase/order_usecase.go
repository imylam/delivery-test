package usecase

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/imylam/delivery-test/domain"
	"github.com/imylam/delivery-test/googlemap"
)

const (
	ErrorOrderTaken                string = "order taken, you are too late"
	statusUpdateOrderStatusSuccess string = "SUCCESS"
)

type orderUsecase struct {
	orderRepo domain.OrderRepository
	mapClient googlemap.MapClient
}

// NewOrderUsecase will create new a orderUsecase object representation of domain.OrderUsecase interface
func NewOrderUsecase(userRepo domain.OrderRepository,
	mapClient googlemap.MapClient) domain.OrderUsecase {

	return &orderUsecase{
		orderRepo: userRepo,
		mapClient: mapClient,
	}
}

func (uc *orderUsecase) PlaceOrder(origins,
	destinations []string) (order *domain.Order, err error) {

	origin := strings.Join(origins, ",")
	dest := strings.Join(destinations, ",")
	dist, err := uc.mapClient.GetDistance(origin, dest)
	if err != nil {
		return
	}

	order = &domain.Order{Distance: dist, Status: domain.StatusUnassigned}
	err = uc.orderRepo.Create(order)
	if err != nil {
		return
	}

	return
}

func (uc *orderUsecase) TakeOrder(id int64) (status string, err error) {
	order, err := uc.orderRepo.FindByID(id)
	if err != nil {
		return
	}
	if order.Status == domain.StatusTaken {
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

func (uc *orderUsecase) ListOrders(page, limit int) (orders *[]domain.Order, err error) {

	offset := (page - 1) * limit
	orders, err = uc.orderRepo.FindRange(limit, offset)

	return
}
