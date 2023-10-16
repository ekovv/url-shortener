package main

import (
	"fmt"
	"go.uber.org/zap"
	"log"
	"url-shortener/config"
	"url-shortener/internal/handler"
	myLog "url-shortener/internal/logger"
	"url-shortener/internal/service"
	"url-shortener/internal/storage"
)

func main() {
	conf := config.New()
	stM, err := storage.New(conf)
	if err != nil {
		fmt.Println(err)
		return
	}
	sr, err := service.NewService(stM, conf)
	if err != nil {
		log.Fatalf("Error creating service: %s", err)
		return
	}
	h := handler.NewHandler(&sr, conf)
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	myLog.Sugar = *logger.Sugar()
	h.Start()
	stM.Close()
}
