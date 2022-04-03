package proto

import (
	"context"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"restaurant-assistant/internal/domain"
	"restaurant-assistant/internal/repository/proto"
	"restaurant-assistant/pkg/orderservice_ra"
	"restaurant-assistant/pkg/validation"
)

type OrderService struct {
	repo *proto.Repository
}

func NewOrderService(repo *proto.Repository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(ctx context.Context, order *orderservice_ra.Order) (*emptypb.Empty, error) {
	if err := s.ValidateOrder(order); err != nil {

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return s.repo.CreateOrder(ctx, order)
}

func (s *OrderService) AddRestaurantFeedback(ctx context.Context, input *orderservice_ra.OrderFeedbackOnRestaurantQuality) (*emptypb.Empty, error) {
	return s.repo.AddRestaurantFeedback(ctx, input)
}

func (s *OrderService) ValidateOrder(order *orderservice_ra.Order) error {
	if ok := validation.IsValidUUID(order.RestaurantID); !ok {
		log.Error().Msg("invalid restaurant id")

		return domain.ErrInvalidRestaurantID

	}

	return nil
}

func (s *OrderService) GetOrderTotal(ctx context.Context, input *orderservice_ra.OrderDishes) (*orderservice_ra.OrderTotal, error) {
	return s.repo.GetOrderTotal(ctx, input)
}
