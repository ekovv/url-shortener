package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"url-shortener/config"
	"url-shortener/internal/app/handler"
	myLog "url-shortener/internal/app/logger"
	"url-shortener/internal/app/service"
	"url-shortener/internal/app/storage"
)

func main() {
	conf := config.New()
	st := storage.NewStorage()
	sr := service.NewService(st, *conf)
	h := handler.NewHandler(&sr)
	router := gin.Default()
	router.Use(myLog.HTTPLogger())
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	myLog.Sugar = *logger.Sugar()
	router.POST("/", h.UpdateAndGetShort)
	router.GET("/:id", h.GetLongURL)
	router.POST("/api/shorten", h.GetShortByJSON)

	err = router.Run(conf.Host)
	if err != nil {
		log.Fatalf("can't run http server: %v", err)
		return
	}
}
