package main

import (
	"github.com/gin-gonic/gin"
	"url-shortener/internal/app/handler"
	"url-shortener/internal/app/service"
	"url-shortener/internal/app/storage"
)

func main() {
	st := storage.NewStorage()
	sr := service.NewService(st)
	h := handler.NewHandler(sr)
	router := gin.Default()
	router.POST("/", h.UpdateAndRetShort)
	router.GET("/:id", h.GetLongURL)

	router.Run("localhost:8080")
}
