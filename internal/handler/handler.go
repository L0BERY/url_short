package handler

import (
	"github.com/gin-gonic/gin"
	"shorturl.com/internal/service"
)

type Handler struct {
	service *service.ShortenerService
}

func NewHandler(service *service.ShortenerService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	router.GET("/", nil)
	router.GET(":/shortCode", nil)

	router.POST("/shorten", nil)

	return router
}
