package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"url-shortener/internal/app/service"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

func (s *Handler) UpdateAndRetShort(c *gin.Context) {
	body, _ := c.GetRawData()
	str := string(body)
	fmt.Println(str)
	short, err := s.service.RetShort(str)
	fmt.Println(short)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.String(http.StatusCreated, short)
}

func (s *Handler) GetLongURL(c *gin.Context) {
	id := c.Param("id")
	long, err := s.service.RetLong(id)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.Status(http.StatusFound)

	// установка заголовка Location на нужный URL
	c.Header("Location", long)

}
