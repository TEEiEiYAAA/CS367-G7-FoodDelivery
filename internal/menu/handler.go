package menu

import "github.com/gin-gonic/gin"

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

// POST /restaurant/{id}/menu (เพิ่มเมนู)
func (h *Handler) CreateMenu(c *gin.Context) {}

// GET /restaurant/{id}/menu (ดูเมนู)
func (h *Handler) GetMenu(c *gin.Context) {}
