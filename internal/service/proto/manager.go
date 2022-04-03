package proto

import (
	"context"
	"restaurant-assistant/internal/repository/proto"
	"restaurant-assistant/pkg/managerservice"
)

type ManagerService struct {
	repo *proto.Repository
}

func NewManagerService(repo *proto.Repository) *ManagerService {
	return &ManagerService{repo:repo}
}

func (*ManagerService) CreateManager (ctx context.Context, manager *managerservice.Manager) (*managerservice.ManagerState, error) {
	ans := managerservice.ManagerState{Success: true}
	return &ans, nil
}