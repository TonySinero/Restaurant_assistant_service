package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"restaurant-assistant/internal/domain"
)

// CreateDish @Summary Manager Create Dish
// @Tags dish
// @Description create dish
// @ModuleID CreateDish
// @Accept  json
// @Produce  json
// @Param input body domain.Dish true "create dish info"
// @Success 200 {object} dataResponse "ok"
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /dish/ [post]
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func (h *Handler) CreateDish(c *gin.Context) {
	restaurantID:= c.GetString("restaurantID")
	var input domain.Dish
	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	id, err := h.services.Dish.CreateDish(input, restaurantID)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dataResponse{
		Data: id,
	})
}

// UpdateDish @Summary Manager Update Dish
// @Tags dish
// @Description update dish
// @ModuleID UpdateDish
// @Accept  json
// @Produce  json
// @Param id path string true "dish id"
// @Param input body domain.UpdateDish true "order update info"
// @Success 200 {object} getDishResponse "ok"
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /dish/{id} [put]
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func (h *Handler) UpdateDish(c *gin.Context) {

	var input domain.UpdateDish

	id := c.Param("id")

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	dish, err := h.services.Dish.UpdateDish(id, input)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getDishResponse{
		Data: dish,
	})
}

// GetAllDishes @Summary Manager Get All dishes
// @Tags dish
// @Description Get all dishes
// @ModuleID GetAllDishes
// @Accept  json
// @Produce  json
// @Success 200 {object} getAllDishesResponse "ok"
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /dish/all [get]
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func (h *Handler) GetAllDishes(c *gin.Context) {
	restaurantID:= c.GetString("restaurantID")
	dishes, err := h.services.Dish.GetAllDishes(restaurantID)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllDishesResponse{
		Data: dishes,
	})
}

// DeleteDish @Summary Manager delete dish
// @Tags dish
// @Description delete dish
// @ModuleID DeleteDish
// @Accept  json
// @Produce  json
// @Param id path string true "dish id"
// @Success 200 {string} string "ok"
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /dish/{id} [delete]
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func (h *Handler) DeleteDish(c *gin.Context) {

	id := c.Param("id")

	err := h.services.Dish.DeleteDish(id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, "The dish has been successfully deleted")
}

// GetDishByID @Summary Manager get dish information
// @Tags dish
// @Description get dish information
// @ModuleID GetDishByID
// @Accept  json
// @Produce  json
// @Param id path string true "dish id"
// @Success 200 {object} getDishResponse "ok"
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /dish/{id} [get]
func (h *Handler) GetDishByID(c *gin.Context) {

	id := c.Param("id")

	dish, err := h.services.Dish.GetDishByID(id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getDishResponse{
		Data: dish,
	})
}

// GetDishByRestaurantID @Summary Manager get all restaurant dishes
// @Tags dish
// @Description Get all restaurant dishes
// @ModuleID GetDishByRestaurantID
// @Accept  json
// @Produce  json
// @Param id path string true "restaurant id"
// @Success 200 {object} getAllDishesResponse "ok"
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /dish/restaurant/{id} [get]
func (h *Handler) GetDishByRestaurantID(c *gin.Context) {

	id := c.Param("id")

	dishes, err := h.services.Dish.GetDishByRestaurantID(id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllDishesResponse{
		Data: dishes,
	})
}

// GetDishWithCategoryByRestaurantID @Summary Manager get all restaurant dishes with category
// @Tags dish
// @Description Get all restaurant dishes with category
// @ModuleID GetDishWithCategoryByRestaurantID
// @Accept  json
// @Produce  json
// @Param id path string true "restaurant id"
// @Success 200 {object} getDishesByRestaurantResponse "ok"
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /restaurant/dishes/{id} [get]
func (h *Handler) GetDishWithCategoryByRestaurantID(c *gin.Context) {
	id := c.Param("id")

	dishes, err := h.services.Dish.GetDishWithCategoryByRestaurantID(id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getDishesByRestaurantResponse{
		Data: dishes,
	})
}

// GetDishesTypes @Summary Manager get all types of dishes
// @Tags dish
// @Description Get all types of dishes
// @ModuleID GetDishesTypes
// @Accept  json
// @Produce  json
// @Success 200 {object} getDishesTypes "ok"
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /dish/types [get]
func (h *Handler) GetDishesTypes(c *gin.Context) {
	types, err := h.services.Dish.GetDishesTypes()
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getDishesTypes{
		Data: types,
	})
}
