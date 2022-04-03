package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"restaurant-assistant/internal/domain"
)

type dataResponse struct {
	Data interface{} `json:"data"`
}

type getDishResponse struct {
	Data domain.GetDishByID `json:"data"`
}

type getAllDishesResponse struct {
	Data []domain.GetAllDishes `json:"data"`
}

type getAllRestaurantsResponse struct {
	Data []domain.GetRestaurant `json:"data"`
}

type getRestaurantsCategoriesResponse struct {
	Data []domain.GetRestaurantCategories `json:"data"`
}

type getDishesByRestaurantResponse struct {
	Data []domain.GetDishesByRestaurant `json:"data"`
}

type getCategories struct {
	Data []string `json:"data"`
}

type getRestaurant struct {
	Data domain.GetRestaurant `json:"data"`
}

type getFeedbacks struct {
	Data []domain.Feedback `json:"data"`
}

type getDishesTypes struct {
	Data []domain.DishesCategory `json:"data"`
}

type response struct {
	Message string `json:"message"`
}

func newResponse(c *gin.Context, statusCode int, message string) {
	log.Error().Msg(message)
	c.AbortWithStatusJSON(statusCode, response{message})
}
