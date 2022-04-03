package proto

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"restaurant-assistant/internal/domain"
)

type RestaurantPostgres struct {
	db *sqlx.DB
}

func NewRestaurantPostgres(db *sqlx.DB) *RestaurantPostgres {
	return &RestaurantPostgres{db: db}
}

func (r *RestaurantPostgres) GetRestaurantsInfo(ctx context.Context, lat, lng float64) ([]domain.GetRestaurant, error) {
	var restaurants []domain.GetRestaurant

	query := fmt.Sprintf(`SELECT id, title, description, rating, time_work_start, time_work_end, medium_price, 
		address, is_active, image FROM restaurants 
		ORDER BY earth_distance(ll_to_earth(latitude, longitude),ll_to_earth($1, $2))`)
	err := r.db.Select(&restaurants, query, lat, lng)
	if err != nil {
		log.Error().Err(err).Msg("Err proto restaurant db")
	}

	return restaurants, err
}
