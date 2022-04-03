package server

import (
	"context"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net"
	"net/http"
	config "restaurant-assistant/configs"
	"restaurant-assistant/internal/service/proto"
	"restaurant-assistant/pkg/managerservice"
	"restaurant-assistant/pkg/orderservice_ra"
	"restaurant-assistant/pkg/restaurantservice"
)

type Server struct {
	server *http.Server
}

type ServerGRPC struct {
	server *grpc.Server
}

func NewServer(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		server: &http.Server{
			Addr:           ":" + cfg.HTTP.Port,
			Handler:        handler,
			ReadTimeout:    cfg.HTTP.ReadTimeout,
			WriteTimeout:   cfg.HTTP.WriteTimeout,
			MaxHeaderBytes: cfg.HTTP.MaxHeaderMegabytes << 20,
		},
	}
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func NewServerGRPC() *ServerGRPC {
	return &ServerGRPC{grpc.NewServer()}
}

func (s *ServerGRPC) RegisterServices(services *proto.Service) {
	orderservice_ra.RegisterOrderServiceServer(s.server, services)
	managerservice.RegisterManagerServiceServer(s.server, services)
	restaurantservice.RegisterRestaurantServiceServer(s.server, services)
}

func (s *ServerGRPC) Run(cfg *config.Config) error {

	lis, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
	if err != nil {
		log.Fatal().Err(err).Msg("error occurred while running grpc connection")

		return err
	}
	if err := s.server.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("error occurred while running grpc server")
		return err
	}
	return nil
}

func (s *ServerGRPC) Shutdown() {
	s.server.Stop()
}
