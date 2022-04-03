package repository

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"restaurant-assistant/internal/domain"
)

type ManagerPostgres struct {
	db *sqlx.DB
}

func NewManagerPostgres(db *sqlx.DB) *ManagerPostgres {
	return &ManagerPostgres{db: db}
}

func (r *ManagerPostgres) GetRestaurantID(managerID int) (string, error) {
	var restaurantID string
	SelectRestaurantID := fmt.Sprintf("SELECT id FROM restaurants WHERE user_id =$1")
	err := r.db.Get(&restaurantID, SelectRestaurantID, managerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", domain.ErrManagerNotAssigned
		}
		log.Error().Err(err).Msg("error occurred while selecting restaurant of manager")
		return "", err
	}
	return restaurantID, nil
}
