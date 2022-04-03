package proto

import (
	"context"
	"github.com/jmoiron/sqlx"
	"google.golang.org/protobuf/types/known/emptypb"
	"restaurant-assistant/internal/domain"
	"restaurant-assistant/pkg/orderservice_ra"
)

type Order interface {
	CreateOrder(ctx context.Context, order *orderservice_ra.Order) (*emptypb.Empty, error)
	GetOrderTotal(ctx context.Context, input *orderservice_ra.OrderDishes) (*orderservice_ra.OrderTotal, error)
	AddRestaurantFeedback(ctx context.Context, input *orderservice_ra.OrderFeedbackOnRestaurantQuality) (*emptypb.Empty, error)
}

type Restaurant interface {
	GetRestaurantsInfo(ctx context.Context, lat, lng float64) ([]domain.GetRestaurant, error)
}

type Repository struct {
	Order
	Restaurant
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Order:      NewOrderPostgres(db),
		Restaurant: NewRestaurantPostgres(db),
	}
}
