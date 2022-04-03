package main

import (
	"context"
	_ "github.com/golang-migrate/migrate"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	config "restaurant-assistant/configs"
	"restaurant-assistant/internal/handler"
	"restaurant-assistant/internal/repository"
	protorepository "restaurant-assistant/internal/repository/proto"
	"restaurant-assistant/internal/service"
	"restaurant-assistant/internal/service/proto"
	"restaurant-assistant/pkg/authProto"
	"restaurant-assistant/pkg/storage"
	"restaurant-assistant/server"
	"syscall"
	"time"
)

// @title Restaurant-Assistant API
// @version 1.0
// @description REST API for Restaurant-Assistant service

// @host 165.232.68.67:8080
// @BasePath /

// @securityDefinitions.apikey Auth
// @in header
// @name Authorization

func main() {
	cfg, err := config.Init("./configs")
	if err != nil {
		log.Fatal().Err(err).Msg("wrong config variables")
	}

	db, err := newPostgresDB(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("err initializing DB")
	}

	storageProvider, err := newStorageProvider(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("wrong config variables")
	}

	tokenManager, err := authProto.NewManager(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("err initializing tokenManager")
	}

	repo := repository.NewRepository(db)
	services := service.NewService(repo, storageProvider, cfg)
	handlers := handler.NewHandler(services, tokenManager)
	srv := server.NewServer(cfg, handlers.InitRoutes())

	go func() {
		if err := srv.Run(); err != http.ErrServerClosed {
			log.Error().Err(err).Msg("error occurred while running http server")
		}
	}()

	repoGRPC := protorepository.NewRepository(db)
	servicesGRPC := proto.NewService(repoGRPC)
	srvGRPC := server.NewServerGRPC()
	srvGRPC.RegisterServices(servicesGRPC)

	//go func() {
	//	for {
	//		if err := repo.Restaurant.RestaurantActivityCheck(); err != nil {
	//			log.Error().Err(err).Msg("")
	//		}
	//		time.Sleep(time.Minute * 30)
	//	}
	//}()

	go func() {
		if err := srvGRPC.Run(cfg); err != nil {
			log.Error().Err(err).Msg("error occurred while running gRPC server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to stop server")
	}

	srvGRPC.Shutdown()

	if err := db.Close(); err != nil {
		log.Fatal().Err(err).Msg("failed to stop connection db")
	}

}

func newPostgresDB(cfg *config.Config) (*sqlx.DB, error) {
	return repository.NewPostgresDB(repository.Config{
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		Username: cfg.Postgres.Username,
		Password: cfg.Postgres.Password,
		DBName:   cfg.Postgres.Dbname,
		SSLMode:  cfg.Postgres.Sslmode,
	})
}

func newStorageProvider(cfg *config.Config) (storage.Provider, error) {
	client, err := minio.New(cfg.FileStorage.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.FileStorage.AccessKey, cfg.FileStorage.SecretKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, err
	}

	provider := storage.NewFileStorage(client, cfg.FileStorage.Bucket, cfg.FileStorage.Endpoint)

	return provider, nil
}
