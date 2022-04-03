package service

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"googlemaps.github.io/maps"
	"restaurant-assistant/internal/domain"
	"restaurant-assistant/internal/repository"
)

const googleKey = "AIzaSyD0i1yTpCTID4kOBJimXbh0LxiHiGbFGkQ"

type RestaurantService struct {
	repo repository.Restaurant
}

func NewRestaurantService(repo repository.Restaurant) *RestaurantService {
	return &RestaurantService{repo: repo}
}

func (s *RestaurantService) CreateRestaurant(input domain.Restaurant) (string, error) {

	if err := s.repo.CheckRestaurantDuplicates(input); err != nil {
		log.Error().Err(err).Msg("restaurant with this number/email already exist")
		return "", err
	}

	var clinetGCM *maps.Client
	if clinetGCM == nil {
		var err error
		clinetGCM, err = maps.NewClient(maps.WithAPIKey(googleKey))
		if err != nil {
			log.Error().Err(err).Msg("maps.NewClient failed: ")
			return "", err
		}
	}
	r := &maps.GeocodingRequest{
		Address: input.Address,
	}
	res, err := clinetGCM.Geocode(context.Background(), r)
	if err != nil || len(res) == 0 {
		log.Error().Err(err).Msg("didn`t find this address")
		return "", err
	}
	input.Latitude = res[0].Geometry.Location.Lat
	input.Longitude = res[0].Geometry.Location.Lng
	return s.repo.CreateRestaurant(input)
}

func (s *RestaurantService) UpdateRestaurant(id string, input domain.UpdateRestaurant) (domain.GetRestaurant, error) {
	if input.Address != nil {
		var clinetGCM *maps.Client
		if clinetGCM == nil {
			var err error
			clinetGCM, err = maps.NewClient(maps.WithAPIKey(googleKey))
			if err != nil {
				log.Error().Err(err).Msg("maps.NewClient failed: ")
			}
		}
		r := &maps.GeocodingRequest{
			Address: *input.Address,
		}
		res, err := clinetGCM.Geocode(context.Background(), r)
		if err != nil || len(res) == 0 {
			log.Error().Err(err).Msg("didn`t find this address")
			return domain.GetRestaurant{}, err
		}
		input.Latitude = res[0].Geometry.Location.Lat
		input.Longitude = res[0].Geometry.Location.Lng
	}

	s.repo.UpdateRestaurant(id, input)

	return s.repo.GetRestaurantById(id)
}

func (s *RestaurantService) GetAllRestaurant(input domain.GetRestaurantOrderBy) ([]domain.GetRestaurant, error) {
	return s.repo.GetAllRestaurant(input)
}

func (s *RestaurantService) GetRestaurantsByCategory(input domain.GetRestaurantOrderBy) ([]domain.GetRestaurant, error) {
	return s.repo.GetRestaurantsByCategory(input)
}

func (s *RestaurantService) DeleteRestaurant(id string) error {
	return s.repo.DeleteRestaurant(id)
}

func (s *RestaurantService) GetRestaurantById(id string) (domain.GetRestaurant, error) {
	return s.repo.GetRestaurantById(id)
}

func (s *RestaurantService) GetNearestRestaurant(input domain.GetRestaurantByAddress) ([]domain.GetRestaurant, error) {
	var clinetGCM *maps.Client
	if clinetGCM == nil {
		var err error
		clinetGCM, err = maps.NewClient(maps.WithAPIKey(googleKey))
		if err != nil {
			log.Error().Err(err).Msg("maps.NewClient failed: ")
		}
	}
	r := &maps.GeocodingRequest{
		Address: input.Address,
	}
	res, err := clinetGCM.Geocode(context.Background(), r)
	if err != nil || len(res) == 0 {
		return nil, fmt.Errorf("didn`t find this address")
	}
	lat := res[0].Geometry.Location.Lat
	lng := res[0].Geometry.Location.Lng
	return s.repo.GetNearestRestaurant(lat, lng)
}

func (s *RestaurantService) GetRestaurantCategoriesWithRestaurants() ([]domain.GetRestaurantCategories, error) {
	return s.repo.GetRestaurantCategoriesWithRestaurants()
}

func (s *RestaurantService) GetRestaurantCategories() ([]string, error) {
	return s.repo.GetRestaurantCategories()
}

func (s *RestaurantService) GetRestaurantFeedbacksById(id string) ([]domain.Feedback, error){
	return s.repo.GetRestaurantFeedbacksById(id)
}
