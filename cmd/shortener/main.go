package main

import (
	"github.com/gin-gonic/gin"
	"log"
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
	router.POST("/", h.UpdateAndGetShort)
	router.GET("/:id", h.GetLongURL)

	err := router.Run(conf.Host)
	if err != nil {
		log.Fatalf("can't run http server: %v", err)
		return
	}
}
