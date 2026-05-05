package menu

import (
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

// POST /restaurant/{id}/menu (เพิ่มเมนู)
func (h *Handler) CreateMenu(c *gin.Context) {
	restaurantID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid restaurant id"})
		return
	}

	var req CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := h.service.CreateMenu(restaurantID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create menu item"})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// GET /restaurant/{id}/menu (ดูเมนู)
func (h *Handler) GetMenu(c *gin.Context) {}
