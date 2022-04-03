package service

import (
	"restaurant-assistant/internal/repository"
)

type ManagerService struct {
	repo repository.Manager
}

func NewManagerService(repo repository.Manager) *ManagerService {
	return &ManagerService{repo: repo}
}

func (s *ManagerService) GetRestaurantID(id int) (string, error) {
	return s.repo.GetRestaurantID(id)
}
