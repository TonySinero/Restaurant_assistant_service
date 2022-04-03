package service

import (
	"restaurant-assistant/internal/domain"
	"restaurant-assistant/internal/repository"
)

type DishService struct {
	repo repository.Dish
}

func NewDishService(repo repository.Dish) *DishService {
	return &DishService{repo: repo}
}

func (s *DishService) CreateDish(input domain.Dish, id string) (string, error) {
	return s.repo.CreateDish(input, id)
}

func (s *DishService) UpdateDish(id string, input domain.UpdateDish) (domain.GetDishByID, error) {
	err := s.repo.UpdateDish(id, input)

	dish, _ := s.repo.GetDishByID(id)

	return dish, err
}

func (s *DishService) GetAllDishes(id string) ([]domain.GetAllDishes, error) {
	return s.repo.GetAllDishes(id)
}

func (s *DishService) DeleteDish(id string) error {
	return s.repo.DeleteDish(id)
}

func (s *DishService) GetDishByID(id string) (domain.GetDishByID, error) {
	return s.repo.GetDishByID(id)
}

func (s *DishService) GetDishByRestaurantID(id string) ([]domain.GetAllDishes, error) {
	return s.repo.GetDishByRestaurantID(id)
}

func (s *DishService) GetDishWithCategoryByRestaurantID(id string) ([]domain.GetDishesByRestaurant, error) {
	return s.repo.GetDishWithCategoryByRestaurantID(id)
}

func (s *DishService) GetDishesTypes() ([]domain.DishesCategory, error) {
	return s.repo.GetDishesTypes()
}