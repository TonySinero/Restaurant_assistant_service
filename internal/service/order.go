package service

import (
	config "restaurant-assistant/configs"
	"restaurant-assistant/internal/domain"
	"restaurant-assistant/internal/repository"
)

type OrderService struct {
	repo repository.Order
	cfg  *config.Config
}

func NewOrderService(repo repository.Order, cfg *config.Config) *OrderService {
	return &OrderService{repo: repo, cfg: cfg}
}

func (s *OrderService) UpdateOrder(id string, input domain.UpdateOrder) error {

	if input.Status != nil && *input.Status == 3 {
		if input.DeliveryType == nil {
			return domain.ErrDelTypeNotSelected
		}

		if *input.DeliveryType == 2 {
			if input.CourierService != nil {
				if err := s.CreateOrderCS(id, input); err != nil {
					return err
				}

			} else {
				return domain.ErrCourServiceNotSelected
			}
		}
	}

	if err := s.repo.UpdateOrder(id, input); err != nil {
		return err
	}

	if err := s.UpdateOrderFD(id, input); err != nil {
		return err
	}
	return nil
}

func (s *OrderService) GetAllOrders(filter *domain.FilterOrder, limit int, offset int) (*[]domain.GetOrder, error) {
	return s.repo.GetAllOrders(filter, limit, offset)
}

func (s *OrderService) GetAllOrderStatuses() []domain.OrderStatus {
	return s.repo.GetAllOrderStatuses()
}

func (s *OrderService) GetAllOrderDeliveryTypes() []domain.OrderDeliveryType {
	return s.repo.GetAllOrderDeliveryTypes()
}

func (s *OrderService) GetOrderByID(id string) (domain.OrderByID, error) {
	return s.repo.GetOrderByID(id)
}

func (s *OrderService) CheckNewOrdersMark(id string) bool {
	return s.repo.CheckNewOrdersMark(id)
}
