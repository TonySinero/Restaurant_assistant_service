package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"restaurant-assistant/internal/domain"
)

// CreateRestaurant @Summary Super Admin Create Restaurant
// @Security superadmin-auth
// @Tags restaurant
// @Description create restaurant
// @ModuleID CreateRestaurant
// @Accept  json
// @Produce  json
// @Param input body domain.Restaurant true "create restaurant"
// @Success 200 {object} dataResponse "ok"
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /restaurant [post]
func (h *Handler) CreateRestaurant(c *gin.Context) {
	var input domain.Restaurant

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	id, err := h.services.Restaurant.CreateRestaurant(input)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "failed to create restaurant")
		return
	}

	c.JSON(http.StatusOK, dataResponse{
		Data: id,
	})
}

// UpdateRestaurant @Summary Manager Update Restaurant
// @Security superadmin-auth
// @Tags restaurant
// @Description update restaurant
// @ModuleID UpdateRestaurant
// @Accept  json
// @Produce  json
// @Param id path string true "restaurant id"
// @Param input body domain.UpdateRestaurant true "create restaurant"
// @Success 200 {object} getRestaurant "ok"
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /restaurant/{id} [put]
func (h *Handler) UpdateRestaurant(c *gin.Context) {

	var input domain.UpdateRestaurant

	id := c.Param("id")

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	restaurant, err := h.services.Restaurant.UpdateRestaurant(id, input)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getRestaurant{
		Data: restaurant,
	})
}

// GetAllRestaurant @Summary Getting restaurants
// @Tags restaurant
// @Description getting restaurants
// @ModuleID GetAllRestaurant
// @Produce  json
// @Param sort_by body domain.GetRestaurantOrderBy string "sort restaurant by rating and title"
// @Success 200 {object} getAllRestaurantsResponse "ok"
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /restaurant/ [get]
func (h *Handler) GetAllRestaurant(c *gin.Context) {

	var input domain.GetRestaurantOrderBy

	if err := c.Bind(&input); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	var restaurants []domain.GetRestaurant
	var err error

	if input.Category != "" {
		restaurants, err = h.services.Restaurant.GetRestaurantsByCategory(input)
		if err != nil {
			newResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		restaurants, err = h.services.Restaurant.GetAllRestaurant(input)
		if err != nil {
			newResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, getAllRestaurantsResponse{
		Data: restaurants,
	})
}

// GetRestaurantCategoriesWithRestaurants @Summary Getting Restaurant Categories With Restaurants
// @Tags restaurant
// @Description getting restaurant categories with restaurants
// @ModuleID GetRestaurantCategoriesWithRestaurants
// @Produce  json
// @Success 200 {object} getRestaurantsCategoriesResponse "ok"
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /restaurant/categories-with-restaurants [get]
func (h *Handler) GetRestaurantCategoriesWithRestaurants(c *gin.Context) {
	categories, err := h.services.Restaurant.GetRestaurantCategoriesWithRestaurants()
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getRestaurantsCategoriesResponse{
		Data: categories,
	})
}

// GetRestaurantCategories @Summary Getting Restaurant Categories
// @Tags restaurant
// @Description getting restaurant categories
// @ModuleID GetRestaurantCategories
// @Produce  json
// @Success 200 {object} getCategories "ok"
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /restaurant/categories [get]
func (h *Handler) GetRestaurantCategories(c *gin.Context) {
	categories, err := h.services.Restaurant.GetRestaurantCategories()
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getCategories{
		Data: categories})
}

// DeleteRestaurant @Summary Delete Restaurant
// @Security superadmin-auth
// @Tags restaurant
// @Description delete restaurant
// @ModuleID DeleteRestaurant
// @Produce  json
// @Param id path string true "restaurant id"
// @Success 200 {object} string "ok"
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /restaurant/{id} [delete]
func (h *Handler) DeleteRestaurant(c *gin.Context) {

	id := c.Param("id")

	err := h.services.Restaurant.DeleteRestaurant(id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "The restaurant has been successfully deleted")
}

// GetRestaurantById @Summary Getting Restaurant By Id
// @Tags restaurant
// @Description getting restaurant by id
// @ModuleID GetRestaurantById
// @Param id path string true "restaurant id"
// @Produce  json
// @Success 200 {object} getRestaurant "ok"
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /restaurant/{id} [get]
func (h *Handler) GetRestaurantById(c *gin.Context) {

	id := c.Param("id")

	restaurant, err := h.services.Restaurant.GetRestaurantById(id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getRestaurant{
		Data: restaurant},
	)
}

// GetNearestRestaurant @Summary Getting restaurants
// @Tags restaurant
// @Description getting restaurants
// @ModuleID GetNearestRestaurant
// @Produce  json
// @Param address body domain.GetRestaurantByAddress string "get nearest restaurants"
// @Success 200 {object} getAllRestaurantsResponse "ok"
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /restaurant/nearest [get]
func (h *Handler) GetNearestRestaurant(c *gin.Context) {

	var input domain.GetRestaurantByAddress

	if err := c.Bind(&input); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	restaurants, err := h.services.Restaurant.GetNearestRestaurant(input)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllRestaurantsResponse{
		Data: restaurants,
	})
}

// GetRestaurantFeedbacksById @Summary Getting Restaurant Feedbacks By ID
// @Security manager-auth
// @Tags restaurant
// @Description getting restaurant feedbacks by iD
// @ModuleID GetRestaurantFeedbacksById
// @Param id path string true "restaurant id"
// @Produce  json
// @Success 200 {object} getFeedbacks "ok"
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /restaurant/feedbacks [get]
func (h *Handler) GetRestaurantFeedbacksById(c *gin.Context) {
	id := c.GetString("restaurantID")

	restaurant, err := h.services.Restaurant.GetRestaurantFeedbacksById(id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getFeedbacks{
		Data: restaurant},
	)
}
