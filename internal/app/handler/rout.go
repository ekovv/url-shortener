package handler

import "github.com/gin-gonic/gin"

func Route(c *gin.Engine, h *Handler) {
	c.POST("/", h.UpdateAndGetShort)
	c.GET("/:id", h.GetLongURL)
	c.POST("/api/shorten", h.GetShortByJSON)
}
