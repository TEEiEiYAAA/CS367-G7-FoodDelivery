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

// @Summary  เพิ่มเมนูในร้าน
// @Tags     Menu
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    id path int true "Restaurant ID"
// @Param    body body CreateMenuRequest true "body"
// @Success  201 {object} Menu
// @Router   /restaurant/{id}/menu [post]
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

// @Summary ดูเมนูของร้าน
// @Tags    Menu
// @Produce json
// @Param   id path int true "Restaurant ID"
// @Success 200 {array} Menu
// @Router  /restaurant/{id}/menu [get]
func (h *Handler) GetMenu(c *gin.Context) {
	restaurantID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid restaurant id"})
		return
	}

	menus, err := h.service.GetMenu(restaurantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get menu"})
		return
	}

	c.JSON(http.StatusOK, menus)
}
