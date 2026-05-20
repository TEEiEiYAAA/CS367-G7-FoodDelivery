package order

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

// @Summary  สร้างคำสั่งซื้อ
// @Tags     Order
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    body body CreateOrderRequest true "body"
// @Success  201 {object} CreateOrderResponse
// @Router   /order [post]
func (h *Handler) CreateOrder(c *gin.Context) {
	// ดึง username จาก JWT ที่ AuthMiddleware set ไว้ใน context
	usernameVal, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	username, ok := usernameVal.(string)
	if !ok || username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
		return
	}

	// Bind + Validate request body
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// เรียก service layer
	orderID, totalPrice, err := h.service.CreateOrder(username, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, CreateOrderResponse{
		ID:         orderID,
		TotalPrice: totalPrice,
		Status:     "pending",
	})
}

// @Summary  ลูกค้ายกเลิกออเดอร์
// @Tags     Order
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    body body CancelOrderRequest true "body"
// @Success  200 {object} map[string]string
// @Router   /order/cancel [put]
func (h *Handler) CancelOrder(c *gin.Context) {
	// ดึง username จาก JWT
	usernameVal, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	username, ok := usernameVal.(string)
	if !ok || username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
		return
	}

	// Bind request body
	var req CancelOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// เรียก service
	if err := h.service.CancelOrder(username, req.OrderID); err != nil {
		switch {
		case containsAny(err.Error(), "not found"):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case containsAny(err.Error(), "forbidden"):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case containsAny(err.Error(), "cannot be cancelled", "expired"):
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order cancelled successfully"})
}

// @Summary ดูรายละเอียดออเดอร์
// @Tags    Order
// @Produce json
// @Param   id path int true "Order ID"
// @Success 200 {object} map[string]interface{}
// @Router  /order/{id} [get]
func (h *Handler) GetOrderByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}
	order, items, err := h.service.GetOrderByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"order": order,
		"items": items,
	})
}

// @Summary  อัปเดตสถานะออเดอร์
// @Tags     Order
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    id path int true "Order ID"
// @Param    body body UpdateOrderStatusRequest true "body"
// @Success  200 {object} map[string]string
// @Router   /order/{id}/status [put]
func (h *Handler) UpdateOrderStatus(c *gin.Context) {
	// แปลง order ID จาก URL param
	idStr := c.Param("id")
	orderID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	// ดึง role จาก JWT
	roleVal, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	role, _ := roleVal.(string)

	// Bind request body
	var req UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// เรียก service
	if err := h.service.UpdateOrderStatus(orderID, req.Status, role); err != nil {
		switch {
		case containsAny(err.Error(), "not found"):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case containsAny(err.Error(), "forbidden"):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case containsAny(err.Error(), "invalid transition"):
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}

// @Summary  มอบหมายไรเดอร์ให้ออเดอร์
// @Tags     Rider
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    id path int true "Order ID"
// @Success  200 {object} map[string]string
// @Router   /order/{id}/assign-rider [post]
func (h *Handler) AssignRider(c *gin.Context) {
	orderID := c.Param("id")

	var body struct {
		RiderID int `json:"riderId"`
	}

	// รับค่า JSON
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	// เรียก service
	err := h.service.AssignRider(orderID, body.RiderID)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Rider assigned successfully",
	})
}

func containsAny(s string, substrings ...string) bool {
	for _, substr := range substrings {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}
