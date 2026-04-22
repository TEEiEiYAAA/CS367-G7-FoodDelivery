package restaurant

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

// POST /restaurant (สร้างร้านอาหาร)
func (h *Handler) CreateRestaurant(c *gin.Context) {}

// GET /restaurant (ดูร้านทั้งหมด)
func (h *Handler) GetRestaurants(c *gin.Context) {}

// GET /restaurant/{id} (ดูข้อมูลของร้านอาหาร)
func (h *Handler) GetRestaurantByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	rest, err := h.service.GetRestaurantByID(id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "restaurant not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, rest)
}

// PUT /restaurant/order/confirm (ยืนยันออเดอร์)
func (h *Handler) ConfirmOrder(c *gin.Context) {}
