package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/http/httptest"
	"url-shortener/internal/service"
)

func ExampleHandler_GetConnection() {
	s := service.Service{}
	h := Handler{service: &s}

	// Создаем новый контекст Gin
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Вызываем функцию GetConnection
	h.GetConnection(c)

	// Проверяем статус ответа
	if w.Code != http.StatusOK {
		log.Fatalf("Expected HTTP 200 OK, got: %v", w.Code)
	}
	// Check connection
}
