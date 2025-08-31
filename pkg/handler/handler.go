package handler

import (
	"github.com/gin-gonic/gin"
	"wb-task-L0/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/api")
	{
		orders := api.Group("/orders")
		{
			orders.POST("/", h.createOrder)
			orders.GET("/", h.getAllOrders)
			orders.GET("/:id", h.getOrderById)
			orders.DELETE("/:id", h.deleteOrder)
		}
	}

	return router
}
