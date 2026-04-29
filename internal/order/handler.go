package order

import "github.com/gin-gonic/gin"

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

// POST /order (สร้างคำสั่งซื้อ)
func (h *Handler) CreateOrder(c *gin.Context) {}

// PUT /order/cancel (ลูกค้ายกเลิกออเดอร์)
func (h *Handler) CancelOrder(c *gin.Context) {}

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