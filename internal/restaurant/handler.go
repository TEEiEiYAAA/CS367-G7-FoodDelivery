package restaurant

import "github.com/gin-gonic/gin"

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

// POST /restaurant (สร้างร้านอาหาร)
func (h *Handler) CreateRestaurant(c *gin.Context) {}

// GET /restaurant (ดูร้านทั้งหมด)
func (h *Handler) GetRestaurants(c *gin.Context) {}

// GET /restaurant/{id} (ดูข้อมูลของร้านอาหาร)
func (h *Handler) GetRestaurantByID(c *gin.Context) {}

// PUT /restaurant/order/confirm (ยืนยันออเดอร์)
func (h *Handler) ConfirmOrder(c *gin.Context) {}
