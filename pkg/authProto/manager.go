package authProto

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	config "restaurant-assistant/configs"
	"restaurant-assistant/internal/domain"
	"time"
)

type TokenManager interface {
	Parse(accessToken string) (*UserRole, error)
	Refresh(refreshToken string) (string, error)
}

type Manager struct {
	TokenManager
	cfg *config.Config
}

func NewManager(cfg *config.Config) (*Manager, error) {
	return &Manager{cfg: cfg}, nil
}

func (m *Manager) Parse(accessToken string) (*UserRole, error) {
	authClient, conn, ctx, err := m.CreateConnectionAuth(m.cfg)
	if err != nil {
		return nil, domain.ErrInternalServer
	}

	defer conn.Close()

	AccessToken := &AccessToken{
		AccessToken: accessToken,
	}

	user, err := authClient.GetUserWithRights(ctx, AccessToken)

	if err != nil {
		return nil, domain.ErrInternalServer
	}

	return user, nil
}

func (m *Manager) CreateConnectionAuth(cfg *config.Config) (AuthClient, *grpc.ClientConn, context.Context, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%s", cfg.GRPCAuth.Host, cfg.GRPCAuth.Port), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Error().Err(err).Msg("error occurred while creating conn to auth service")
		return nil, nil, ctx, err
	}

	authClient := NewAuthClient(conn)

	return authClient, conn, ctx, nil
}
