package service

import (
	"context"
	config "restaurant-assistant/configs"
	"restaurant-assistant/internal/domain"
	"restaurant-assistant/internal/repository"
	"restaurant-assistant/pkg/courierProto"
	"restaurant-assistant/pkg/storage"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Restaurant interface {
	CreateRestaurant(input domain.Restaurant) (string, error)
	UpdateRestaurant(id string, input domain.UpdateRestaurant) (domain.GetRestaurant, error)
	GetAllRestaurant(input domain.GetRestaurantOrderBy) ([]domain.GetRestaurant, error)
	DeleteRestaurant(id string) error
	GetRestaurantById(id string) (domain.GetRestaurant, error)
	GetNearestRestaurant(input domain.GetRestaurantByAddress) ([]domain.GetRestaurant, error)
	GetRestaurantCategoriesWithRestaurants() ([]domain.GetRestaurantCategories, error)
	GetRestaurantCategories() ([]string, error)
	GetRestaurantsByCategory(input domain.GetRestaurantOrderBy) ([]domain.GetRestaurant, error)
	GetRestaurantFeedbacksById(id string) ([]domain.Feedback, error)
}

type Dish interface {
	CreateDish(input domain.Dish, id string) (string, error)
	UpdateDish(id string, input domain.UpdateDish) (domain.GetDishByID, error)
	GetAllDishes(id string) ([]domain.GetAllDishes, error)
	DeleteDish(id string) error
	GetDishByID(id string) (domain.GetDishByID, error)
	GetDishByRestaurantID(id string) ([]domain.GetAllDishes, error)
	GetDishWithCategoryByRestaurantID(id string) ([]domain.GetDishesByRestaurant, error)
	GetDishesTypes() ([]domain.DishesCategory, error)
}

type File interface {
	UploadAndSaveFile(ctx context.Context, file domain.File, uuid string, path string) (string, error)
}

type Order interface {
	UpdateOrder(id string, input domain.UpdateOrder) error
	GetAllOrders(filter *domain.FilterOrder, limit int, offset int) (*[]domain.GetOrder, error)
	GetAllOrderStatuses() []domain.OrderStatus
	GetAllOrderDeliveryTypes() []domain.OrderDeliveryType
	GetOrderByID(id string) (domain.OrderByID, error)
	GetAllOrderDeliveryServicesCS() (*courierProto.ServicesResponse, error)
	CheckNewOrdersMark(id string) bool
}

type Manager interface {
	GetRestaurantID(id int) (string, error)
}

type Service struct {
	Restaurant
	Dish
	File
	Order
	Manager
}

func NewService(repos *repository.Repository, StorageProvider storage.Provider, cfg *config.Config) *Service {
	return &Service{
		Restaurant: NewRestaurantService(repos.Restaurant),
		Dish:       NewDishService(repos.Dish),
		File:       NewFileService(repos.File, StorageProvider),
		Order:      NewOrderService(repos.Order, cfg),
		Manager:    NewManagerService(repos.Manager),
	}
}
