package main

import (
	"github.com/gin-gonic/gin"
	"url-shortener/config"
	"url-shortener/internal/app/handler"
	"url-shortener/internal/app/service"
	"url-shortener/internal/app/storage"
)

func main() {
	conf := config.New()
	st := storage.NewStorage()
	sr := service.NewService(st, *conf)
	h := handler.NewHandler(&sr)
	router := gin.Default()
	router.POST("/", h.UpdateAndRetShort)
	router.GET("/:id", h.GetLongURL)

	router.Run(conf.Host)
}
