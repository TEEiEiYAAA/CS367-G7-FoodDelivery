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
