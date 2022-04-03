package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"restaurant-assistant/internal/domain"
	"restaurant-assistant/pkg/courierProto"
	"restaurant-assistant/pkg/validation"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultLimit  = 20
	DefaultOffset = 0
)

type GetAllOrderDeliveryServices struct {
	Data []*courierProto.DeliveryService `json:"data"`
}

type GetAllOrders struct {
	Data []domain.GetOrder `json:"data"`
}

type GetAllOrderStatuses struct {
	Data []domain.OrderStatus `json:"data"`
}

type GetOrderByID struct {
	Data []domain.OrderByID `json:"data"`
}

type GetAllOrderDeliveryTypes struct {
	Data []domain.OrderDeliveryType `json:"data"`
}

type CheckNewOrdersMark struct {
	Data bool `json:"data"`
}

// GetOrderByID @Summary Manager get order information
// @Security manager-auth
// @Tags order
// @Description get order information
// @ModuleID GetOrderByID
// @Accept  json
// @Produce  json
// @Param id path string true "order id"
// @Success 200 {object} GetOrderByID "ok"
// @Success 400 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /order/{id} [get]
func (h *Handler) GetOrderByID(c *gin.Context) {
	id := c.Param("id")

	if ok := validation.IsValidUUID(id); !ok {
		log.Error().Msg("invalid order uuid")
		newResponse(c, http.StatusBadRequest, "invalid order uuid")
		return
	}

	order, err := h.services.Order.GetOrderByID(id)

	if err == domain.ErrOrderNotFound {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "failed to get order")
		return
	}

	c.JSON(http.StatusOK, dataResponse{
		Data: order})

	return
}

// UpdateOrder @Summary Manager Update Order
// @Security manager-auth
// @Tags order
// @Description update order status
// @ModuleID UpdateOrder
// @Accept  json
// @Produce  json
// @Param input body domain.UpdateOrder true "order update info"
// @Success 200 {string} string "ok"
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /order [put]
func (h *Handler) UpdateOrder(c *gin.Context) {
	var input domain.UpdateOrder

	id := c.Param("id")

	if err := c.BindJSON(&input); err != nil {
		log.Error().Err(err)
		newResponse(c, http.StatusBadRequest, "wrong value of order uuid")
		return
	}

	err := h.services.Order.UpdateOrder(id, input)
	if err != nil {
		switch err {
		case domain.ErrCourServiceNotSelected:
			newResponse(c, http.StatusBadRequest, err.Error())
			return

		case domain.ErrDelTypeNotSelected:
			newResponse(c, http.StatusBadRequest, err.Error())
			return

		default:
			newResponse(c, http.StatusInternalServerError, "failed to update order in db")
			return
		}
	}

	c.JSON(http.StatusOK, nil)

	return
}

// GetAllOrders @Summary Manager Get Orders
// @Security manager-auth
// @Tags order
// @Description Get orders with filter, ordered by created time
// @ModuleID GetAllOrders
// @Accept  json
// @Produce  json
// @Param        start    query     string  false  "select orders with created date > start" Format(date)
// @Param        end    query     string  false  "select orders with created date < end" Format(date)
// @Param        dishName    query     string  false  "dish name"
// @Param        status    query     int  false  "order status" Enums(1, 2, 3, 4, 5, 6)
// @Param        limit    query     int  false  "number of returned orders" default(20)
// @Param        page    query     int  false  "page" default(1)
// @Success 200 {object} GetAllOrders "ok"
// @Failure default {object} dataResponse
// @Router /order [get]
func (h *Handler) GetAllOrders(c *gin.Context) {
	restaurantID:= c.GetString("restaurantID")
	//fmt.Println(c.Query("restaurantID"))
	if ok := validation.IsValidUUID(restaurantID); !ok {
		restaurantID = ""
	}

	dishName := c.Query("dishName")
	if dishName != "" {
		dishName = strings.ToUpper(dishName) + "%"
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit <= 0 {
		limit = DefaultLimit
	}

	var offset int
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
		offset = DefaultOffset
	} else {
		offset = (page - 1) * limit
	}

	updatedStatus, err := strconv.Atoi(c.Query("status"))
	if err != nil {
		updatedStatus = 0
	}

	updatedDateStart, err := time.Parse("2006-01-02", c.Query("start"))
	if err != nil {
		updatedDateStart = time.Time{}
	}

	updatedDateEnd, err := time.Parse("2006-01-02", c.Query("end"))
	if err != nil {
		updatedDateEnd = time.Time{}

	}

	filter := &domain.FilterOrder{
		RestaurantID: restaurantID,
		Status:       updatedStatus,
		DateStart:    updatedDateStart,
		DateEnd:      updatedDateEnd,
		DishName:     dishName,
	}

	allOrders, _ := h.services.GetAllOrders(filter, limit, offset)

	c.JSON(http.StatusOK, dataResponse{
		Data: allOrders,
	})

	return
}

// GetAllOrderStatuses @Summary Manager get all order status
// @Tags order
// @Description Get all order status
// @ModuleID GetAllOrderOrders
// @Accept  json
// @Produce  json
// @Success 200 {object} GetAllOrderStatuses "ok"
// @Router /order/status [get]
func (h *Handler) GetAllOrderStatuses(c *gin.Context) {

	allStatus := h.services.GetAllOrderStatuses()

	c.JSON(http.StatusOK, dataResponse{
		Data: allStatus,
	})

	return
}

// GetAllOrderDeliveryTypes @Summary Manager get all order delivery types
// @Tags order
// @Description get all order delivery types
// @ModuleID GetAllOrderDeliveryTypes
// @Accept  json
// @Produce  json
// @Success 200 {object} GetAllOrderDeliveryTypes "ok"
// @Router /order/delivery-type [get]
func (h *Handler) GetAllOrderDeliveryTypes(c *gin.Context) {

	allDeliveryTypes := h.services.GetAllOrderDeliveryTypes()

	c.JSON(http.StatusOK, dataResponse{
		Data: allDeliveryTypes,
	})

	return
}

// GetAllOrderDeliveryServices @Summary Manager get all delivery services
// @Security manager-auth
// @Tags order
// @Description Get all delivery services
// @ModuleID GetAllOrderDeliveryServices
// @Accept  json
// @Produce  json
// @Success 200 {object} GetAllOrderDeliveryServices "ok"
// @Router /order/delivery-service [get]
// @Failure 500 {object} response
// @Failure default {object} response
func (h *Handler) GetAllOrderDeliveryServices(c *gin.Context) {

	allDeliveryServices, err := h.services.GetAllOrderDeliveryServicesCS()

	if err != nil {
		newResponse(c, http.StatusInternalServerError, "failed to get all order delivery services")
		return
	}

	c.JSON(http.StatusOK, GetAllOrderDeliveryServices{
		Data: allDeliveryServices.Services,
	})

	return
}



// CheckNewOrdersMark @Summary Manager can accept status of presence of new orders
// @Security manager-auth
// @Tags order
// @Description Get status of presence of active orders
// @ModuleID CheckNewOrdersMark
// @Accept  json
// @Produce  json
// @Success 200 {object} CheckNewOrdersMark "ok"
// @Router /order/check-new-orders-mark [get]
func (h *Handler) CheckNewOrdersMark(c *gin.Context) {
	restaurantID:= c.GetString("restaurantID")
	//id := c.Param("id")

	status := h.services.CheckNewOrdersMark(restaurantID)

	c.JSON(http.StatusOK, dataResponse{
		Data: status,
	})

	return
}

// CheckNewOrdersMarkAuth @Summary Manager can accept status of presence of new orders
// @Security manager-auth
// @Tags order
// @Description Get status of presence of active orders
// @ModuleID CheckNewOrdersMarkAuth
// @Accept  json
// @Produce  json
// @Success 200 {object} CheckNewOrdersMark "ok"
// @Router /authenticated/order/check-new-orders-mark [get]
func (h *Handler) CheckNewOrdersMarkAuth(c *gin.Context) {

	id := c.GetString("restaurantID")

	status := h.services.CheckNewOrdersMark(id)

	c.JSON(http.StatusOK, dataResponse{
		Data: status,
	})

	return
}