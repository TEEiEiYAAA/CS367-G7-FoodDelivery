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

// @Summary  สร้างร้านอาหาร
// @Tags     Restaurant
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    body body CreateRestaurantRequest true "body"
// @Success  201 {object} Restaurant
// @Router   /restaurant [post]
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

// @Summary ดูร้านอาหารทั้งหมด
// @Tags    Restaurant
// @Produce json
// @Success 200 {array} Restaurant
// @Router  /restaurant [get]
func (h *Handler) GetRestaurants(c *gin.Context) {
	list, err := h.service.GetRestaurants()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch restaurants"})
		return
	}
	c.JSON(http.StatusOK, list)
}

// @Summary ดูข้อมูลร้านอาหาร
// @Tags    Restaurant
// @Produce json
// @Param   id path int true "Restaurant ID"
// @Success 200 {object} Restaurant
// @Router  /restaurant/{id} [get]
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

// @Summary  ยืนยันออเดอร์ (ฝั่งร้าน)
// @Tags     Restaurant
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    body body ConfirmOrderRequest true "body"
// @Success  200 {object} map[string]string
// @Router   /restaurant/order/confirm [put]
func (h *Handler) ConfirmOrder(c *gin.Context) {
	var req ConfirmOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	ownerUsername := c.GetString("username")
	if ownerUsername == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "username not found in token"})
		return
	}

	if err := h.service.ConfirmOrder(req.OrderID, ownerUsername); err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found or not authorized"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order confirmed"})
}
