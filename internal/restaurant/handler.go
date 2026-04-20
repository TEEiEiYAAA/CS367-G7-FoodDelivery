package restaurant

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) CreateRestaurant(c *gin.Context) {
	var req CreateRestaurantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username := c.GetString("username")
	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "username not found in token"})
		return
	}

	rest := Restaurant{
		Name:          req.Name,
		Address:       req.Address,
		OwnerUsername: username,
	}

	created, err := h.service.CreateRestaurant(rest)
	if err != nil {
		if errors.Is(err, ErrInvalidRestaurant) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create restaurant"})
		return
	}

	c.JSON(http.StatusCreated, created)
}

func (h *Handler) GetRestaurants(c *gin.Context) {
	list, err := h.service.GetRestaurants()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch restaurants"})
		return
	}
	c.JSON(http.StatusOK, list)
}

// GET /restaurant/{id} (ดูข้อมูลของร้านอาหาร)
func (h *Handler) GetRestaurantByID(c *gin.Context) {}

// PUT /restaurant/order/confirm (ยืนยันออเดอร์)
func (h *Handler) ConfirmOrder(c *gin.Context) {}
