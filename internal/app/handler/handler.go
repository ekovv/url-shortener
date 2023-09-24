package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"url-shortener/internal/app/domains"
)

type Handler struct {
	service domains.UseCase
}

func NewHandler(service domains.UseCase) *Handler {
	return &Handler{service: service}
}

func (s *Handler) UpdateAndGetShort(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return

	}
	str := string(body)
	fmt.Println(str)
	short, err := s.service.GetShort(str)
	fmt.Println(short)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.String(http.StatusCreated, short)
}

func (s *Handler) GetLongURL(c *gin.Context) {
	id := c.Param("id")
	long, err := s.service.GetLong(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	c.Status(http.StatusTemporaryRedirect)
	c.Header("Location", long)

}
