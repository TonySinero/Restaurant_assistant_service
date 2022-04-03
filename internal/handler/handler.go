package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "restaurant-assistant/docs"
	"restaurant-assistant/internal/service"
	"restaurant-assistant/pkg/authProto"
)

type Handler struct {
	services     *service.Service
	tokenManager *authProto.Manager
}

func NewHandler(services *service.Service, tokenManager *authProto.Manager) *Handler {
	return &Handler{services: services,
		tokenManager: tokenManager}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	router.Use(
		corsMiddleware,
	)

	restaurant := router.Group("/restaurant")
	{
		restaurant.POST("/",  h.CreateRestaurant)
		restaurant.PUT("/:id",h.UpdateRestaurant)
		restaurant.GET("/", h.GetAllRestaurant)
		restaurant.GET("/:id", h.GetRestaurantById)
		restaurant.GET("/categories-with-restaurants", h.GetRestaurantCategoriesWithRestaurants)
		restaurant.GET("/categories", h.GetRestaurantCategories)
		restaurant.DELETE("/:id",h.DeleteRestaurant)
		restaurant.POST("/image/:id", h.UploadImage)
		restaurant.GET("/nearest", h.GetNearestRestaurant)
		restaurant.GET("/dishes/:id", h.GetDishWithCategoryByRestaurantID)
		restaurant.GET("/feedbacks",h.GetRestaurantFeedbacksById)
	}

	order := router.Group("/order")
	{
		order.PUT("/:id", h.ManagerIdentity, h.UpdateOrder)
		order.GET("/status", h.GetAllOrderStatuses)
		order.GET("/", h.ManagerIdentity, h.GetAllOrders)
		order.GET("/check-new-orders-mark", h.ManagerIdentity, h.CheckNewOrdersMark)
		order.GET("/delivery-type",  h.GetAllOrderDeliveryTypes)
		order.GET("/:id", h.ManagerIdentity, h.GetOrderByID)
		order.GET("/delivery-service", h.ManagerIdentity, h.GetAllOrderDeliveryServices)
	}

	dish := router.Group("/dish")
	{
		dish.POST("/", h.ManagerIdentity, h.CreateDish)
		dish.PUT("/:id", h.ManagerIdentity, h.UpdateDish)
		dish.GET("/:id", h.GetDishByID)
		dish.DELETE("/:id", h.ManagerIdentity, h.DeleteDish)
		dish.POST("/image/:id", h.ManagerIdentity, h.UploadImage)
		dish.GET("/restaurant/:id", h.GetDishByRestaurantID)
		dish.GET("/all", h.ManagerIdentity, h.GetAllDishes)
		dish.GET("/types", h.GetDishesTypes)
	}

	authenticated := router.Group("/authenticated", h.ManagerIdentity)
	{
		authOrder := authenticated.Group("/order")
		{
			authOrder.GET("/check-new-orders-mark", h.CheckNewOrdersMarkAuth)
			authOrder.GET("", h.GetAllOrders)
		}
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
