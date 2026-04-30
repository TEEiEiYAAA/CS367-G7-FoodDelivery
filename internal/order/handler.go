package order

import (
	"net/http"
	"strings"

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
		OrderID:    orderID,
		TotalPrice: totalPrice,
		Status:     "pending",
	})
}

// PUT /order/cancel (ลูกค้ายกเลิกออเดอร์)
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

// GET /order/{id} (ดูรายละเอียดออเดอร์)
func (h *Handler) GetOrderByID(c *gin.Context) {}

// PUT /order/{id}/status (อัปเดตสถานะออเดอร์ เช่น รับออเดอร์ กำลังทำ ทำเสร็จ กำลังจัดส่ง)
func (h *Handler) UpdateOrderStatus(c *gin.Context) {}

// POST /order/{id}/assign-rider (มอบหมายไรเดอร์)
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