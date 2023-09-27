package main

import (
	"github.com/gin-contrib/gzip"
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
	stM := storage.NewStorage()
	stF := storage.NewFileStorage(conf.Storage)
	sr := service.NewService(stM, *stF, *conf)
	h := handler.NewHandler(&sr)
	router := gin.Default()
	router.Use(h.Decompressed())
	router.Use(gzip.Gzip(gzip.DefaultCompression))
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
