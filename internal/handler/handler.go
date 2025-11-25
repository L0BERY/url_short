package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"shorturl.com/internal/service"
	"shorturl.com/pkg/config"
)

type Handler struct {
	service *service.ShortenerService
	config  *config.Config
}

func NewHandler(service *service.ShortenerService, config *config.Config) *Handler {
	return &Handler{
		service: service,
		config:  config,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	router.LoadHTMLGlob("web/templates/*")
	router.Static("/static", "./web/static")

	router.GET("/", h.home)
	router.GET("/:shortCode", h.redirect)
	// router.GET("/:shortCode/stats", h.stats)

	router.POST("/shorten", h.shorten)

	return router
}

func (h *Handler) home(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "URL shortener",
	})
}

func (h *Handler) redirect(c *gin.Context) {
	shortCode := c.Param("shortCode")
	originalURL, err := h.service.GetURL(shortCode)
	if err != nil {
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"message": "URL not found",
		})
		return
	}
	c.Redirect(http.StatusMovedPermanently, originalURL)
}

// func (h *Handler) stats(c *gin.Context) {
// 	shortCode := c.Param("shortCode")

// 	clickCount, createdAt, err := h.service.GetStats(shortCode)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"error": "URL not found",
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"short_code":       shortCode,
// 		"click_count":      clickCount,
// 		"created_at":       createdAt.Format(time.RFC3339),
// 		"created_at_human": createdAt.Format("02.01.2006 15:04"),
// 	})
// }

func (h *Handler) shorten(c *gin.Context) {
	var request struct {
		URL string `json:"url" binding:"required,url"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid URL",
		})
		return
	}

	shortCode, err := h.service.SaveURL(request.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "server error",
		})
		return
	}

	shortURL := h.config.BaseURL + "/" + shortCode
	c.JSON(http.StatusOK, gin.H{
		"short_url":  shortURL,
		"short_code": shortCode,
	})

}
