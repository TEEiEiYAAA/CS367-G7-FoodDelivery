package rider

import "github.com/gin-gonic/gin"

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

// POST /order/{id}/assign-rider (มอบหมายไรเดอร์)
func (h *Handler) AssignRider(c *gin.Context) {}
