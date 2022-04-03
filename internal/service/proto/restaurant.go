package proto

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"googlemaps.github.io/maps"
	"restaurant-assistant/internal/repository/proto"
	"restaurant-assistant/pkg/restaurantservice"
)

const googleKey = "AIzaSyD0i1yTpCTID4kOBJimXbh0LxiHiGbFGkQ"

type RestaurantService struct {
	repo *proto.Repository
}

func NewRestaurantService(repo *proto.Repository) *RestaurantService {
	return &RestaurantService{repo: repo}
}

func (s *RestaurantService) GetUserAddress(ctx context.Context, address *restaurantservice.UserAddress) (*restaurantservice.NearestRestaurants, error) {
	logrus.Info(fmt.Sprintf("get lng,lat by address %s", address.Address))
	var clinetGCM *maps.Client
	if clinetGCM == nil {
		var err error
		clinetGCM, err = maps.NewClient(maps.WithAPIKey(googleKey))
		if err != nil {
			logrus.Info("maps.NewClient failed", err)
		}
	}
	add := address.Address
	logrus.Info("this adddress" + add)
	fmt.Println(address.Address)
	r := &maps.GeocodingRequest{
		Address: add,
	}
	res, err := clinetGCM.Geocode(context.Background(), r)
	if err != nil || len(res) == 0 {
		return nil, fmt.Errorf("didn`t find this address")
	}
	lat := res[0].Geometry.Location.Lat
	logrus.Info(lat)
	lng := res[0].Geometry.Location.Lng
	logrus.Info(lng)
	restaurants, err := s.repo.GetRestaurantsInfo(ctx, lat, lng)
	if err != nil {
		return nil, err
	}

	protoRestaurants := make([]*restaurantservice.NearestRestaurant, 0, len(restaurants))
	for _, val := range restaurants {

		protoRestaurants = append(protoRestaurants,
			&restaurantservice.NearestRestaurant{ID: val.ID, Title: val.Title, Description: "*val.Description",
				Rating: *val.Rating, MediumPrice: *val.MediumPrice, TimeWorkStart: timestamppb.New(val.TimeWorkStart.UTC()),
				TimeWorkEnd: timestamppb.New(val.TimeWorkEnd.UTC()), Address: *val.Address, IsActive: *val.IsActive,
				Image: *val.Image,
			})
	}


	return &restaurantservice.NearestRestaurants{List: protoRestaurants}, err
}
