package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"os"
	"restaurant-assistant/internal/domain"
	"restaurant-assistant/internal/repository"
	"restaurant-assistant/pkg/storage"
	"strings"
)

type FileService struct {
	repo    repository.File
	storage storage.Provider
}

func NewFileService(repo repository.File, storage storage.Provider) *FileService {
	return &FileService{repo: repo, storage: storage}
}

func (s *FileService) Save(link string, uuid string, path string) error {
	return s.repo.Create(link, uuid, path)
}

func (s *FileService) CheckUUID(uuid string, path string) error {
	return s.repo.CheckUUID(uuid, path)
}

func (s *FileService) UploadAndSaveFile(ctx context.Context, file domain.File, uuid string, path string) (string, error) {
	defer removeFile(file.Name)

	if err := s.CheckUUID(uuid, path); err != nil {
		return "", err
	}

	link, err := s.upload(ctx, file)
	if err != nil {
		return "", err
	}

	if err := s.Save(link, uuid, path); err != nil {
		return "", err
	}

	return link, nil
}

func (s *FileService) upload(ctx context.Context, file domain.File) (string, error) {
	f, err := os.Open(file.Name)
	if err != nil {
		log.Error().Err(err)
		return "", err
	}

	defer f.Close()

	return s.storage.Upload(ctx, storage.UploadInput{
		File:        f,
		Size:        file.Size,
		ContentType: file.ContentType,
		Name:        s.generateFilename(file),
	})
}

func (s *FileService) generateFilename(file domain.File) string {
	filename := fmt.Sprintf("%s.%s", uuid.New().String(), getFileExtension(file.Name))
	fileNameParts := strings.Split(file.Name, "--") // first part is restaurantID

	return fmt.Sprintf("%s/%s", fileNameParts[0], filename)
}

func getFileExtension(filename string) string {
	parts := strings.Split(filename, ".")

	return parts[len(parts)-1]
}

func removeFile(filename string) {
	if err := os.Remove(filename); err != nil {
		log.Error().Err(err).Msg("removeFile():")
	}
}
