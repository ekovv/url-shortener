package main

import (
	"fmt"
	"go.uber.org/zap"
	"url-shortener/config"
	"url-shortener/internal/app/handler"
	myLog "url-shortener/internal/app/logger"
	"url-shortener/internal/app/service"
	"url-shortener/internal/app/storage"
)

func main() {
	conf := config.New()
	stM := storage.NewStorage()
	stF, err := storage.NewFileStorage(conf.Storage)
	if err != nil {
		fmt.Println(err)
		return
	}
	sr := service.NewService(stM, *stF, *conf)
	h := handler.NewHandler(&sr, *conf)

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	myLog.Sugar = *logger.Sugar()
	h.Start()
}
