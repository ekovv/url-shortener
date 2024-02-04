package handler

import "github.com/gin-gonic/gin"

// Route sa
func Route(c *gin.Engine, h *Handler) {
	c.POST("/", h.UpdateAndGetShort)
	c.GET("/:id", h.GetLongURL)
	c.POST("/api/shorten", h.GetShortByJSON)
	c.GET("/ping")
	c.POST("/api/shorten/batch", h.GetBatch)
	c.GET("/api/user/urls", h.GetAll)
	c.DELETE("/api/user/urls", h.Del)
}
