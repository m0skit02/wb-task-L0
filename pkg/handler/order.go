package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"wb-task-L0/pkg/models"
)

// @Summary Create order
// @Tags orders
// @Description create order
// @ID create-order
// @Accept  json
// @Produce  json
// @Param input body models.Order true "order info"
// @Success 200 {object} map[string]string
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /orders [post]
func (h *Handler) createOrder(c *gin.Context) {
	var input models.Order
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	order, err := h.services.Order.Create(&input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": order.OrderUID,
	})
}

func (h *Handler) createOrderWithAssociations(c *gin.Context) {
	var input models.Order
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// сохраняем заказ вместе с ассоциациями
	if err := h.services.Order.CreateOrderWithAssociations(c.Request.Context(), &input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// берём OrderUID из input, т.к. метод возвращает только error
	c.JSON(http.StatusOK, gin.H{
		"id": input.OrderUID,
	})
}

type getAllOrdersResponse struct {
	Data []models.Order `json:"data"`
}

// @Summary Get All Orders
// @Tags orders
// @Description get all orders
// @ID get-all-orders
// @Accept  json
// @Produce  json
// @Success 200 {object} getAllOrdersResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /orders [get]
func (h *Handler) getAllOrders(c *gin.Context) {
	orders, err := h.services.Order.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllOrdersResponse{
		Data: orders,
	})
}

// @Summary Get Order By Id
// @Tags orders
// @Description get order by id
// @ID get-order-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "order id"
// @Success 200 {object} models.Order
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /orders/{id} [get]
func (h *Handler) getOrderById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	order, err := h.services.Order.GetByID(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, order)
}

// @Summary Delete Order
// @Tags orders
// @Description delete order by id
// @ID delete-order
// @Accept  json
// @Produce  json
// @Param id path string true "order id"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /orders/{id} [delete]
func (h *Handler) deleteOrder(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	err := h.services.Order.Delete(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
