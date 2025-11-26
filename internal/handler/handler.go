package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"shorturl.com/internal/service"
	"shorturl.com/pkg/config"
	"shorturl.com/pkg/logger"
)

type Handler struct {
	service *service.ShortenerService
	config  *config.Config
	logger  logger.Logger
}

func NewHandler(service *service.ShortenerService, config *config.Config, logger logger.Logger) *Handler {
	return &Handler{
		service: service,
		config:  config,
		logger:  logger,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	router.LoadHTMLGlob("web/templates/*")
	router.Static("/static", "./web/static")

	router.GET("/", h.home)
	router.GET("/:shortCode", h.redirect)
	router.GET("/:shortCode/stats", h.stats)

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

	h.logger.Debug("redirect attempt",
		logger.String("short_code", shortCode),
		logger.String("ip", c.ClientIP()),
	)

	originalURL, err := h.service.GetURL(shortCode)
	if err != nil {
		h.logger.Warn("URL not found",
			logger.String("short_code", shortCode),
			logger.Err(err),
		)
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"message": "URL not found",
		})
		return
	}
	h.logger.Info("redirect successful",
		logger.String("short_code", shortCode),
		logger.String("original_url", originalURL),
		logger.String("ip", c.ClientIP()),
	)
	c.Redirect(http.StatusMovedPermanently, originalURL)
}

func (h *Handler) stats(c *gin.Context) {
	shortCode := c.Param("shortCode")

	h.logger.Debug("stats request",
		logger.String("short_code", shortCode),
	)

	clickCount, createdAt, err := h.service.GetStats(shortCode)
	if err != nil {
		h.logger.Warn("stats not found",
			logger.String("short_code", shortCode),
			logger.Err(err),
		)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "URL not found",
		})
		return
	}

	h.logger.Info("stats retrieved",
		logger.String("short_code", shortCode),
		logger.Int("click_count", clickCount),
	)

	c.JSON(http.StatusOK, gin.H{
		"short_code":       shortCode,
		"click_count":      clickCount,
		"created_at":       createdAt.Format(time.RFC3339),
		"created_at_human": createdAt.Format("02.01.2006 15:04"),
	})
}

func (h *Handler) shorten(c *gin.Context) {
	var request struct {
		URL string `json:"url" binding:"required,url"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Warn("invalid URL request",
			logger.String("url", request.URL),
			logger.Err(err),
			logger.String("ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid URL",
		})
		return
	}

	h.logger.Info("shorten request",
		logger.String("original_url", request.URL),
		logger.String("ip", c.ClientIP()),
	)

	shortCode, err := h.service.SaveURL(request.URL)
	if err != nil {
		h.logger.Error("failed to save URL",
			logger.String("url", request.URL),
			logger.Err(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "server error",
		})
		return
	}

	shortURL := h.config.BaseURL + "/" + shortCode

	h.logger.Info("URL shortened successfully",
		logger.String("original_url", request.URL),
		logger.String("short_code", shortCode),
		logger.String("short_url", shortURL),
	)

	c.JSON(http.StatusOK, gin.H{
		"short_url":  shortURL,
		"short_code": shortCode,
	})

}
