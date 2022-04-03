package repository

import (
	"github.com/jmoiron/sqlx"
	"restaurant-assistant/internal/domain"
)

type Restaurant interface {
	CreateRestaurant(input domain.Restaurant) (string, error)
	UpdateRestaurant(id string, input domain.UpdateRestaurant) error
	GetAllRestaurant(input domain.GetRestaurantOrderBy) ([]domain.GetRestaurant, error)
	DeleteRestaurant(id string) error
	GetRestaurantById(id string) (domain.GetRestaurant, error)
	GetNearestRestaurant(lat, lng float64) ([]domain.GetRestaurant, error)
	GetRestaurantCategoriesWithRestaurants() ([]domain.GetRestaurantCategories, error)
	GetRestaurantCategories() ([]string, error)
	RestaurantActivityCheck() error
	GetRestaurantsByCategory(input domain.GetRestaurantOrderBy) ([]domain.GetRestaurant, error)
	CheckRestaurantDuplicates(input domain.Restaurant) error
	GetRestaurantFeedbacksById(id string) ([]domain.Feedback, error)
}

type Dish interface {
	CreateDish(input domain.Dish, id string) (string, error)
	UpdateDish(id string, input domain.UpdateDish) error
	GetAllDishes(id string) ([]domain.GetAllDishes, error)
	DeleteDish(id string) error
	GetDishByID(id string) (domain.GetDishByID, error)
	GetDishByRestaurantID(id string) ([]domain.GetAllDishes, error)
	GetDishWithCategoryByRestaurantID(id string) ([]domain.GetDishesByRestaurant, error)
	GetDishesTypes() ([]domain.DishesCategory, error)
}

type File interface {
	Create(link string, uuid string, path string) error
	CheckUUID(uuid string, path string) error
}

type Order interface {
	UpdateOrder(id string, input domain.UpdateOrder) error
	GetAllOrders(filter *domain.FilterOrder, limit int, offset int) (*[]domain.GetOrder, error)
	GetAllOrderStatuses() []domain.OrderStatus
	GetAllOrderDeliveryTypes() []domain.OrderDeliveryType
	GetOrderByID(id string) (domain.OrderByID, error)
	GetRestaurantByID(id string) (domain.RestaurantToCourier, error)
	CheckNewOrdersMark(id string) bool
}

type Manager interface {
	GetRestaurantID(id int) (string, error)
}

type Repository struct {
	Restaurant
	Dish
	File
	Order
	Manager
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Restaurant: NewRestaurantPostgres(db),
		Dish:       NewDishPostgres(db),
		File:       NewFilePostgres(db),
		Order:      NewOrderPostgres(db),
		Manager:    NewManagerPostgres(db),
	}
}
