package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikemeh/ecommerce-api/models"
	"gorm.io/gorm"
)

type OrderController struct {
	DB *gorm.DB
}

func NewOrderController(db *gorm.DB) *OrderController {
	return &OrderController{DB: db}
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order with items
// @Tags orders
// @Accept json
// @Produce json
// @Param order body models.Order true "Order info"
// @Success 201 {object} models.Order
// @Failure 400 {object} map[string]string
// @Security ApiKeyAuth
// @Router /orders [post]
func (oc *OrderController) CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")
	order.UserID = userID.(uint)
	order.Status = "Pending"

	if err := oc.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (oc *OrderController) GetOrders(c *gin.Context) {
	userID, _ := c.Get("userID")
	var orders []models.Order
	if err := oc.DB.Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetOrder godoc
// @Summary Get an order
// @Description Get an order by its ID
// @Tags orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} models.Order
// @Failure 404 {object} map[string]string
// @Security ApiKeyAuth
// @Router /orders/{id} [get]
func (oc *OrderController) GetOrder(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")

	var order models.Order
	if err := oc.DB.Where("id = ? AND user_id = ?", id, userID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (oc *OrderController) CancelOrder(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")

	var order models.Order
	if err := oc.DB.Where("id = ? AND user_id = ?", id, userID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	if order.Status != "Pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only pending orders can be cancelled"})
		return
	}

	order.Status = "Cancelled"
	if err := oc.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel order"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (oc *OrderController) UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")

	var statusUpdate struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order models.Order
	if err := oc.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	order.Status = statusUpdate.Status
	if err := oc.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	c.JSON(http.StatusOK, order)
}
