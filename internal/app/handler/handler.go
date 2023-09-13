package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"url-shortener/internal/app/service"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

func (s *Handler) UpdateAndRetShort(c *gin.Context) {
	body, _ := c.GetRawData() // получаем тело запроса в виде []byte
	str := string(body)       // преобразовываем []byte в string
	fmt.Println(str)          // выводим полученную строку в консоль

	c.Status(200)
}
