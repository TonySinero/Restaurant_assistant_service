package proto

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"restaurant-assistant/internal/repository/proto"
	"restaurant-assistant/pkg/managerservice"
	"restaurant-assistant/pkg/orderservice_ra"
	"restaurant-assistant/pkg/restaurantservice"
)

type Order interface {
	CreateOrder(ctx context.Context, order *orderservice_ra.Order) (*emptypb.Empty, error)
	GetOrderTotal(ctx context.Context, input *orderservice_ra.OrderDishes) (*orderservice_ra.OrderTotal, error)
	AddRestaurantFeedback(ctx context.Context, input *orderservice_ra.OrderFeedbackOnRestaurantQuality) (*emptypb.Empty, error)
}

type Manager interface {
	CreateManager(ctx context.Context, manager *managerservice.Manager) (*managerservice.ManagerState, error)
}

type Restaurant interface {
	GetUserAddress(ctx context.Context, address *restaurantservice.UserAddress) (*restaurantservice.NearestRestaurants, error)
}

type Service struct {
	Order
	orderservice_ra.UnsafeOrderServiceServer
	Manager
	managerservice.UnsafeManagerServiceServer
	Restaurant
	restaurantservice.UnsafeRestaurantServiceServer
}

func NewService(repo *proto.Repository) *Service {
	return &Service{
		Order:      NewOrderService(repo),
		Manager:    NewManagerService(repo),
		Restaurant: NewRestaurantService(repo),
	}
}
