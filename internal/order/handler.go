package order

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

// POST /order (สร้างคำสั่งซื้อ)
func (h *Handler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customerUsername := c.GetString("username")
	if customerUsername == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "username not found in token"})
		return
	}

	order, err := h.service.CreateOrder(req, customerUsername)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create order"})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// PUT /order/cancel (ลูกค้ายกเลิกออเดอร์)
func (h *Handler) CancelOrder(c *gin.Context) {
	var body struct {
		OrderID int `json:"order_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customerUsername := c.GetString("username")
	if customerUsername == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "username not found in token"})
		return
	}

	err := h.service.CancelOrder(body.OrderID, customerUsername)
	if err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
			return
		}
		if errors.Is(err, ErrGracePeriodExpired) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot cancel: grace period has expired"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not cancel order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order cancelled successfully"})
}

// GET /order/:id (ดูรายละเอียดออเดอร์)
func (h *Handler) GetOrderByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	order, err := h.service.GetOrderByID(id)
	if err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch order"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// PUT /order/:id/status (อัปเดตสถานะออเดอร์)
func (h *Handler) UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")

	var body struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateOrderStatus(orderID, body.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status updated successfully"})
}

// POST /order/:id/assign-rider (มอบหมายไรเดอร์)
func (h *Handler) AssignRider(c *gin.Context) {
	orderID := c.Param("id")

	var body struct {
		RiderID int `json:"riderId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AssignRider(orderID, body.RiderID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "rider assigned successfully"})
}